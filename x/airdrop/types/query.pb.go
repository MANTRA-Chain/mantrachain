// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mantrachain/airdrop/v1beta1/query.proto

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
	return fileDescriptor_6afeeb0a7c7b18c9, []int{0}
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
	return fileDescriptor_6afeeb0a7c7b18c9, []int{1}
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

// QueryGetCampaignRequest
type QueryGetCampaignRequest struct {
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryGetCampaignRequest) Reset()         { *m = QueryGetCampaignRequest{} }
func (m *QueryGetCampaignRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetCampaignRequest) ProtoMessage()    {}
func (*QueryGetCampaignRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afeeb0a7c7b18c9, []int{2}
}
func (m *QueryGetCampaignRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetCampaignRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetCampaignRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetCampaignRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetCampaignRequest.Merge(m, src)
}
func (m *QueryGetCampaignRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetCampaignRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetCampaignRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetCampaignRequest proto.InternalMessageInfo

func (m *QueryGetCampaignRequest) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

// QueryGetCampaignResponse
type QueryGetCampaignResponse struct {
	Campaign Campaign `protobuf:"bytes,1,opt,name=campaign,proto3" json:"campaign"`
}

func (m *QueryGetCampaignResponse) Reset()         { *m = QueryGetCampaignResponse{} }
func (m *QueryGetCampaignResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetCampaignResponse) ProtoMessage()    {}
func (*QueryGetCampaignResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afeeb0a7c7b18c9, []int{3}
}
func (m *QueryGetCampaignResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetCampaignResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetCampaignResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetCampaignResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetCampaignResponse.Merge(m, src)
}
func (m *QueryGetCampaignResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetCampaignResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetCampaignResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetCampaignResponse proto.InternalMessageInfo

func (m *QueryGetCampaignResponse) GetCampaign() Campaign {
	if m != nil {
		return m.Campaign
	}
	return Campaign{}
}

// QueryAllCampaignRequest
type QueryAllCampaignRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllCampaignRequest) Reset()         { *m = QueryAllCampaignRequest{} }
func (m *QueryAllCampaignRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllCampaignRequest) ProtoMessage()    {}
func (*QueryAllCampaignRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afeeb0a7c7b18c9, []int{4}
}
func (m *QueryAllCampaignRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllCampaignRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllCampaignRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllCampaignRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllCampaignRequest.Merge(m, src)
}
func (m *QueryAllCampaignRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllCampaignRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllCampaignRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllCampaignRequest proto.InternalMessageInfo

func (m *QueryAllCampaignRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAllCampaignResponse
type QueryAllCampaignResponse struct {
	Campaign   []Campaign          `protobuf:"bytes,1,rep,name=campaign,proto3" json:"campaign"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllCampaignResponse) Reset()         { *m = QueryAllCampaignResponse{} }
func (m *QueryAllCampaignResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllCampaignResponse) ProtoMessage()    {}
func (*QueryAllCampaignResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6afeeb0a7c7b18c9, []int{5}
}
func (m *QueryAllCampaignResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllCampaignResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllCampaignResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllCampaignResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllCampaignResponse.Merge(m, src)
}
func (m *QueryAllCampaignResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllCampaignResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllCampaignResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllCampaignResponse proto.InternalMessageInfo

func (m *QueryAllCampaignResponse) GetCampaign() []Campaign {
	if m != nil {
		return m.Campaign
	}
	return nil
}

func (m *QueryAllCampaignResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "mantrachain.airdrop.v1beta1.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "mantrachain.airdrop.v1beta1.QueryParamsResponse")
	proto.RegisterType((*QueryGetCampaignRequest)(nil), "mantrachain.airdrop.v1beta1.QueryGetCampaignRequest")
	proto.RegisterType((*QueryGetCampaignResponse)(nil), "mantrachain.airdrop.v1beta1.QueryGetCampaignResponse")
	proto.RegisterType((*QueryAllCampaignRequest)(nil), "mantrachain.airdrop.v1beta1.QueryAllCampaignRequest")
	proto.RegisterType((*QueryAllCampaignResponse)(nil), "mantrachain.airdrop.v1beta1.QueryAllCampaignResponse")
}

func init() {
	proto.RegisterFile("mantrachain/airdrop/v1beta1/query.proto", fileDescriptor_6afeeb0a7c7b18c9)
}

var fileDescriptor_6afeeb0a7c7b18c9 = []byte{
	// 527 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x4f, 0x6b, 0x13, 0x41,
	0x18, 0xc6, 0x33, 0x69, 0x0d, 0x75, 0x0a, 0x82, 0x63, 0xc1, 0x12, 0x65, 0x95, 0x2d, 0x31, 0x35,
	0xda, 0x19, 0x1b, 0xed, 0x07, 0x48, 0x84, 0xe4, 0xe4, 0xbf, 0xc5, 0x93, 0xe0, 0x61, 0xb2, 0x19,
	0xb6, 0x03, 0xbb, 0x33, 0xdb, 0xdd, 0x89, 0x58, 0xc4, 0x8b, 0x9f, 0x40, 0xf0, 0xe2, 0x47, 0x50,
	0xf4, 0xe0, 0xc7, 0xe8, 0xb1, 0xe0, 0xc5, 0x93, 0x48, 0x22, 0xf8, 0x25, 0x3c, 0x48, 0x66, 0x66,
	0xd3, 0xc4, 0x0d, 0x9b, 0xd5, 0x4b, 0x18, 0x86, 0xe7, 0x79, 0xde, 0xdf, 0x3b, 0xef, 0x9b, 0x85,
	0xcd, 0x88, 0x0a, 0x95, 0x50, 0xff, 0x90, 0x72, 0x41, 0x28, 0x4f, 0x86, 0x89, 0x8c, 0xc9, 0x8b,
	0xfd, 0x01, 0x53, 0x74, 0x9f, 0x1c, 0x8d, 0x58, 0x72, 0x8c, 0xe3, 0x44, 0x2a, 0x89, 0xae, 0xcc,
	0x09, 0xb1, 0x15, 0x62, 0x2b, 0xac, 0x5f, 0xa4, 0x11, 0x17, 0x92, 0xe8, 0x5f, 0xa3, 0xaf, 0xb7,
	0x7c, 0x99, 0x46, 0x32, 0x25, 0x03, 0x9a, 0x32, 0x13, 0x34, 0x8b, 0x8d, 0x69, 0xc0, 0x05, 0x55,
	0x5c, 0x0a, 0xab, 0xdd, 0x0a, 0x64, 0x20, 0xf5, 0x91, 0x4c, 0x4f, 0xf6, 0xf6, 0x6a, 0x20, 0x65,
	0x10, 0x32, 0x42, 0x63, 0x4e, 0xa8, 0x10, 0x52, 0x69, 0x4b, 0x9a, 0xe5, 0x17, 0x81, 0xfb, 0x34,
	0x8a, 0x29, 0x0f, 0xb2, 0xfc, 0xdd, 0x22, 0x6d, 0x4c, 0x13, 0x1a, 0xd9, 0x54, 0x77, 0x0b, 0xa2,
	0x27, 0x53, 0xd6, 0xc7, 0xfa, 0xd2, 0x63, 0x47, 0x23, 0x96, 0x2a, 0xf7, 0x39, 0xbc, 0xb4, 0x70,
	0x9b, 0xc6, 0x52, 0xa4, 0x0c, 0xf5, 0x60, 0xcd, 0x98, 0xb7, 0xc1, 0x75, 0xb0, 0xbb, 0xd9, 0xde,
	0xc1, 0x05, 0x6f, 0x84, 0x8d, 0xb9, 0x7b, 0xfe, 0xe4, 0xfb, 0xb5, 0xca, 0x87, 0x5f, 0x5f, 0x5a,
	0xc0, 0xb3, 0x6e, 0xf7, 0x26, 0xbc, 0xac, 0xe3, 0xfb, 0x4c, 0xdd, 0xb7, 0xe0, 0xb6, 0x32, 0xba,
	0x00, 0xab, 0x7c, 0xa8, 0xe3, 0xd7, 0xbd, 0x2a, 0x1f, 0xba, 0x3e, 0xdc, 0xce, 0x4b, 0x2d, 0x4e,
	0x1f, 0x6e, 0x64, 0x7d, 0x5b, 0xa0, 0x46, 0x21, 0x50, 0x16, 0xd0, 0x5d, 0x9f, 0x22, 0x79, 0x33,
	0xb3, 0x4b, 0x2d, 0x4f, 0x27, 0x0c, 0xff, 0xe6, 0xe9, 0x41, 0x78, 0x36, 0x3d, 0x5b, 0xe5, 0x06,
	0x36, 0xa3, 0xc6, 0xd3, 0x51, 0x63, 0xb3, 0x33, 0x67, 0x4d, 0x07, 0xcc, 0x7a, 0xbd, 0x39, 0xa7,
	0xfb, 0x19, 0xd8, 0x46, 0x16, 0x6a, 0x2c, 0x6d, 0x64, 0xed, 0xbf, 0x1b, 0x41, 0xfd, 0x05, 0xda,
	0xaa, 0xa6, 0x6d, 0xae, 0xa4, 0x35, 0x14, 0xf3, 0xb8, 0xed, 0xdf, 0x6b, 0xf0, 0x9c, 0xc6, 0x45,
	0xef, 0x01, 0xac, 0x99, 0x49, 0x22, 0x52, 0x08, 0x95, 0x5f, 0xa3, 0xfa, 0x9d, 0xf2, 0x06, 0xc3,
	0xe0, 0xde, 0x7a, 0xf3, 0xf5, 0xe7, 0xbb, 0x6a, 0x03, 0xed, 0x90, 0xd5, 0x1b, 0x8c, 0x3e, 0x01,
	0xb8, 0x91, 0x3d, 0x05, 0xba, 0xb7, 0xba, 0x56, 0x7e, 0xdd, 0xea, 0x07, 0xff, 0xe8, 0xb2, 0x98,
	0x6d, 0x8d, 0x79, 0x1b, 0xb5, 0x48, 0x99, 0x3f, 0x25, 0x79, 0xc5, 0x87, 0xaf, 0xd1, 0x47, 0x00,
	0x37, 0xb3, 0xa0, 0x4e, 0x18, 0x96, 0x01, 0xce, 0xef, 0x63, 0x19, 0xe0, 0x25, 0x1b, 0xe6, 0xee,
	0x69, 0xe0, 0x26, 0x6a, 0x94, 0x02, 0xee, 0x3e, 0x3a, 0x19, 0x3b, 0xe0, 0x74, 0xec, 0x80, 0x1f,
	0x63, 0x07, 0xbc, 0x9d, 0x38, 0x95, 0xd3, 0x89, 0x53, 0xf9, 0x36, 0x71, 0x2a, 0xcf, 0x0e, 0x02,
	0xae, 0x0e, 0x47, 0x03, 0xec, 0xcb, 0x88, 0x3c, 0xe8, 0x3c, 0x7c, 0xea, 0x75, 0xf6, 0x7a, 0x5c,
	0x50, 0xe1, 0xb3, 0x85, 0xe4, 0x97, 0xb3, 0x6c, 0x75, 0x1c, 0xb3, 0x74, 0x50, 0xd3, 0x5f, 0x9b,
	0xbb, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4e, 0xce, 0xc8, 0x00, 0x7e, 0x05, 0x00, 0x00,
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
	// Queries a list of Campaign items.
	Campaign(ctx context.Context, in *QueryGetCampaignRequest, opts ...grpc.CallOption) (*QueryGetCampaignResponse, error)
	// Queries a list of Campaign items.
	CampaignAll(ctx context.Context, in *QueryAllCampaignRequest, opts ...grpc.CallOption) (*QueryAllCampaignResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.airdrop.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Campaign(ctx context.Context, in *QueryGetCampaignRequest, opts ...grpc.CallOption) (*QueryGetCampaignResponse, error) {
	out := new(QueryGetCampaignResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.airdrop.v1beta1.Query/Campaign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CampaignAll(ctx context.Context, in *QueryAllCampaignRequest, opts ...grpc.CallOption) (*QueryAllCampaignResponse, error) {
	out := new(QueryAllCampaignResponse)
	err := c.cc.Invoke(ctx, "/mantrachain.airdrop.v1beta1.Query/CampaignAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Campaign items.
	Campaign(context.Context, *QueryGetCampaignRequest) (*QueryGetCampaignResponse, error)
	// Queries a list of Campaign items.
	CampaignAll(context.Context, *QueryAllCampaignRequest) (*QueryAllCampaignResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) Campaign(ctx context.Context, req *QueryGetCampaignRequest) (*QueryGetCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Campaign not implemented")
}
func (*UnimplementedQueryServer) CampaignAll(ctx context.Context, req *QueryAllCampaignRequest) (*QueryAllCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CampaignAll not implemented")
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
		FullMethod: "/mantrachain.airdrop.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Campaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetCampaignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Campaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.airdrop.v1beta1.Query/Campaign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Campaign(ctx, req.(*QueryGetCampaignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CampaignAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllCampaignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CampaignAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mantrachain.airdrop.v1beta1.Query/CampaignAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CampaignAll(ctx, req.(*QueryAllCampaignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.airdrop.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Campaign",
			Handler:    _Query_Campaign_Handler,
		},
		{
			MethodName: "CampaignAll",
			Handler:    _Query_CampaignAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/airdrop/v1beta1/query.proto",
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

func (m *QueryGetCampaignRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetCampaignRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetCampaignRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetCampaignResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetCampaignResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetCampaignResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Campaign.MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAllCampaignRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllCampaignRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllCampaignRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *QueryAllCampaignResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllCampaignResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllCampaignResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Campaign) > 0 {
		for iNdEx := len(m.Campaign) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Campaign[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryGetCampaignRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovQuery(uint64(m.Id))
	}
	return n
}

func (m *QueryGetCampaignResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Campaign.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllCampaignRequest) Size() (n int) {
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

func (m *QueryAllCampaignResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Campaign) > 0 {
		for _, e := range m.Campaign {
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
func (m *QueryGetCampaignRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetCampaignRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetCampaignRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *QueryGetCampaignResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetCampaignResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetCampaignResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Campaign", wireType)
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
			if err := m.Campaign.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllCampaignRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllCampaignRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllCampaignRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryAllCampaignResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllCampaignResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllCampaignResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Campaign", wireType)
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
			m.Campaign = append(m.Campaign, Campaign{})
			if err := m.Campaign[len(m.Campaign)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
