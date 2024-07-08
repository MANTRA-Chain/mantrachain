// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mantrachain/did/v1/tx.proto

package didv1

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
	Msg_UpdateParams_FullMethodName                 = "/mantrachain.did.v1.Msg/UpdateParams"
	Msg_CreateDidDocument_FullMethodName            = "/mantrachain.did.v1.Msg/CreateDidDocument"
	Msg_UpdateDidDocument_FullMethodName            = "/mantrachain.did.v1.Msg/UpdateDidDocument"
	Msg_AddVerification_FullMethodName              = "/mantrachain.did.v1.Msg/AddVerification"
	Msg_RevokeVerification_FullMethodName           = "/mantrachain.did.v1.Msg/RevokeVerification"
	Msg_SetVerificationRelationships_FullMethodName = "/mantrachain.did.v1.Msg/SetVerificationRelationships"
	Msg_AddService_FullMethodName                   = "/mantrachain.did.v1.Msg/AddService"
	Msg_DeleteService_FullMethodName                = "/mantrachain.did.v1.Msg/DeleteService"
	Msg_AddController_FullMethodName                = "/mantrachain.did.v1.Msg/AddController"
	Msg_DeleteController_FullMethodName             = "/mantrachain.did.v1.Msg/DeleteController"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	// CreateDidDocument defines a method for creating a new identity.
	CreateDidDocument(ctx context.Context, in *MsgCreateDidDocument, opts ...grpc.CallOption) (*MsgCreateDidDocumentResponse, error)
	// UpdateDidDocument defines a method for updating an identity.
	UpdateDidDocument(ctx context.Context, in *MsgUpdateDidDocument, opts ...grpc.CallOption) (*MsgUpdateDidDocumentResponse, error)
	// AddVerificationMethod adds a new verification method
	AddVerification(ctx context.Context, in *MsgAddVerification, opts ...grpc.CallOption) (*MsgAddVerificationResponse, error)
	// RevokeVerification remove the verification method and all associated
	// verification Relations
	RevokeVerification(ctx context.Context, in *MsgRevokeVerification, opts ...grpc.CallOption) (*MsgRevokeVerificationResponse, error)
	// SetVerificationRelationships overwrite current verification relationships
	SetVerificationRelationships(ctx context.Context, in *MsgSetVerificationRelationships, opts ...grpc.CallOption) (*MsgSetVerificationRelationshipsResponse, error)
	// AddService add a new service
	AddService(ctx context.Context, in *MsgAddService, opts ...grpc.CallOption) (*MsgAddServiceResponse, error)
	// DeleteService delete an existing service
	DeleteService(ctx context.Context, in *MsgDeleteService, opts ...grpc.CallOption) (*MsgDeleteServiceResponse, error)
	// AddService add a new service
	AddController(ctx context.Context, in *MsgAddController, opts ...grpc.CallOption) (*MsgAddControllerResponse, error)
	// DeleteService delete an existing service
	DeleteController(ctx context.Context, in *MsgDeleteController, opts ...grpc.CallOption) (*MsgDeleteControllerResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateDidDocument(ctx context.Context, in *MsgCreateDidDocument, opts ...grpc.CallOption) (*MsgCreateDidDocumentResponse, error) {
	out := new(MsgCreateDidDocumentResponse)
	err := c.cc.Invoke(ctx, Msg_CreateDidDocument_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateDidDocument(ctx context.Context, in *MsgUpdateDidDocument, opts ...grpc.CallOption) (*MsgUpdateDidDocumentResponse, error) {
	out := new(MsgUpdateDidDocumentResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateDidDocument_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddVerification(ctx context.Context, in *MsgAddVerification, opts ...grpc.CallOption) (*MsgAddVerificationResponse, error) {
	out := new(MsgAddVerificationResponse)
	err := c.cc.Invoke(ctx, Msg_AddVerification_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RevokeVerification(ctx context.Context, in *MsgRevokeVerification, opts ...grpc.CallOption) (*MsgRevokeVerificationResponse, error) {
	out := new(MsgRevokeVerificationResponse)
	err := c.cc.Invoke(ctx, Msg_RevokeVerification_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetVerificationRelationships(ctx context.Context, in *MsgSetVerificationRelationships, opts ...grpc.CallOption) (*MsgSetVerificationRelationshipsResponse, error) {
	out := new(MsgSetVerificationRelationshipsResponse)
	err := c.cc.Invoke(ctx, Msg_SetVerificationRelationships_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddService(ctx context.Context, in *MsgAddService, opts ...grpc.CallOption) (*MsgAddServiceResponse, error) {
	out := new(MsgAddServiceResponse)
	err := c.cc.Invoke(ctx, Msg_AddService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeleteService(ctx context.Context, in *MsgDeleteService, opts ...grpc.CallOption) (*MsgDeleteServiceResponse, error) {
	out := new(MsgDeleteServiceResponse)
	err := c.cc.Invoke(ctx, Msg_DeleteService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddController(ctx context.Context, in *MsgAddController, opts ...grpc.CallOption) (*MsgAddControllerResponse, error) {
	out := new(MsgAddControllerResponse)
	err := c.cc.Invoke(ctx, Msg_AddController_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeleteController(ctx context.Context, in *MsgDeleteController, opts ...grpc.CallOption) (*MsgDeleteControllerResponse, error) {
	out := new(MsgDeleteControllerResponse)
	err := c.cc.Invoke(ctx, Msg_DeleteController_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	// CreateDidDocument defines a method for creating a new identity.
	CreateDidDocument(context.Context, *MsgCreateDidDocument) (*MsgCreateDidDocumentResponse, error)
	// UpdateDidDocument defines a method for updating an identity.
	UpdateDidDocument(context.Context, *MsgUpdateDidDocument) (*MsgUpdateDidDocumentResponse, error)
	// AddVerificationMethod adds a new verification method
	AddVerification(context.Context, *MsgAddVerification) (*MsgAddVerificationResponse, error)
	// RevokeVerification remove the verification method and all associated
	// verification Relations
	RevokeVerification(context.Context, *MsgRevokeVerification) (*MsgRevokeVerificationResponse, error)
	// SetVerificationRelationships overwrite current verification relationships
	SetVerificationRelationships(context.Context, *MsgSetVerificationRelationships) (*MsgSetVerificationRelationshipsResponse, error)
	// AddService add a new service
	AddService(context.Context, *MsgAddService) (*MsgAddServiceResponse, error)
	// DeleteService delete an existing service
	DeleteService(context.Context, *MsgDeleteService) (*MsgDeleteServiceResponse, error)
	// AddService add a new service
	AddController(context.Context, *MsgAddController) (*MsgAddControllerResponse, error)
	// DeleteService delete an existing service
	DeleteController(context.Context, *MsgDeleteController) (*MsgDeleteControllerResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) CreateDidDocument(context.Context, *MsgCreateDidDocument) (*MsgCreateDidDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDidDocument not implemented")
}
func (UnimplementedMsgServer) UpdateDidDocument(context.Context, *MsgUpdateDidDocument) (*MsgUpdateDidDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDidDocument not implemented")
}
func (UnimplementedMsgServer) AddVerification(context.Context, *MsgAddVerification) (*MsgAddVerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddVerification not implemented")
}
func (UnimplementedMsgServer) RevokeVerification(context.Context, *MsgRevokeVerification) (*MsgRevokeVerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeVerification not implemented")
}
func (UnimplementedMsgServer) SetVerificationRelationships(context.Context, *MsgSetVerificationRelationships) (*MsgSetVerificationRelationshipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetVerificationRelationships not implemented")
}
func (UnimplementedMsgServer) AddService(context.Context, *MsgAddService) (*MsgAddServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddService not implemented")
}
func (UnimplementedMsgServer) DeleteService(context.Context, *MsgDeleteService) (*MsgDeleteServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteService not implemented")
}
func (UnimplementedMsgServer) AddController(context.Context, *MsgAddController) (*MsgAddControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddController not implemented")
}
func (UnimplementedMsgServer) DeleteController(context.Context, *MsgDeleteController) (*MsgDeleteControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteController not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateDidDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateDidDocument)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateDidDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateDidDocument_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateDidDocument(ctx, req.(*MsgCreateDidDocument))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateDidDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateDidDocument)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateDidDocument(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateDidDocument_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateDidDocument(ctx, req.(*MsgUpdateDidDocument))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddVerification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_AddVerification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddVerification(ctx, req.(*MsgAddVerification))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RevokeVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRevokeVerification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RevokeVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RevokeVerification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RevokeVerification(ctx, req.(*MsgRevokeVerification))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetVerificationRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetVerificationRelationships)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetVerificationRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_SetVerificationRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetVerificationRelationships(ctx, req.(*MsgSetVerificationRelationships))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_AddService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddService(ctx, req.(*MsgAddService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeleteService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeleteService)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeleteService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DeleteService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeleteService(ctx, req.(*MsgDeleteService))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddController_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddController)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddController(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_AddController_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddController(ctx, req.(*MsgAddController))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeleteController_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeleteController)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeleteController(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_DeleteController_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeleteController(ctx, req.(*MsgDeleteController))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mantrachain.did.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "CreateDidDocument",
			Handler:    _Msg_CreateDidDocument_Handler,
		},
		{
			MethodName: "UpdateDidDocument",
			Handler:    _Msg_UpdateDidDocument_Handler,
		},
		{
			MethodName: "AddVerification",
			Handler:    _Msg_AddVerification_Handler,
		},
		{
			MethodName: "RevokeVerification",
			Handler:    _Msg_RevokeVerification_Handler,
		},
		{
			MethodName: "SetVerificationRelationships",
			Handler:    _Msg_SetVerificationRelationships_Handler,
		},
		{
			MethodName: "AddService",
			Handler:    _Msg_AddService_Handler,
		},
		{
			MethodName: "DeleteService",
			Handler:    _Msg_DeleteService_Handler,
		},
		{
			MethodName: "AddController",
			Handler:    _Msg_AddController_Handler,
		},
		{
			MethodName: "DeleteController",
			Handler:    _Msg_DeleteController_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mantrachain/did/v1/tx.proto",
}