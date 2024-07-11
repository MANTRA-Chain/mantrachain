// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mantrachain/coinfactory/v1beta1/query.proto

package coinfactoryv1beta1

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

const (
	Query_Params_FullMethodName                  = "/mantrachain.coinfactory.v1beta1.Query/Params"
	Query_DenomsFromCreator_FullMethodName       = "/mantrachain.coinfactory.v1beta1.Query/DenomsFromCreator"
	Query_DenomAuthorityMetadata_FullMethodName  = "/mantrachain.coinfactory.v1beta1.Query/DenomAuthorityMetadata"
	Query_DenomAuthorityMetadata2_FullMethodName = "/mantrachain.coinfactory.v1beta1.Query/DenomAuthorityMetadata2"
	Query_Balance_FullMethodName                 = "/mantrachain.coinfactory.v1beta1.Query/Balance"
	Query_SupplyOf_FullMethodName                = "/mantrachain.coinfactory.v1beta1.Query/SupplyOf"
	Query_DenomMetadata_FullMethodName           = "/mantrachain.coinfactory.v1beta1.Query/DenomMetadata"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// DenomsFromCreator defines a gRPC query method for fetching all
	// denominations created by a specific admin/creator.
	DenomsFromCreator(ctx context.Context, in *QueryDenomsFromCreatorRequest, opts ...grpc.CallOption) (*QueryDenomsFromCreatorResponse, error)
	// DenomAuthorityMetadata defines a gRPC query method for fetching
	// DenomAuthorityMetadata for a particular denom.
	DenomAuthorityMetadata(ctx context.Context, in *QueryDenomAuthorityMetadataRequest, opts ...grpc.CallOption) (*QueryDenomAuthorityMetadataResponse, error)
	// DenomAuthorityMetadata2
	DenomAuthorityMetadata2(ctx context.Context, in *QueryDenomAuthorityMetadata2Request, opts ...grpc.CallOption) (*QueryDenomAuthorityMetadata2Response, error)
	// Balance
	Balance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error)
	// SupplyOf
	SupplyOf(ctx context.Context, in *QuerySupplyOfRequest, opts ...grpc.CallOption) (*QuerySupplyOfResponse, error)
	// DenomMetadata
	DenomMetadata(ctx context.Context, in *QueryDenomMetadataRequest, opts ...grpc.CallOption) (*QueryDenomMetadataResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, Query_Params_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomsFromCreator(ctx context.Context, in *QueryDenomsFromCreatorRequest, opts ...grpc.CallOption) (*QueryDenomsFromCreatorResponse, error) {
	out := new(QueryDenomsFromCreatorResponse)
	err := c.cc.Invoke(ctx, Query_DenomsFromCreator_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomAuthorityMetadata(ctx context.Context, in *QueryDenomAuthorityMetadataRequest, opts ...grpc.CallOption) (*QueryDenomAuthorityMetadataResponse, error) {
	out := new(QueryDenomAuthorityMetadataResponse)
	err := c.cc.Invoke(ctx, Query_DenomAuthorityMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomAuthorityMetadata2(ctx context.Context, in *QueryDenomAuthorityMetadata2Request, opts ...grpc.CallOption) (*QueryDenomAuthorityMetadata2Response, error) {
	out := new(QueryDenomAuthorityMetadata2Response)
	err := c.cc.Invoke(ctx, Query_DenomAuthorityMetadata2_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Balance(ctx context.Context, in *QueryBalanceRequest, opts ...grpc.CallOption) (*QueryBalanceResponse, error) {
	out := new(QueryBalanceResponse)
	err := c.cc.Invoke(ctx, Query_Balance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) SupplyOf(ctx context.Context, in *QuerySupplyOfRequest, opts ...grpc.CallOption) (*QuerySupplyOfResponse, error) {
	out := new(QuerySupplyOfResponse)
	err := c.cc.Invoke(ctx, Query_SupplyOf_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DenomMetadata(ctx context.Context, in *QueryDenomMetadataRequest, opts ...grpc.CallOption) (*QueryDenomMetadataResponse, error) {
	out := new(QueryDenomMetadataResponse)
	err := c.cc.Invoke(ctx, Query_DenomMetadata_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// DenomsFromCreator defines a gRPC query method for fetching all
	// denominations created by a specific admin/creator.
	DenomsFromCreator(context.Context, *QueryDenomsFromCreatorRequest) (*QueryDenomsFromCreatorResponse, error)
	// DenomAuthorityMetadata defines a gRPC query method for fetching
	// DenomAuthorityMetadata for a particular denom.
	DenomAuthorityMetadata(context.Context, *QueryDenomAuthorityMetadataRequest) (*QueryDenomAuthorityMetadataResponse, error)
	// DenomAuthorityMetadata2
	DenomAuthorityMetadata2(context.Context, *QueryDenomAuthorityMetadata2Request) (*QueryDenomAuthorityMetadata2Response, error)
	// Balance
	Balance(context.Context, *QueryBalanceRequest) (*QueryBalanceResponse, error)
	// SupplyOf
	SupplyOf(context.Context, *QuerySupplyOfRequest) (*QuerySupplyOfResponse, error)
	// DenomMetadata
	DenomMetadata(context.Context, *QueryDenomMetadataRequest) (*QueryDenomMetadataResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) DenomsFromCreator(context.Context, *QueryDenomsFromCreatorRequest) (*QueryDenomsFromCreatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomsFromCreator not implemented")
}
func (UnimplementedQueryServer) DenomAuthorityMetadata(context.Context, *QueryDenomAuthorityMetadataRequest) (*QueryDenomAuthorityMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomAuthorityMetadata not implemented")
}
func (UnimplementedQueryServer) DenomAuthorityMetadata2(context.Context, *QueryDenomAuthorityMetadata2Request) (*QueryDenomAuthorityMetadata2Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomAuthorityMetadata2 not implemented")
}
func (UnimplementedQueryServer) Balance(context.Context, *QueryBalanceRequest) (*QueryBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
}
func (UnimplementedQueryServer) SupplyOf(context.Context, *QuerySupplyOfRequest) (*QuerySupplyOfResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SupplyOf not implemented")
}
func (UnimplementedQueryServer) DenomMetadata(context.Context, *QueryDenomMetadataRequest) (*QueryDenomMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DenomMetadata not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
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
		FullMethod: Query_Params_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomsFromCreator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDenomsFromCreatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomsFromCreator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DenomsFromCreator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomsFromCreator(ctx, req.(*QueryDenomsFromCreatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomAuthorityMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDenomAuthorityMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomAuthorityMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DenomAuthorityMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomAuthorityMetadata(ctx, req.(*QueryDenomAuthorityMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomAuthorityMetadata2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDenomAuthorityMetadata2Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomAuthorityMetadata2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DenomAuthorityMetadata2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomAuthorityMetadata2(ctx, req.(*QueryDenomAuthorityMetadata2Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Balance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Balance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Balance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Balance(ctx, req.(*QueryBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_SupplyOf_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySupplyOfRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).SupplyOf(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_SupplyOf_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).SupplyOf(ctx, req.(*QuerySupplyOfRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DenomMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDenomMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DenomMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_DenomMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DenomMetadata(ctx, req.(*QueryDenomMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.coinfactory.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "DenomsFromCreator",
			Handler:    _Query_DenomsFromCreator_Handler,
		},
		{
			MethodName: "DenomAuthorityMetadata",
			Handler:    _Query_DenomAuthorityMetadata_Handler,
		},
		{
			MethodName: "DenomAuthorityMetadata2",
			Handler:    _Query_DenomAuthorityMetadata2_Handler,
		},
		{
			MethodName: "Balance",
			Handler:    _Query_Balance_Handler,
		},
		{
			MethodName: "SupplyOf",
			Handler:    _Query_SupplyOf_Handler,
		},
		{
			MethodName: "DenomMetadata",
			Handler:    _Query_DenomMetadata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/coinfactory/v1beta1/query.proto",
}
