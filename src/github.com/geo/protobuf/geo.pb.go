// Code generated by protoc-gen-go. DO NOT EDIT.
// source: geo.proto

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	geo.proto

It has these top-level messages:
	Content
	IP
	MaxmindCityInfo
	IpipCityInfo
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Content struct {
	Body    []byte            `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	Headers map[string]string `protobuf:"bytes,2,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Content) Reset()                    { *m = Content{} }
func (m *Content) String() string            { return proto.CompactTextString(m) }
func (*Content) ProtoMessage()               {}
func (*Content) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Content) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Content) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

type IP struct {
	IP   string `protobuf:"bytes,1,opt,name=IP" json:"IP,omitempty"`
	Lang string `protobuf:"bytes,2,opt,name=Lang" json:"Lang,omitempty"`
}

func (m *IP) Reset()                    { *m = IP{} }
func (m *IP) String() string            { return proto.CompactTextString(m) }
func (*IP) ProtoMessage()               {}
func (*IP) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *IP) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *IP) GetLang() string {
	if m != nil {
		return m.Lang
	}
	return ""
}

type MaxmindCityInfo struct {
	Message     string  `protobuf:"bytes,1,opt,name=Message" json:"Message,omitempty"`
	City        string  `protobuf:"bytes,2,opt,name=City" json:"City,omitempty"`
	Subdivision string  `protobuf:"bytes,3,opt,name=Subdivision" json:"Subdivision,omitempty"`
	Country     string  `protobuf:"bytes,4,opt,name=Country" json:"Country,omitempty"`
	Zone        string  `protobuf:"bytes,5,opt,name=Zone" json:"Zone,omitempty"`
	Latitude    float64 `protobuf:"fixed64,6,opt,name=Latitude" json:"Latitude,omitempty"`
	Longitude   float64 `protobuf:"fixed64,7,opt,name=Longitude" json:"Longitude,omitempty"`
}

func (m *MaxmindCityInfo) Reset()                    { *m = MaxmindCityInfo{} }
func (m *MaxmindCityInfo) String() string            { return proto.CompactTextString(m) }
func (*MaxmindCityInfo) ProtoMessage()               {}
func (*MaxmindCityInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MaxmindCityInfo) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *MaxmindCityInfo) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *MaxmindCityInfo) GetSubdivision() string {
	if m != nil {
		return m.Subdivision
	}
	return ""
}

func (m *MaxmindCityInfo) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *MaxmindCityInfo) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func (m *MaxmindCityInfo) GetLatitude() float64 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *MaxmindCityInfo) GetLongitude() float64 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

type IpipCityInfo struct {
	Message   string `protobuf:"bytes,1,opt,name=Message" json:"Message,omitempty"`
	Country   string `protobuf:"bytes,2,opt,name=Country" json:"Country,omitempty"`
	Province  string `protobuf:"bytes,3,opt,name=Province" json:"Province,omitempty"`
	City      string `protobuf:"bytes,4,opt,name=City" json:"City,omitempty"`
	Org       string `protobuf:"bytes,5,opt,name=Org" json:"Org,omitempty"`
	ISP       string `protobuf:"bytes,6,opt,name=ISP" json:"ISP,omitempty"`
	Latitude  string `protobuf:"bytes,7,opt,name=Latitude" json:"Latitude,omitempty"`
	Longitude string `protobuf:"bytes,8,opt,name=Longitude" json:"Longitude,omitempty"`
	TimeZone  string `protobuf:"bytes,9,opt,name=TimeZone" json:"TimeZone,omitempty"`
	UTC       string `protobuf:"bytes,10,opt,name=UTC" json:"UTC,omitempty"`
	ChinaCode string `protobuf:"bytes,11,opt,name=ChinaCode" json:"ChinaCode,omitempty"`
	PhoneCode string `protobuf:"bytes,12,opt,name=PhoneCode" json:"PhoneCode,omitempty"`
	ISO2      string `protobuf:"bytes,13,opt,name=ISO2" json:"ISO2,omitempty"`
	Continent string `protobuf:"bytes,14,opt,name=Continent" json:"Continent,omitempty"`
}

func (m *IpipCityInfo) Reset()                    { *m = IpipCityInfo{} }
func (m *IpipCityInfo) String() string            { return proto.CompactTextString(m) }
func (*IpipCityInfo) ProtoMessage()               {}
func (*IpipCityInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *IpipCityInfo) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *IpipCityInfo) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *IpipCityInfo) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *IpipCityInfo) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *IpipCityInfo) GetOrg() string {
	if m != nil {
		return m.Org
	}
	return ""
}

func (m *IpipCityInfo) GetISP() string {
	if m != nil {
		return m.ISP
	}
	return ""
}

func (m *IpipCityInfo) GetLatitude() string {
	if m != nil {
		return m.Latitude
	}
	return ""
}

func (m *IpipCityInfo) GetLongitude() string {
	if m != nil {
		return m.Longitude
	}
	return ""
}

func (m *IpipCityInfo) GetTimeZone() string {
	if m != nil {
		return m.TimeZone
	}
	return ""
}

func (m *IpipCityInfo) GetUTC() string {
	if m != nil {
		return m.UTC
	}
	return ""
}

func (m *IpipCityInfo) GetChinaCode() string {
	if m != nil {
		return m.ChinaCode
	}
	return ""
}

func (m *IpipCityInfo) GetPhoneCode() string {
	if m != nil {
		return m.PhoneCode
	}
	return ""
}

func (m *IpipCityInfo) GetISO2() string {
	if m != nil {
		return m.ISO2
	}
	return ""
}

func (m *IpipCityInfo) GetContinent() string {
	if m != nil {
		return m.Continent
	}
	return ""
}

func init() {
	proto.RegisterType((*Content)(nil), "protobuf.content")
	proto.RegisterType((*IP)(nil), "protobuf.IP")
	proto.RegisterType((*MaxmindCityInfo)(nil), "protobuf.MaxmindCityInfo")
	proto.RegisterType((*IpipCityInfo)(nil), "protobuf.IpipCityInfo")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Geo service

type GeoClient interface {
	Echo(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error)
	GetMaxmindCity(ctx context.Context, in *IP, opts ...grpc.CallOption) (*MaxmindCityInfo, error)
	GetIpipCity(ctx context.Context, in *IP, opts ...grpc.CallOption) (*IpipCityInfo, error)
}

type geoClient struct {
	cc *grpc.ClientConn
}

func NewGeoClient(cc *grpc.ClientConn) GeoClient {
	return &geoClient{cc}
}

func (c *geoClient) Echo(ctx context.Context, in *Content, opts ...grpc.CallOption) (*Content, error) {
	out := new(Content)
	err := grpc.Invoke(ctx, "/protobuf.Geo/echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoClient) GetMaxmindCity(ctx context.Context, in *IP, opts ...grpc.CallOption) (*MaxmindCityInfo, error) {
	out := new(MaxmindCityInfo)
	err := grpc.Invoke(ctx, "/protobuf.Geo/GetMaxmindCity", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geoClient) GetIpipCity(ctx context.Context, in *IP, opts ...grpc.CallOption) (*IpipCityInfo, error) {
	out := new(IpipCityInfo)
	err := grpc.Invoke(ctx, "/protobuf.Geo/GetIpipCity", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Geo service

type GeoServer interface {
	Echo(context.Context, *Content) (*Content, error)
	GetMaxmindCity(context.Context, *IP) (*MaxmindCityInfo, error)
	GetIpipCity(context.Context, *IP) (*IpipCityInfo, error)
}

func RegisterGeoServer(s *grpc.Server, srv GeoServer) {
	s.RegisterService(&_Geo_serviceDesc, srv)
}

func _Geo_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Geo/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).Echo(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

func _Geo_GetMaxmindCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).GetMaxmindCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Geo/GetMaxmindCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).GetMaxmindCity(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

func _Geo_GetIpipCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IP)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeoServer).GetIpipCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Geo/GetIpipCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeoServer).GetIpipCity(ctx, req.(*IP))
	}
	return interceptor(ctx, in, info, handler)
}

var _Geo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Geo",
	HandlerType: (*GeoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "echo",
			Handler:    _Geo_Echo_Handler,
		},
		{
			MethodName: "GetMaxmindCity",
			Handler:    _Geo_GetMaxmindCity_Handler,
		},
		{
			MethodName: "GetIpipCity",
			Handler:    _Geo_GetIpipCity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "geo.proto",
}

func init() { proto.RegisterFile("geo.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 481 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x6e, 0xd3, 0x40,
	0x14, 0xae, 0x1d, 0xb7, 0x8e, 0x5f, 0x4c, 0x28, 0x23, 0x84, 0x06, 0x0b, 0xa1, 0xc8, 0xab, 0xac,
	0xbc, 0x08, 0x42, 0xaa, 0xca, 0xd2, 0x42, 0xc1, 0x52, 0xaa, 0x5a, 0x4e, 0xba, 0x61, 0xe7, 0xc4,
	0x53, 0x67, 0x04, 0x9d, 0x89, 0x9c, 0x71, 0x44, 0xce, 0xc1, 0x29, 0x38, 0x0a, 0x77, 0xe1, 0x10,
	0xe8, 0xcd, 0xd8, 0x8e, 0xdb, 0xb2, 0x60, 0x93, 0x7c, 0x3f, 0xef, 0x7d, 0x33, 0x5f, 0x32, 0xe0,
	0x95, 0x4c, 0x46, 0xbb, 0x4a, 0x2a, 0x49, 0x86, 0xfa, 0x6b, 0x5d, 0xdf, 0x87, 0x3f, 0x2d, 0x70,
	0x37, 0x52, 0x28, 0x26, 0x14, 0x21, 0xe0, 0xac, 0x65, 0x71, 0xa4, 0xd6, 0xc4, 0x9a, 0xfa, 0x99,
	0xc6, 0xe4, 0x0a, 0xdc, 0x2d, 0xcb, 0x0b, 0x56, 0xed, 0xa9, 0x3d, 0x19, 0x4c, 0x47, 0xb3, 0xf7,
	0x51, 0xbb, 0x1b, 0x35, 0x7b, 0xd1, 0x17, 0x33, 0xf0, 0x59, 0xa8, 0xea, 0x98, 0xb5, 0xe3, 0xc1,
	0x35, 0xf8, 0x7d, 0x83, 0x5c, 0xc2, 0xe0, 0x1b, 0x33, 0xe1, 0x5e, 0x86, 0x90, 0xbc, 0x86, 0xf3,
	0x43, 0xfe, 0xbd, 0x66, 0xd4, 0xd6, 0x9a, 0x21, 0xd7, 0xf6, 0x95, 0x15, 0x4e, 0xc1, 0x4e, 0x52,
	0x32, 0xc6, 0xcf, 0x66, 0x01, 0x39, 0x01, 0x67, 0x91, 0x8b, 0xb2, 0x19, 0xd7, 0x38, 0xfc, 0x6d,
	0xc1, 0xcb, 0x9b, 0xfc, 0xc7, 0x03, 0x17, 0x45, 0xcc, 0xd5, 0x31, 0x11, 0xf7, 0x92, 0x50, 0x70,
	0x6f, 0xd8, 0x7e, 0x9f, 0x97, 0xac, 0x59, 0x6e, 0x29, 0x26, 0xe0, 0x54, 0x9b, 0x80, 0x98, 0x4c,
	0x60, 0xb4, 0xac, 0xd7, 0x05, 0x3f, 0xf0, 0x3d, 0x97, 0x82, 0x0e, 0xb4, 0xd5, 0x97, 0x30, 0x2f,
	0x96, 0x35, 0x96, 0xa0, 0x8e, 0xc9, 0x6b, 0x28, 0xe6, 0x7d, 0x95, 0x82, 0xd1, 0x73, 0x93, 0x87,
	0x98, 0x04, 0x30, 0x5c, 0xe4, 0x8a, 0xab, 0xba, 0x60, 0xf4, 0x62, 0x62, 0x4d, 0xad, 0xac, 0xe3,
	0xe4, 0x1d, 0x78, 0x0b, 0x29, 0x4a, 0x63, 0xba, 0xda, 0x3c, 0x09, 0xe1, 0x1f, 0x1b, 0xfc, 0x64,
	0xc7, 0x77, 0xff, 0x51, 0xa4, 0x77, 0x25, 0xfb, 0xf1, 0x95, 0x02, 0x18, 0xa6, 0x95, 0x3c, 0x70,
	0xb1, 0x61, 0x4d, 0x97, 0x8e, 0x77, 0xf5, 0x9d, 0x5e, 0xfd, 0x4b, 0x18, 0xdc, 0x56, 0x65, 0xd3,
	0x00, 0x21, 0x2a, 0xc9, 0x32, 0xd5, 0x77, 0xf7, 0x32, 0x84, 0x8f, 0x2a, 0xb9, 0x26, 0xf3, 0xdf,
	0x95, 0x86, 0xda, 0x3c, 0x09, 0xb8, 0xb9, 0xe2, 0x0f, 0x4c, 0xff, 0x48, 0x9e, 0xd9, 0x6c, 0x39,
	0x9e, 0x73, 0xb7, 0x8a, 0x29, 0x98, 0x73, 0xee, 0x56, 0x31, 0x66, 0xc5, 0x5b, 0x2e, 0xf2, 0x58,
	0x16, 0x8c, 0x8e, 0x4c, 0x56, 0x27, 0xa0, 0x9b, 0x6e, 0xa5, 0x60, 0xda, 0xf5, 0x8d, 0xdb, 0x09,
	0xd8, 0x2d, 0x59, 0xde, 0xce, 0xe8, 0x0b, 0xd3, 0x0d, 0xb1, 0xce, 0x93, 0x42, 0x71, 0xc1, 0x84,
	0xa2, 0xe3, 0x26, 0xaf, 0x15, 0x66, 0xbf, 0x2c, 0x18, 0xcc, 0x99, 0x24, 0x11, 0x38, 0x6c, 0xb3,
	0x95, 0xe4, 0xd5, 0xb3, 0x97, 0x1d, 0x3c, 0x97, 0xc2, 0x33, 0xf2, 0x09, 0xc6, 0x73, 0xa6, 0x7a,
	0x8f, 0x8e, 0xf8, 0xa7, 0xb1, 0x24, 0x0d, 0xde, 0x9e, 0xd8, 0x93, 0x97, 0x19, 0x9e, 0x91, 0x8f,
	0x30, 0x9a, 0x33, 0xd5, 0xfe, 0xcb, 0x4f, 0x36, 0xdf, 0xf4, 0x58, 0xef, 0x1d, 0x84, 0x67, 0xeb,
	0x0b, 0x6d, 0x7c, 0xf8, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x90, 0x19, 0xfd, 0xc4, 0x03, 0x00,
	0x00,
}
