// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mantrachain/lpfarm/v1beta1/query.proto

package lpfarmv1beta1

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
	Query_Params_FullMethodName                 = "/mantrachain.lpfarm.v1beta1.Query/Params"
	Query_QueryPlans_FullMethodName             = "/mantrachain.lpfarm.v1beta1.Query/QueryPlans"
	Query_QueryPlan_FullMethodName              = "/mantrachain.lpfarm.v1beta1.Query/QueryPlan"
	Query_QueryFarm_FullMethodName              = "/mantrachain.lpfarm.v1beta1.Query/QueryFarm"
	Query_QueryPositions_FullMethodName         = "/mantrachain.lpfarm.v1beta1.Query/QueryPositions"
	Query_QueryPosition_FullMethodName          = "/mantrachain.lpfarm.v1beta1.Query/QueryPosition"
	Query_QueryHistoricalRewards_FullMethodName = "/mantrachain.lpfarm.v1beta1.Query/QueryHistoricalRewards"
	Query_QueryTotalRewards_FullMethodName      = "/mantrachain.lpfarm.v1beta1.Query/QueryTotalRewards"
	Query_QueryRewards_FullMethodName           = "/mantrachain.lpfarm.v1beta1.Query/QueryRewards"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	QueryPlans(ctx context.Context, in *QueryPlansRequest, opts ...grpc.CallOption) (*QueryPlansResponse, error)
	QueryPlan(ctx context.Context, in *QueryPlanRequest, opts ...grpc.CallOption) (*QueryPlanResponse, error)
	QueryFarm(ctx context.Context, in *QueryFarmRequest, opts ...grpc.CallOption) (*QueryFarmResponse, error)
	QueryPositions(ctx context.Context, in *QueryPositionsRequest, opts ...grpc.CallOption) (*QueryPositionsResponse, error)
	QueryPosition(ctx context.Context, in *QueryPositionRequest, opts ...grpc.CallOption) (*QueryPositionResponse, error)
	QueryHistoricalRewards(ctx context.Context, in *QueryHistoricalRewardsRequest, opts ...grpc.CallOption) (*QueryHistoricalRewardsResponse, error)
	QueryTotalRewards(ctx context.Context, in *QueryTotalRewardsRequest, opts ...grpc.CallOption) (*QueryTotalRewardsResponse, error)
	QueryRewards(ctx context.Context, in *QueryRewardsRequest, opts ...grpc.CallOption) (*QueryRewardsResponse, error)
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

func (c *queryClient) QueryPlans(ctx context.Context, in *QueryPlansRequest, opts ...grpc.CallOption) (*QueryPlansResponse, error) {
	out := new(QueryPlansResponse)
	err := c.cc.Invoke(ctx, Query_QueryPlans_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryPlan(ctx context.Context, in *QueryPlanRequest, opts ...grpc.CallOption) (*QueryPlanResponse, error) {
	out := new(QueryPlanResponse)
	err := c.cc.Invoke(ctx, Query_QueryPlan_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryFarm(ctx context.Context, in *QueryFarmRequest, opts ...grpc.CallOption) (*QueryFarmResponse, error) {
	out := new(QueryFarmResponse)
	err := c.cc.Invoke(ctx, Query_QueryFarm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryPositions(ctx context.Context, in *QueryPositionsRequest, opts ...grpc.CallOption) (*QueryPositionsResponse, error) {
	out := new(QueryPositionsResponse)
	err := c.cc.Invoke(ctx, Query_QueryPositions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryPosition(ctx context.Context, in *QueryPositionRequest, opts ...grpc.CallOption) (*QueryPositionResponse, error) {
	out := new(QueryPositionResponse)
	err := c.cc.Invoke(ctx, Query_QueryPosition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryHistoricalRewards(ctx context.Context, in *QueryHistoricalRewardsRequest, opts ...grpc.CallOption) (*QueryHistoricalRewardsResponse, error) {
	out := new(QueryHistoricalRewardsResponse)
	err := c.cc.Invoke(ctx, Query_QueryHistoricalRewards_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryTotalRewards(ctx context.Context, in *QueryTotalRewardsRequest, opts ...grpc.CallOption) (*QueryTotalRewardsResponse, error) {
	out := new(QueryTotalRewardsResponse)
	err := c.cc.Invoke(ctx, Query_QueryTotalRewards_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) QueryRewards(ctx context.Context, in *QueryRewardsRequest, opts ...grpc.CallOption) (*QueryRewardsResponse, error) {
	out := new(QueryRewardsResponse)
	err := c.cc.Invoke(ctx, Query_QueryRewards_FullMethodName, in, out, opts...)
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
	QueryPlans(context.Context, *QueryPlansRequest) (*QueryPlansResponse, error)
	QueryPlan(context.Context, *QueryPlanRequest) (*QueryPlanResponse, error)
	QueryFarm(context.Context, *QueryFarmRequest) (*QueryFarmResponse, error)
	QueryPositions(context.Context, *QueryPositionsRequest) (*QueryPositionsResponse, error)
	QueryPosition(context.Context, *QueryPositionRequest) (*QueryPositionResponse, error)
	QueryHistoricalRewards(context.Context, *QueryHistoricalRewardsRequest) (*QueryHistoricalRewardsResponse, error)
	QueryTotalRewards(context.Context, *QueryTotalRewardsRequest) (*QueryTotalRewardsResponse, error)
	QueryRewards(context.Context, *QueryRewardsRequest) (*QueryRewardsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) QueryPlans(context.Context, *QueryPlansRequest) (*QueryPlansResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPlans not implemented")
}
func (UnimplementedQueryServer) QueryPlan(context.Context, *QueryPlanRequest) (*QueryPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPlan not implemented")
}
func (UnimplementedQueryServer) QueryFarm(context.Context, *QueryFarmRequest) (*QueryFarmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryFarm not implemented")
}
func (UnimplementedQueryServer) QueryPositions(context.Context, *QueryPositionsRequest) (*QueryPositionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPositions not implemented")
}
func (UnimplementedQueryServer) QueryPosition(context.Context, *QueryPositionRequest) (*QueryPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryPosition not implemented")
}
func (UnimplementedQueryServer) QueryHistoricalRewards(context.Context, *QueryHistoricalRewardsRequest) (*QueryHistoricalRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryHistoricalRewards not implemented")
}
func (UnimplementedQueryServer) QueryTotalRewards(context.Context, *QueryTotalRewardsRequest) (*QueryTotalRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryTotalRewards not implemented")
}
func (UnimplementedQueryServer) QueryRewards(context.Context, *QueryRewardsRequest) (*QueryRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRewards not implemented")
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

func _Query_QueryPlans_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPlansRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryPlans(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryPlans_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryPlans(ctx, req.(*QueryPlansRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryPlan_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryPlan(ctx, req.(*QueryPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryFarm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryFarmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryFarm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryFarm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryFarm(ctx, req.(*QueryFarmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryPositions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPositionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryPositions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryPositions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryPositions(ctx, req.(*QueryPositionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryPosition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryPosition(ctx, req.(*QueryPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryHistoricalRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryHistoricalRewardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryHistoricalRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryHistoricalRewards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryHistoricalRewards(ctx, req.(*QueryHistoricalRewardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryTotalRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTotalRewardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryTotalRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryTotalRewards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryTotalRewards(ctx, req.(*QueryTotalRewardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_QueryRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRewardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).QueryRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_QueryRewards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).QueryRewards(ctx, req.(*QueryRewardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.lpfarm.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "QueryPlans",
			Handler:    _Query_QueryPlans_Handler,
		},
		{
			MethodName: "QueryPlan",
			Handler:    _Query_QueryPlan_Handler,
		},
		{
			MethodName: "QueryFarm",
			Handler:    _Query_QueryFarm_Handler,
		},
		{
			MethodName: "QueryPositions",
			Handler:    _Query_QueryPositions_Handler,
		},
		{
			MethodName: "QueryPosition",
			Handler:    _Query_QueryPosition_Handler,
		},
		{
			MethodName: "QueryHistoricalRewards",
			Handler:    _Query_QueryHistoricalRewards_Handler,
		},
		{
			MethodName: "QueryTotalRewards",
			Handler:    _Query_QueryTotalRewards_Handler,
		},
		{
			MethodName: "QueryRewards",
			Handler:    _Query_QueryRewards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/lpfarm/v1beta1/query.proto",
}