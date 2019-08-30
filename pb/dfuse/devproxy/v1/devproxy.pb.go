// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dfuse/devproxy/v1/devproxy.proto

package devproxy

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d14fe04e0ad063e, []int{0}
}

func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (m *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(m, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

type ListResponse struct {
	Servers              []string `protobuf:"bytes,1,rep,name=servers,proto3" json:"servers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}
func (*ListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4d14fe04e0ad063e, []int{1}
}

func (m *ListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListResponse.Unmarshal(m, b)
}
func (m *ListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListResponse.Marshal(b, m, deterministic)
}
func (m *ListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListResponse.Merge(m, src)
}
func (m *ListResponse) XXX_Size() int {
	return xxx_messageInfo_ListResponse.Size(m)
}
func (m *ListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListResponse proto.InternalMessageInfo

func (m *ListResponse) GetServers() []string {
	if m != nil {
		return m.Servers
	}
	return nil
}

func init() {
	proto.RegisterType((*ListRequest)(nil), "dfuse.devproxy.v1.ListRequest")
	proto.RegisterType((*ListResponse)(nil), "dfuse.devproxy.v1.ListResponse")
}

func init() { proto.RegisterFile("dfuse/devproxy/v1/devproxy.proto", fileDescriptor_4d14fe04e0ad063e) }

var fileDescriptor_4d14fe04e0ad063e = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0x49, 0x2b, 0x2d,
	0x4e, 0xd5, 0x4f, 0x49, 0x2d, 0x2b, 0x28, 0xca, 0xaf, 0xa8, 0xd4, 0x2f, 0x33, 0x84, 0xb3, 0xf5,
	0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x04, 0xc1, 0x2a, 0xf4, 0xe0, 0xa2, 0x65, 0x86, 0x4a, 0xbc,
	0x5c, 0xdc, 0x3e, 0x99, 0xc5, 0x25, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x4a, 0x1a, 0x5c,
	0x3c, 0x10, 0x6e, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x90, 0x04, 0x17, 0x7b, 0x71, 0x6a, 0x51,
	0x59, 0x6a, 0x51, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x67, 0x10, 0x8c, 0x6b, 0x14, 0xc5, 0xc5,
	0xe1, 0x02, 0x35, 0x47, 0xc8, 0x0f, 0x62, 0x48, 0x30, 0x44, 0x4a, 0x48, 0x4e, 0x0f, 0xc3, 0x1e,
	0x3d, 0x24, 0x4b, 0xa4, 0xe4, 0x71, 0xca, 0x43, 0x6c, 0x75, 0xe2, 0x8a, 0xe2, 0x80, 0xc9, 0x25,
	0xb1, 0x81, 0x9d, 0x6e, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x5f, 0x64, 0x32, 0x5d, 0xde, 0x00,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DevproxyClient is the client API for Devproxy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DevproxyClient interface {
	ListServers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type devproxyClient struct {
	cc *grpc.ClientConn
}

func NewDevproxyClient(cc *grpc.ClientConn) DevproxyClient {
	return &devproxyClient{cc}
}

func (c *devproxyClient) ListServers(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/dfuse.devproxy.v1.Devproxy/ListServers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DevproxyServer is the server API for Devproxy service.
type DevproxyServer interface {
	ListServers(context.Context, *ListRequest) (*ListResponse, error)
}

func RegisterDevproxyServer(s *grpc.Server, srv DevproxyServer) {
	s.RegisterService(&_Devproxy_serviceDesc, srv)
}

func _Devproxy_ListServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevproxyServer).ListServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dfuse.devproxy.v1.Devproxy/ListServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevproxyServer).ListServers(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Devproxy_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dfuse.devproxy.v1.Devproxy",
	HandlerType: (*DevproxyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListServers",
			Handler:    _Devproxy_ListServers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dfuse/devproxy/v1/devproxy.proto",
}
