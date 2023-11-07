// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/txfees/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
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
	return fileDescriptor_7f92659b59977175, []int{0}
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
	return fileDescriptor_7f92659b59977175, []int{1}
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

type QueryGetFeeTokenRequest struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *QueryGetFeeTokenRequest) Reset()         { *m = QueryGetFeeTokenRequest{} }
func (m *QueryGetFeeTokenRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetFeeTokenRequest) ProtoMessage()    {}
func (*QueryGetFeeTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f92659b59977175, []int{2}
}
func (m *QueryGetFeeTokenRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetFeeTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetFeeTokenRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetFeeTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetFeeTokenRequest.Merge(m, src)
}
func (m *QueryGetFeeTokenRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetFeeTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetFeeTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetFeeTokenRequest proto.InternalMessageInfo

func (m *QueryGetFeeTokenRequest) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

type QueryGetFeeTokenResponse struct {
	FeeToken FeeToken `protobuf:"bytes,1,opt,name=feeToken,proto3" json:"feeToken"`
}

func (m *QueryGetFeeTokenResponse) Reset()         { *m = QueryGetFeeTokenResponse{} }
func (m *QueryGetFeeTokenResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetFeeTokenResponse) ProtoMessage()    {}
func (*QueryGetFeeTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f92659b59977175, []int{3}
}
func (m *QueryGetFeeTokenResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetFeeTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetFeeTokenResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetFeeTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetFeeTokenResponse.Merge(m, src)
}
func (m *QueryGetFeeTokenResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetFeeTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetFeeTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetFeeTokenResponse proto.InternalMessageInfo

func (m *QueryGetFeeTokenResponse) GetFeeToken() FeeToken {
	if m != nil {
		return m.FeeToken
	}
	return FeeToken{}
}

type QueryAllFeeTokenRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllFeeTokenRequest) Reset()         { *m = QueryAllFeeTokenRequest{} }
func (m *QueryAllFeeTokenRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllFeeTokenRequest) ProtoMessage()    {}
func (*QueryAllFeeTokenRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f92659b59977175, []int{4}
}
func (m *QueryAllFeeTokenRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllFeeTokenRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllFeeTokenRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllFeeTokenRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllFeeTokenRequest.Merge(m, src)
}
func (m *QueryAllFeeTokenRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllFeeTokenRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllFeeTokenRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllFeeTokenRequest proto.InternalMessageInfo

func (m *QueryAllFeeTokenRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QueryAllFeeTokenResponse struct {
	FeeToken   []FeeToken          `protobuf:"bytes,1,rep,name=feeToken,proto3" json:"feeToken"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllFeeTokenResponse) Reset()         { *m = QueryAllFeeTokenResponse{} }
func (m *QueryAllFeeTokenResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllFeeTokenResponse) ProtoMessage()    {}
func (*QueryAllFeeTokenResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7f92659b59977175, []int{5}
}
func (m *QueryAllFeeTokenResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllFeeTokenResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllFeeTokenResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllFeeTokenResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllFeeTokenResponse.Merge(m, src)
}
func (m *QueryAllFeeTokenResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllFeeTokenResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllFeeTokenResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllFeeTokenResponse proto.InternalMessageInfo

func (m *QueryAllFeeTokenResponse) GetFeeToken() []FeeToken {
	if m != nil {
		return m.FeeToken
	}
	return nil
}

func (m *QueryAllFeeTokenResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "mantrachain.txfees.v1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "mantrachain.txfees.v1.QueryParamsResponse")
	proto.RegisterType((*QueryGetFeeTokenRequest)(nil), "mantrachain.txfees.v1.QueryGetFeeTokenRequest")
	proto.RegisterType((*QueryGetFeeTokenResponse)(nil), "mantrachain.txfees.v1.QueryGetFeeTokenResponse")
	proto.RegisterType((*QueryAllFeeTokenRequest)(nil), "mantrachain.txfees.v1.QueryAllFeeTokenRequest")
	proto.RegisterType((*QueryAllFeeTokenResponse)(nil), "mantrachain.txfees.v1.QueryAllFeeTokenResponse")
}

func init() { proto.RegisterFile("mantrachain/txfees/v1/query.proto", fileDescriptor_7f92659b59977175) }

var fileDescriptor_7f92659b59977175 = []byte{
	// 517 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x3f, 0x6f, 0xd3, 0x40,
	0x18, 0xc6, 0xe3, 0x96, 0x46, 0xe5, 0xba, 0x1d, 0x41, 0x54, 0x11, 0x75, 0xe0, 0xa4, 0xf2, 0xa7,
	0x12, 0x77, 0x72, 0x41, 0xaa, 0x10, 0x03, 0x72, 0x87, 0x74, 0xa2, 0x2a, 0x56, 0x27, 0x24, 0x84,
	0x2e, 0xe1, 0x8d, 0x6b, 0x61, 0xdf, 0xb9, 0xb9, 0x4b, 0xd4, 0x0a, 0xb1, 0x30, 0x30, 0x23, 0xf1,
	0x1d, 0x60, 0xe1, 0x83, 0x74, 0xac, 0xc4, 0xc2, 0x84, 0x50, 0xc2, 0xc2, 0xb7, 0x40, 0xb9, 0x3b,
	0x43, 0x42, 0x5c, 0xb7, 0xb0, 0xd9, 0xe7, 0xe7, 0x79, 0xde, 0xdf, 0xeb, 0xf7, 0xb5, 0xd1, 0xcd,
	0x8c, 0x0b, 0xdd, 0xe7, 0xdd, 0x03, 0x9e, 0x08, 0xa6, 0x8f, 0x7a, 0x00, 0x8a, 0x0d, 0x03, 0x76,
	0x38, 0x80, 0xfe, 0x31, 0xcd, 0xfb, 0x52, 0x4b, 0x7c, 0x75, 0x4a, 0x42, 0xad, 0x84, 0x0e, 0x83,
	0x66, 0x23, 0x96, 0xb1, 0x34, 0x0a, 0x36, 0xb9, 0xb2, 0xe2, 0xe6, 0xf5, 0x58, 0xca, 0x38, 0x05,
	0xc6, 0xf3, 0x84, 0x71, 0x21, 0xa4, 0xe6, 0x3a, 0x91, 0x42, 0xb9, 0xa7, 0x1b, 0x5d, 0xa9, 0x32,
	0xa9, 0x58, 0x87, 0x2b, 0xb0, 0x35, 0xd8, 0x30, 0xe8, 0x80, 0xe6, 0x01, 0xcb, 0x79, 0x9c, 0x08,
	0x23, 0x76, 0x5a, 0x52, 0x4e, 0x96, 0xf3, 0x3e, 0xcf, 0x8a, 0xbc, 0xf5, 0x72, 0x4d, 0x0f, 0xe0,
	0x85, 0x96, 0xaf, 0xc0, 0x45, 0x91, 0x06, 0xc2, 0x4f, 0x27, 0xc5, 0xf6, 0x8c, 0x37, 0x82, 0xc3,
	0x01, 0x28, 0x4d, 0x22, 0x74, 0x65, 0xe6, 0x54, 0xe5, 0x52, 0x28, 0xc0, 0x8f, 0x50, 0xdd, 0xd6,
	0x58, 0xf5, 0x6e, 0x78, 0x77, 0x56, 0x36, 0xd7, 0x68, 0x69, 0xff, 0xd4, 0xda, 0xb6, 0x2f, 0x9d,
	0x7c, 0x6b, 0xd5, 0x22, 0x67, 0x21, 0x0c, 0x5d, 0x33, 0x99, 0x3b, 0xa0, 0xdb, 0x00, 0xfb, 0x13,
	0x06, 0x57, 0x0e, 0x37, 0xd0, 0xd2, 0x4b, 0x10, 0x32, 0x33, 0xb1, 0x97, 0x23, 0x7b, 0x43, 0x9e,
	0xa3, 0xd5, 0x79, 0x83, 0x23, 0x09, 0xd1, 0x72, 0xcf, 0x9d, 0x39, 0x96, 0xd6, 0x19, 0x2c, 0x85,
	0xd5, 0xd1, 0xfc, 0xb6, 0x11, 0xee, 0x78, 0xc2, 0x34, 0xfd, 0x9b, 0xa7, 0x8d, 0xd0, 0x9f, 0x77,
	0xee, 0xf2, 0x6f, 0x51, 0x3b, 0x20, 0x3a, 0x19, 0x10, 0xb5, 0x4b, 0xe0, 0x06, 0x44, 0xf7, 0x78,
	0x0c, 0xce, 0x1b, 0x4d, 0x39, 0xc9, 0x47, 0xcf, 0xb5, 0x30, 0x53, 0xa3, 0xb4, 0x85, 0xc5, 0xff,
	0x68, 0x01, 0xef, 0xcc, 0x70, 0x2e, 0x18, 0xce, 0xdb, 0xe7, 0x72, 0xda, 0xfa, 0xd3, 0xa0, 0x9b,
	0x3f, 0x17, 0xd1, 0x92, 0x01, 0xc5, 0xef, 0x3c, 0x54, 0xb7, 0xe3, 0xc3, 0x77, 0xcf, 0xc0, 0x99,
	0xdf, 0x97, 0xe6, 0xc6, 0x45, 0xa4, 0xb6, 0x2e, 0x59, 0x7f, 0xfb, 0xe5, 0xc7, 0x87, 0x85, 0x16,
	0x5e, 0x63, 0x55, 0x5b, 0x8c, 0x3f, 0x7b, 0x68, 0xb9, 0x68, 0x1c, 0xd3, 0xaa, 0xfc, 0xf9, 0x85,
	0x6a, 0xb2, 0x0b, 0xeb, 0x1d, 0xd4, 0x63, 0x03, 0xf5, 0x10, 0x6f, 0xb1, 0x27, 0xe1, 0xee, 0x7e,
	0x14, 0xde, 0x6b, 0x27, 0x82, 0x8b, 0x2e, 0xb0, 0x73, 0xbe, 0x22, 0xf6, 0xda, 0xec, 0xea, 0x1b,
	0xfc, 0xc9, 0x43, 0x2b, 0x45, 0x6a, 0x98, 0xa6, 0xd5, 0xc4, 0xf3, 0x2b, 0x57, 0x4d, 0x5c, 0xb2,
	0x3e, 0x64, 0xcb, 0x10, 0x07, 0x98, 0xfd, 0x23, 0xf1, 0xf6, 0xee, 0xc9, 0xc8, 0xf7, 0x4e, 0x47,
	0xbe, 0xf7, 0x7d, 0xe4, 0x7b, 0xef, 0xc7, 0x7e, 0xed, 0x74, 0xec, 0xd7, 0xbe, 0x8e, 0xfd, 0xda,
	0xb3, 0x07, 0x71, 0xa2, 0x0f, 0x06, 0x1d, 0xda, 0x95, 0x59, 0x55, 0xe8, 0x51, 0x11, 0xab, 0x8f,
	0x73, 0x50, 0x9d, 0xba, 0xf9, 0x91, 0xdc, 0xff, 0x15, 0x00, 0x00, 0xff, 0xff, 0x67, 0xf8, 0x0e,
	0xbd, 0x2f, 0x05, 0x00, 0x00,
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
	// Queries a list of FeeToken items.
	FeeToken(ctx context.Context, in *QueryGetFeeTokenRequest, opts ...grpc.CallOption) (*QueryGetFeeTokenResponse, error)
	FeeTokenAll(ctx context.Context, in *QueryAllFeeTokenRequest, opts ...grpc.CallOption) (*QueryAllFeeTokenResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.txfees.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FeeToken(ctx context.Context, in *QueryGetFeeTokenRequest, opts ...grpc.CallOption) (*QueryGetFeeTokenResponse, error) {
	out := new(QueryGetFeeTokenResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.txfees.v1.Query/FeeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FeeTokenAll(ctx context.Context, in *QueryAllFeeTokenRequest, opts ...grpc.CallOption) (*QueryAllFeeTokenResponse, error) {
	out := new(QueryAllFeeTokenResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.txfees.v1.Query/FeeTokenAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of FeeToken items.
	FeeToken(context.Context, *QueryGetFeeTokenRequest) (*QueryGetFeeTokenResponse, error)
	FeeTokenAll(context.Context, *QueryAllFeeTokenRequest) (*QueryAllFeeTokenResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) FeeToken(ctx context.Context, req *QueryGetFeeTokenRequest) (*QueryGetFeeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeeToken not implemented")
}
func (*UnimplementedQueryServer) FeeTokenAll(ctx context.Context, req *QueryAllFeeTokenRequest) (*QueryAllFeeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeeTokenAll not implemented")
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
		FullMethod: "/mantrachain.txfees.v1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FeeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetFeeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FeeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.txfees.v1.Query/FeeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FeeToken(ctx, req.(*QueryGetFeeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FeeTokenAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllFeeTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FeeTokenAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.txfees.v1.Query/FeeTokenAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FeeTokenAll(ctx, req.(*QueryAllFeeTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.txfees.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "FeeToken",
			Handler:    _Query_FeeToken_Handler,
		},
		{
			MethodName: "FeeTokenAll",
			Handler:    _Query_FeeTokenAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/txfees/v1/query.proto",
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

func (m *QueryGetFeeTokenRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetFeeTokenRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetFeeTokenRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetFeeTokenResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetFeeTokenResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetFeeTokenResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.FeeToken.MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAllFeeTokenRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllFeeTokenRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllFeeTokenRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *QueryAllFeeTokenResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllFeeTokenResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllFeeTokenResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.FeeToken) > 0 {
		for iNdEx := len(m.FeeToken) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FeeToken[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryGetFeeTokenRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetFeeTokenResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.FeeToken.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllFeeTokenRequest) Size() (n int) {
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

func (m *QueryAllFeeTokenResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.FeeToken) > 0 {
		for _, e := range m.FeeToken {
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
func (m *QueryGetFeeTokenRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetFeeTokenRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetFeeTokenRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
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
func (m *QueryGetFeeTokenResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetFeeTokenResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetFeeTokenResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeToken", wireType)
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
			if err := m.FeeToken.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllFeeTokenRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllFeeTokenRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllFeeTokenRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryAllFeeTokenResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllFeeTokenResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllFeeTokenResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeToken", wireType)
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
			m.FeeToken = append(m.FeeToken, FeeToken{})
			if err := m.FeeToken[len(m.FeeToken)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
