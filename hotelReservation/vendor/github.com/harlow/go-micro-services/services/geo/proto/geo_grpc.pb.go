// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: geo.proto

package geo

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GeoClient is the client API for Geo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GeoClient interface {
	// Finds the hotels contained nearby the current lat/lon.
	Nearby(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error)
}

type geoClient struct {
	cc grpc.ClientConnInterface
}

func NewGeoClient(cc grpc.ClientConnInterface) GeoClient {
	return &geoClient{cc}
}

func (c *geoClient) Nearby(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/geo.Geo/Nearby", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GeoServer is the server API for Geo service.
// All implementations must embed UnimplementedGeoServer
// for forward compatibility
type GeoServer interface {
	// Finds the hotels contained nearby the current lat/lon.
	Nearby(context.Context, *Request) (*Result, error)
	mustEmbedUnimplementedGeoServer()
}

// UnimplementedGeoServer must be embedded to have forward compatible implementations.
type UnimplementedGeoServer struct {
}

func (UnimplementedGeoServer) Nearby(context.Context, *Request) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Nearby not implemented")
}
func (UnimplementedGeoServer) mustEmbedUnimplementedGeoServer() {}

// UnsafeGeoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GeoServer will
// result in compilation errors.
type UnsafeGeoServer interface {
	mustEmbedUnimplementedGeoServer()
}

func RegisterGeoServer(s grpc.ServiceRegistrar, srv GeoServer) {
	s.RegisterService(&Geo_ServiceDesc, srv)
}

func _Geo_Nearby_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).Nearby(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geo.Geo/Nearby",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).Nearby(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Geo_ServiceDesc is the grpc.ServiceDesc for Geo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Geo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "geo.Geo",
	HandlerType: (*GeoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Nearby",
			Handler:    _Geo_Nearby_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "geo.proto",
}
