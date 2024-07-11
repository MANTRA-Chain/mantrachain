// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/bridge/v1beta1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params holds all the parameters of this module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// QueryGetBridgedRequest
type QueryGetBridgedRequest struct {
	EthTxHash string `protobuf:"bytes,1,opt,name=ethTxHash,proto3" json:"ethTxHash,omitempty"`
}

func (m *QueryGetBridgedRequest) Reset()         { *m = QueryGetBridgedRequest{} }
func (m *QueryGetBridgedRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetBridgedRequest) ProtoMessage()    {}
func (*QueryGetBridgedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{2}
}
func (m *QueryGetBridgedRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetBridgedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetBridgedRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetBridgedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetBridgedRequest.Merge(m, src)
}
func (m *QueryGetBridgedRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetBridgedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetBridgedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetBridgedRequest proto.InternalMessageInfo

func (m *QueryGetBridgedRequest) GetEthTxHash() string {
	if m != nil {
		return m.EthTxHash
	}
	return ""
}

// QueryGetBridgedResponse
type QueryGetBridgedResponse struct {
	Bridged Bridged `protobuf:"bytes,1,opt,name=bridged,proto3" json:"bridged"`
}

func (m *QueryGetBridgedResponse) Reset()         { *m = QueryGetBridgedResponse{} }
func (m *QueryGetBridgedResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetBridgedResponse) ProtoMessage()    {}
func (*QueryGetBridgedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{3}
}
func (m *QueryGetBridgedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetBridgedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetBridgedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetBridgedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetBridgedResponse.Merge(m, src)
}
func (m *QueryGetBridgedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetBridgedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetBridgedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetBridgedResponse proto.InternalMessageInfo

func (m *QueryGetBridgedResponse) GetBridged() Bridged {
	if m != nil {
		return m.Bridged
	}
	return Bridged{}
}

// QueryAllBridgedRequest
type QueryAllBridgedRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllBridgedRequest) Reset()         { *m = QueryAllBridgedRequest{} }
func (m *QueryAllBridgedRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllBridgedRequest) ProtoMessage()    {}
func (*QueryAllBridgedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{4}
}
func (m *QueryAllBridgedRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllBridgedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllBridgedRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllBridgedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllBridgedRequest.Merge(m, src)
}
func (m *QueryAllBridgedRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllBridgedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllBridgedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllBridgedRequest proto.InternalMessageInfo

func (m *QueryAllBridgedRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAllBridgedResponse
type QueryAllBridgedResponse struct {
	Bridged    []Bridged           `protobuf:"bytes,1,rep,name=bridged,proto3" json:"bridged"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllBridgedResponse) Reset()         { *m = QueryAllBridgedResponse{} }
func (m *QueryAllBridgedResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllBridgedResponse) ProtoMessage()    {}
func (*QueryAllBridgedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7bef850b17f5ee6f, []int{5}
}
func (m *QueryAllBridgedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllBridgedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllBridgedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllBridgedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllBridgedResponse.Merge(m, src)
}
func (m *QueryAllBridgedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllBridgedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllBridgedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllBridgedResponse proto.InternalMessageInfo

func (m *QueryAllBridgedResponse) GetBridged() []Bridged {
	if m != nil {
		return m.Bridged
	}
	return nil
}

func (m *QueryAllBridgedResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "mantrachain.bridge.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "mantrachain.bridge.v1beta1.QueryParamsResponse")
	proto.RegisterType((*QueryGetBridgedRequest)(nil), "mantrachain.bridge.v1beta1.QueryGetBridgedRequest")
	proto.RegisterType((*QueryGetBridgedResponse)(nil), "mantrachain.bridge.v1beta1.QueryGetBridgedResponse")
	proto.RegisterType((*QueryAllBridgedRequest)(nil), "mantrachain.bridge.v1beta1.QueryAllBridgedRequest")
	proto.RegisterType((*QueryAllBridgedResponse)(nil), "mantrachain.bridge.v1beta1.QueryAllBridgedResponse")
}

func init() {
	proto.RegisterFile("mantrachain/bridge/v1beta1/query.proto", fileDescriptor_7bef850b17f5ee6f)
}

var fileDescriptor_7bef850b17f5ee6f = []byte{
	// 529 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xc1, 0x6b, 0x13, 0x4f,
	0x14, 0xc7, 0x33, 0xed, 0xef, 0x97, 0x92, 0xf1, 0xe4, 0x58, 0x54, 0x42, 0x59, 0xcb, 0x28, 0x69,
	0x09, 0x3a, 0x43, 0x53, 0xf1, 0x9e, 0x88, 0xad, 0x17, 0x8b, 0x2e, 0x3d, 0x89, 0x88, 0xb3, 0xe9,
	0xb0, 0x59, 0xd8, 0x9d, 0xd9, 0xee, 0x4e, 0xa4, 0x45, 0xbc, 0xf8, 0x17, 0x08, 0xe2, 0xd1, 0xab,
	0x7a, 0xf4, 0xcf, 0xe8, 0xb1, 0xe0, 0xc5, 0x93, 0x48, 0x22, 0xe8, 0x9f, 0x21, 0x99, 0x79, 0x69,
	0x12, 0xb7, 0x6e, 0xa2, 0x97, 0x30, 0x79, 0xfb, 0xfd, 0xbe, 0xf7, 0x79, 0xf3, 0xde, 0x2e, 0x6e,
	0x24, 0x42, 0x99, 0x4c, 0x74, 0x7b, 0x22, 0x52, 0x3c, 0xc8, 0xa2, 0x83, 0x50, 0xf2, 0xe7, 0x5b,
	0x81, 0x34, 0x62, 0x8b, 0x1f, 0xf6, 0x65, 0x76, 0xcc, 0xd2, 0x4c, 0x1b, 0x4d, 0xea, 0x53, 0x3a,
	0xe6, 0x74, 0x0c, 0x74, 0xf5, 0x8b, 0x22, 0x89, 0x94, 0xe6, 0xf6, 0xd7, 0xc9, 0xeb, 0xcd, 0xae,
	0xce, 0x13, 0x9d, 0xf3, 0x40, 0xe4, 0xd2, 0xe5, 0x39, 0xcb, 0x9a, 0x8a, 0x30, 0x52, 0xc2, 0x44,
	0x5a, 0x81, 0x76, 0x35, 0xd4, 0xa1, 0xb6, 0x47, 0x3e, 0x3a, 0x41, 0x74, 0x2d, 0xd4, 0x3a, 0x8c,
	0x25, 0x17, 0x69, 0xc4, 0x85, 0x52, 0xda, 0x58, 0x4b, 0x0e, 0x4f, 0x37, 0x4b, 0xb0, 0xdd, 0xdf,
	0x03, 0x50, 0x6e, 0x94, 0x28, 0x53, 0x91, 0x89, 0x04, 0x52, 0xd2, 0x55, 0x4c, 0x1e, 0x8d, 0x40,
	0x1f, 0xda, 0xa0, 0x2f, 0x0f, 0xfb, 0x32, 0x37, 0xf4, 0x09, 0xbe, 0x34, 0x13, 0xcd, 0x53, 0xad,
	0x72, 0x49, 0xee, 0xe1, 0xaa, 0x33, 0x5f, 0x45, 0xeb, 0x68, 0xf3, 0x42, 0x8b, 0xb2, 0x3f, 0xdf,
	0x0f, 0x73, 0xde, 0x4e, 0xed, 0xe4, 0xeb, 0xb5, 0xca, 0xc7, 0x1f, 0x9f, 0x9a, 0xc8, 0x07, 0x33,
	0xbd, 0x83, 0x2f, 0xdb, 0xec, 0xbb, 0xd2, 0x74, 0x1c, 0x35, 0xd4, 0x25, 0x6b, 0xb8, 0x26, 0x4d,
	0x6f, 0xff, 0xe8, 0xbe, 0xc8, 0x7b, 0xb6, 0x46, 0xcd, 0x9f, 0x04, 0xe8, 0x53, 0x7c, 0xa5, 0xe0,
	0x03, 0xb2, 0xbb, 0x78, 0x05, 0x2e, 0x00, 0xd0, 0xae, 0x97, 0xa1, 0x81, 0xbb, 0xf3, 0xdf, 0x88,
	0xcd, 0x1f, 0x3b, 0xe9, 0x33, 0xe0, 0x6a, 0xc7, 0xf1, 0x6f, 0x5c, 0x3b, 0x18, 0x4f, 0x06, 0x08,
	0x15, 0x1a, 0xcc, 0x4d, 0x9b, 0x8d, 0xa6, 0xcd, 0xdc, 0xd6, 0x4c, 0x7a, 0x0f, 0x25, 0x78, 0xfd,
	0x29, 0x27, 0xfd, 0x80, 0xa0, 0x85, 0xe9, 0x12, 0xe7, 0xb5, 0xb0, 0xfc, 0x6f, 0x2d, 0x90, 0xdd,
	0x19, 0xd0, 0x25, 0x0b, 0xba, 0x31, 0x17, 0xd4, 0x11, 0x4c, 0x93, 0xb6, 0x7e, 0x2e, 0xe3, 0xff,
	0x2d, 0x29, 0x79, 0x8b, 0x70, 0xd5, 0xcd, 0x92, 0xb0, 0x32, 0xa2, 0xe2, 0x1a, 0xd5, 0xf9, 0xc2,
	0x7a, 0x47, 0x40, 0x9b, 0xaf, 0x3e, 0x7f, 0x7f, 0xb3, 0x74, 0x83, 0x50, 0x3e, 0x77, 0x7f, 0xc9,
	0x7b, 0x84, 0x57, 0xe0, 0x16, 0x48, 0x6b, 0x6e, 0xa1, 0xc2, 0xae, 0xd5, 0xb7, 0xff, 0xca, 0x03,
	0x80, 0x2d, 0x0b, 0x78, 0x93, 0x34, 0xcf, 0x07, 0x1c, 0xbf, 0x85, 0xfc, 0xc5, 0xd9, 0xd6, 0xbe,
	0x24, 0xef, 0x10, 0xc6, 0x90, 0xa7, 0x1d, 0xc7, 0x0b, 0xb0, 0x16, 0xf6, 0x6f, 0x01, 0xd6, 0xe2,
	0x42, 0xd1, 0x86, 0x65, 0x5d, 0x27, 0x5e, 0x39, 0x6b, 0x67, 0xef, 0x64, 0xe0, 0xa1, 0xd3, 0x81,
	0x87, 0xbe, 0x0d, 0x3c, 0xf4, 0x7a, 0xe8, 0x55, 0x4e, 0x87, 0x5e, 0xe5, 0xcb, 0xd0, 0xab, 0x3c,
	0xbe, 0x1d, 0x46, 0xa6, 0xd7, 0x0f, 0x58, 0x57, 0x27, 0xfc, 0x41, 0x7b, 0x6f, 0xdf, 0x6f, 0xdf,
	0xda, 0x89, 0x94, 0x50, 0x5d, 0x39, 0x93, 0xf2, 0x68, 0x9c, 0xd4, 0x1c, 0xa7, 0x32, 0x0f, 0xaa,
	0xf6, 0xcb, 0xb2, 0xfd, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x28, 0x40, 0xba, 0x74, 0x65, 0x05, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of Bridged items.
	Bridged(ctx context.Context, in *QueryGetBridgedRequest, opts ...grpc.CallOption) (*QueryGetBridgedResponse, error)
	// BridgedAll
	BridgedAll(ctx context.Context, in *QueryAllBridgedRequest, opts ...grpc.CallOption) (*QueryAllBridgedResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.bridge.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Bridged(ctx context.Context, in *QueryGetBridgedRequest, opts ...grpc.CallOption) (*QueryGetBridgedResponse, error) {
	out := new(QueryGetBridgedResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.bridge.v1beta1.Query/Bridged", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BridgedAll(ctx context.Context, in *QueryAllBridgedRequest, opts ...grpc.CallOption) (*QueryAllBridgedResponse, error) {
	out := new(QueryAllBridgedResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.bridge.v1beta1.Query/BridgedAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Bridged items.
	Bridged(context.Context, *QueryGetBridgedRequest) (*QueryGetBridgedResponse, error)
	// BridgedAll
	BridgedAll(context.Context, *QueryAllBridgedRequest) (*QueryAllBridgedResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) Bridged(ctx context.Context, req *QueryGetBridgedRequest) (*QueryGetBridgedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bridged not implemented")
}
func (*UnimplementedQueryServer) BridgedAll(ctx context.Context, req *QueryAllBridgedRequest) (*QueryAllBridgedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BridgedAll not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.bridge.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Bridged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetBridgedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Bridged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.bridge.v1beta1.Query/Bridged",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Bridged(ctx, req.(*QueryGetBridgedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BridgedAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllBridgedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BridgedAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.bridge.v1beta1.Query/BridgedAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BridgedAll(ctx, req.(*QueryAllBridgedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.bridge.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Bridged",
			Handler:    _Query_Bridged_Handler,
		},
		{
			MethodName: "BridgedAll",
			Handler:    _Query_BridgedAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/bridge/v1beta1/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryGetBridgedRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetBridgedRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetBridgedRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EthTxHash) > 0 {
		i -= len(m.EthTxHash)
		copy(dAtA[i:], m.EthTxHash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.EthTxHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetBridgedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetBridgedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetBridgedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Bridged.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryAllBridgedRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllBridgedRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllBridgedRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllBridgedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllBridgedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllBridgedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Bridged) > 0 {
		for iNdEx := len(m.Bridged) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Bridged[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryGetBridgedRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EthTxHash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetBridgedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Bridged.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllBridgedRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllBridgedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Bridged) > 0 {
		for _, e := range m.Bridged {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetBridgedRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetBridgedRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetBridgedRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthTxHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthTxHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetBridgedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetBridgedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetBridgedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bridged", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Bridged.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllBridgedRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllBridgedRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllBridgedRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllBridgedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllBridgedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllBridgedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bridged", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bridged = append(m.Bridged, Bridged{})
			if err := m.Bridged[len(m.Bridged)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
