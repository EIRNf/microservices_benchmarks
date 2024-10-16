package shmgrpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpchan"
	"github.com/fullstorydev/grpchan/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	grpcproto "google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// type ServeMux struct {
// 	mu    sync.RWMutex
// 	m     map[string]muxEntry
// 	es    []muxEntry
// 	hosts bool
// }

// type muxEntry struct {
// 	h       Handler
// 	pattern string
// }

// type Handler interface {
// 	ServeMethod(,)
// }

// Server is grpc over shared memory. It
// acts as a grpc ServiceRegistrar
type Server struct {
	// mux              ServeMux
	handlers         grpchan.HandlerMap
	basePath         string
	ShmQueueInfo     *QueueInfo
	responseQueue    *Queue
	requestQeuue     *Queue
	opts             handlerOpts
	unaryInterceptor grpc.UnaryServerInterceptor
	// shutdownChannel  chan bool
}

var _ grpc.ServiceRegistrar = (*Server)(nil)

var (
	sSerReqStruct   ShmMessage
	sSerReqWritten  bool = false
	sSerRespData    [600]byte
	sSerRespLen     int
	sSerRespWritten bool = false

	sserPayload        []byte
	sserPayloadWritten bool = false
)

// ServerOption is an option used when constructing a NewServer.
type ServerOption interface {
	apply(*Server)
}

type serverOptFunc func(*Server)

func (fn serverOptFunc) apply(s *Server) {
	fn(s)
}

type handlerOpts struct {
	errFunc func(context.Context, *status.Status, http.ResponseWriter)
}

func NewServer(basePath string, opts ...ServerOption) *Server {
	var s Server
	s.basePath = basePath
	s.handlers = grpchan.HandlerMap{}

	//Gather Keys based on name
	requestKey, responseKey := GatherShmKeys(basePath)

	//Initilize ShmQueueInfo
	//Client -> Server Shm
	requestShmid, requestShmaddr := InitializeShmRegion(requestKey, Size, uintptr(ServerSegFlag))
	//Server -> Client Shm
	responseShmid, responseShmaddr := InitializeShmRegion(responseKey, Size, uintptr(ServerSegFlag))

	qi := QueueInfo{
		RequestShmid:    requestShmid,
		RequestShmaddr:  requestShmaddr,
		ResponseShmid:   responseShmid,
		ResponseShmaddr: responseShmaddr,
	}

	s.ShmQueueInfo = &qi

	//Initiliaze Queue
	s.requestQeuue = initializeQueue(s.ShmQueueInfo.RequestShmaddr)
	s.responseQueue = initializeQueue(s.ShmQueueInfo.ResponseShmaddr)

	return &s
}

// Register Service, also gets generated by protoc, as Register(SERVICE NAME)Server
func (s *Server) RegisterService(desc *grpc.ServiceDesc, svr interface{}) {
	s.handlers.RegisterService(desc, svr)
}

func (s *Server) HandleMethods(svr interface{}) {

	s.unaryInterceptor = nil

	requestQueue := GetQueue(s.ShmQueueInfo.RequestShmaddr)
	responseQueue := GetQueue(s.ShmQueueInfo.ResponseShmaddr)

	for {

		message, err := consumeMessage(requestQueue)
		if err != nil {
			break
			//the channel has been shut down
		}

		messageTag := message.Header.Tag
		slice := message.Data[0:message.Header.Size]
		var message_req_meta ShmMessage
		if !sSerReqWritten {
			json.Unmarshal(slice, &message_req_meta)
			sSerReqStruct = message_req_meta
			sSerReqWritten = true
			if err != nil {
				// return err
				status.Errorf(codes.Unknown, "Codec Marshalling error: %s ", err.Error())
			}
		} else {
			message_req_meta = sSerReqStruct
		}

		//Parse bytes into object

		payload_buffer := unsafeGetBytes(message_req_meta.Payload)

		fullName := message_req_meta.Method
		strs := strings.SplitN(fullName[1:], "/", 2)
		serviceName := strs[0]
		methodName := strs[1]
		clientCtx := message_req_meta.Context
		if clientCtx == nil { // Temp in case of empty context.
			clientCtx = context.Background()
		}

		clientCtx, cancel, err := contextFromHeaders(clientCtx, message_req_meta.Headers)
		if err != nil {
			// writeError(w, http.StatusBadRequest)
			return
		}

		defer cancel()

		//Get Service Descriptor and Handler
		sd, handler := s.handlers.QueryService(serviceName)
		if sd == nil {
			// service name not found
			status.Errorf(codes.Unimplemented, "service %s not implemented", message_req_meta.Method)
		}

		//Get Method Descriptor
		md := internal.FindUnaryMethod(methodName, sd.Methods)
		if md == nil {
			// method name not found
			status.Errorf(codes.Unimplemented, "method %s/%s not implemented", serviceName, methodName)
		}

		//Get Codec for content type.
		codec := encoding.GetCodec(grpcproto.Name)

		// Function to unmarshal payload using proto
		dec := func(msg interface{}) error {
			val := payload_buffer
			if err := codec.Unmarshal(val, msg); err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
			return nil
		}

		// Implements server transport stream
		sts := internal.UnaryServerTransportStream{Name: methodName}
		ctx := grpc.NewContextWithServerTransportStream(makeServerContext(clientCtx), &sts)

		//Get resp write back
		resp, err := md.Handler(
			handler,
			ctx,
			dec,
			s.unaryInterceptor)
		if err != nil {
			status.Errorf(codes.Unknown, "Handler error: %s ", err.Error())
			//TODO: Error code must be sent back to client
		}

		var resp_buffer []byte
		if !sserPayloadWritten {
			resp_buffer, err := codec.Marshal(resp)
			if err != nil {
				// return err
			}
			sserPayload = resp_buffer
			sserPayloadWritten = true
		}

		resp_buffer = sserPayload
		if err != nil {
			status.Errorf(codes.Unknown, "Codec Marshalling error: %s ", err.Error())
		}

		message_resp := &ShmMessage{
			Method:   methodName,
			Context:  ctx,
			Headers:  sts.GetHeaders(),
			Trailers: sts.GetTrailers(),
			Payload:  ByteSlice2String(resp_buffer),
		}

		var serializedMessage []byte
		var data [600]byte
		if !sSerRespWritten {
			serializedMessage, err = json.Marshal(message_resp)
			sSerRespLen = copy(sSerRespData[:], serializedMessage)
			data = sSerRespData
			sSerRespWritten = true
			if err != nil {
				status.Errorf(codes.Unknown, "Codec Marshalling error: %s ", err.Error())
			}
		} else {
			data = sSerRespData
		}

		message_response := Message{
			Header: MessageHeader{
				Size: int32(sSerRespLen),
				Tag:  messageTag},
			Data: data,
		}

		produceMessage(responseQueue, message_response)

		if !NO_SERIALIZATION {
			sSerReqWritten = false
			sSerRespWritten = false
			sserPayloadWritten = false

		}

	}

}

// Shutdown the server
func (s *Server) Stop() {
	StopPollingQueue(s.requestQeuue)
	StopPollingQueue(s.responseQueue)

	defer Detach(s.ShmQueueInfo.RequestShmaddr)
	defer Detach(s.ShmQueueInfo.ResponseShmaddr)

	defer Remove(s.ShmQueueInfo.RequestShmid)
	defer Remove(s.ShmQueueInfo.ResponseShmid)
	// close(responseQueue.DetachQueue)

}

// contextFromHeaders returns a child of the given context that is populated
// using the given headers. The headers are converted to incoming metadata that
// can be retrieved via metadata.FromIncomingContext. If the headers contain a
// GRPC timeout, that is used to create a timeout for the returned context.
func contextFromHeaders(parent context.Context, md metadata.MD) (context.Context, context.CancelFunc, error) {
	cancel := func() {} // default to no-op

	ctx := metadata.NewIncomingContext(parent, md)

	// deadline propagation
	// timeout := h.Get("GRPC-Timeout")
	// if timeout != "" {
	// 	// See GRPC wire format, "Timeout" component of request: https://grpc.io/docs/guides/wire.html#requests
	// 	suffix := timeout[len(timeout)-1]
	// 	if timeoutVal, err := strconv.ParseInt(timeout[:len(timeout)-1], 10, 64); err == nil {
	// 		var unit time.Duration
	// 		switch suffix {
	// 		case 'H':
	// 			unit = time.Hour
	// 		case 'M':
	// 			unit = time.Minute
	// 		case 'S':
	// 			unit = time.Second
	// 		case 'm':
	// 			unit = time.Millisecond
	// 		case 'u':
	// 			unit = time.Microsecond
	// 		case 'n':
	// 			unit = time.Nanosecond
	// 		}
	// 		if unit != 0 {
	// 			ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutVal)*unit)
	// 		}
	// 	}
	// }
	return ctx, cancel, nil
}

// noValuesContext wraps a context but prevents access to its values. This is
// useful when you need a child context only to propagate cancellations and
// deadlines, but explicitly *not* to propagate values.
type noValuesContext struct {
	context.Context
}

func makeServerContext(ctx context.Context) context.Context {
	// We don't want the server have any of the values in the client's context
	// since that can inadvertently leak state from the client to the server.
	// But we do want a child context, just so that request deadlines and client
	// cancellations work seamlessly.
	newCtx := context.Context(noValuesContext{ctx})

	if meta, ok := metadata.FromOutgoingContext(ctx); ok {
		newCtx = metadata.NewIncomingContext(newCtx, meta)
	}
	// newCtx = peer.NewContext(newCtx, &inprocessPeer)
	// newCtx = context.WithValue(newCtx, &clientContextKey, ctx)
	return newCtx
}

// func handleMethod(svr interface{}, serviceName string, desc *grpc.MethodDesc) func() {

// 	fullMethod := fmt.Sprintf("/%s/%s", serviceName, desc.MethodName)
// 	fmt.Println(fullMethod)
// 	return f()
// }
