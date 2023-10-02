// Code generated by protoc-gen-go. DO NOT EDIT.
// source: geo.proto

/*
Package __geo is a generated protocol buffer package.

It is generated from these files:
	geo.proto

It has these top-level messages:
	Request
	Result
*/
package geo

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The latitude and longitude of the current location.
type Request struct {
	Lat float32 `protobuf:"fixed32,1,opt,name=lat" json:"lat,omitempty"`
	Lon float32 `protobuf:"fixed32,2,opt,name=lon" json:"lon,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetLat() float32 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *Request) GetLon() float32 {
	if m != nil {
		return m.Lon
	}
	return 0
}

type Result struct {
	HotelIds []string `protobuf:"bytes,1,rep,name=hotelIds" json:"hotelIds,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetHotelIds() []string {
	if m != nil {
		return m.HotelIds
	}
	return nil
}

func init() {
	proto.RegisterType((*Request)(nil), "geo.Request")
	proto.RegisterType((*Result)(nil), "geo.Result")
}

func init() { proto.RegisterFile("geo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x4f, 0xcd, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x4f, 0xcd, 0x57, 0xd2, 0xe5, 0x62, 0x0f, 0x4a,
	0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x12, 0xe0, 0x62, 0xce, 0x49, 0x2c, 0x91, 0x60, 0x54, 0x60,
	0xd4, 0x60, 0x0a, 0x02, 0x31, 0xc1, 0x22, 0xf9, 0x79, 0x12, 0x4c, 0x50, 0x91, 0xfc, 0x3c, 0x25,
	0x15, 0x2e, 0xb6, 0xa0, 0xd4, 0xe2, 0xd2, 0x9c, 0x12, 0x21, 0x29, 0x2e, 0x8e, 0x8c, 0xfc, 0x92,
	0xd4, 0x1c, 0xcf, 0x94, 0x62, 0x09, 0x46, 0x05, 0x66, 0x0d, 0xce, 0x20, 0x38, 0xdf, 0x48, 0x8b,
	0x8b, 0xd9, 0x3d, 0x35, 0x5f, 0x48, 0x99, 0x8b, 0xcd, 0x2f, 0x35, 0xb1, 0x28, 0xa9, 0x52, 0x88,
	0x47, 0x0f, 0x64, 0x2d, 0xd4, 0x22, 0x29, 0x6e, 0x28, 0x0f, 0x64, 0x8e, 0x13, 0x7b, 0x14, 0xab,
	0x9e, 0x75, 0x7a, 0x6a, 0x7e, 0x12, 0x1b, 0xd8, 0x55, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xa7, 0xde, 0x87, 0x5d, 0xa2, 0x00, 0x00, 0x00,
}
