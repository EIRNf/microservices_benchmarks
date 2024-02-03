package search

import (
	// "encoding/json"
	"fmt"

	"github.com/EIRNf/notnets_grpc"
	"github.com/harlow/go-micro-services/dialer"

	// F"io/ioutil"

	"github.com/rs/zerolog/log"

	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/harlow/go-micro-services/registry"
	geo "github.com/harlow/go-micro-services/services/geo/proto"
	rate "github.com/harlow/go-micro-services/services/rate/proto"
	pb "github.com/harlow/go-micro-services/services/search/proto"
	"github.com/harlow/go-micro-services/tls"
	opentracing "github.com/opentracing/opentracing-go"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/keepalive"
)

const name = "srv-search"

// Server implments the search service
type Server struct {
	geoClient  geo.GeoClient
	rateClient rate.RateClient

	Tracer     opentracing.Tracer
	Port       int
	IpAddr     string
	KnativeDns string
	Registry   *registry.Client
	uuid       string

	pb.UnimplementedSearchServer
}

// Run starts the server
func (s *Server) Run() error {

	if s.Port == 0 {
		return fmt.Errorf("server port must be set")
	}

	s.uuid = uuid.New().String()

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: 120 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.Tracer),
		),
	}

	if tlsopt := tls.GetServerOpt(); tlsopt != nil {
		opts = append(opts, tlsopt)
	}

	// svc := &SearchServer{}
	// srv := grpc.NewServer(opts...)
	srv := notnets_grpc.NewNotnetsServer()
	pb.RegisterSearchServer(srv, s)

	var err error
	// init grpc clients
	if err := s.initGeoClientShm("srv-geo"); err != nil {
		return err
	}
	if err := s.initRateClient("srv-rate"); err != nil {
		return err
	}

	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	// if err != nil {
	// 	log.Fatal().Msgf("failed to listen: %v", err)
	// }
	lis := notnets_grpc.Listen(name)

	// register with consul
	// jsonFile, err := os.Open("config.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer jsonFile.Close()

	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// var result map[string]string
	// json.Unmarshal([]byte(byteValue), &result)

	err = s.Registry.Register(name, s.uuid, s.IpAddr, s.Port)
	if err != nil {
		return fmt.Errorf("failed register: %v", err)
	}
	log.Info().Msg("Successfully registered in consul")

	return srv.Serve(lis)
}

// Shutdown cleans up any processes
func (s *Server) Shutdown() {
	s.Registry.Deregister(s.uuid)
}

// func (s *Server) initGeoClient(name string) error {
// 	conn, err := s.getGprcConn(name)
// 	if err != nil {
// 		return fmt.Errorf("dialer error: %v", err)
// 	}
// 	s.geoClient = geo.NewGeoClient(conn)
// 	return nil
// }

func (s *Server) initGeoClientShm(name string) error {
	// Construct Channel with necessary parameters to talk to the Server
	// cc := shmgrpc.NewChannel(s.IpAddr+":"+fmt.Sprint(s.Port), name)

	// s.Registry.Config.Address
	// cc, err := notnets_grpc.Dial(s.IpAddr, s.IpAddr+":"+fmt.Sprint(s.Port)+name)
	cc, err := notnets_grpc.Dial(s.IpAddr+"srv-search", name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	time.Sleep(10 * time.Second)
	s.geoClient = geo.NewGeoClient(cc)
	return nil
}

func (s *Server) initRateClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.rateClient = rate.NewRateClient(conn)
	return nil
}

func (s *Server) getGprcConn(name string) (*grpc.ClientConn, error) {

	// Make another ClientConn with round_robin policy.
	// return grpc.Dial(
	// 	fmt.Sprintf("%s:///%s", resolver.GetDefaultScheme(), name),
	// 	grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )

	//Use dialer to minimize code reuse
	return dialer.Dial(
		name,
		s.Registry,
		s.Tracer,
	)
	// if s.KnativeDns != "" {
	// 	return dialer.Dial(
	// 		fmt.Sprintf("%s.%s", name, s.KnativeDns),
	// 		dialer.WithTracer(s.Tracer))
	// } else {
	// 	return dialer.Dial(
	// 		name,
	// 		dialer.WithTracer(s.Tracer),
	// 		dialer.WithBalancer(s.Registry.Client),
	// 	)
	// }
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func (s *Server) Nearby(ctx context.Context, req *pb.NearbyRequest) (*pb.SearchResult, error) {
	// find nearby hotels
	log.Trace().Msg("in Search Nearby")

	log.Trace().Msgf("nearby lat = %f", req.Lat)
	log.Trace().Msgf("nearby lon = %f", req.Lon)

	nearby, err := s.geoClient.Nearby(ctx, &geo.Request{
		Lat: req.Lat,
		Lon: req.Lon,
	})
	if err != nil {
		return nil, err
	}

	for _, hid := range nearby.HotelIds {
		log.Trace().Msgf("get Nearby hotelId = %s", hid)
	}

	// find rates for hotels
	rates, err := s.rateClient.GetRates(ctx, &rate.Request{
		HotelIds: nearby.HotelIds,
		InDate:   req.InDate,
		OutDate:  req.OutDate,
	})
	if err != nil {
		return nil, err
	}

	// TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
	// * price (best discount?)
	// * reviews

	// build the response
	res := new(pb.SearchResult)
	for _, ratePlan := range rates.RatePlans {
		log.Trace().Msgf("get RatePlan HotelId = %s, Code = %s", ratePlan.HotelId, ratePlan.Code)
		res.HotelIds = append(res.HotelIds, ratePlan.HotelId)
	}
	return res, nil
}
