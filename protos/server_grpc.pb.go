// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: protos/server.proto

package protos

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

// AiRetreatGoClient is the client API for AiRetreatGo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AiRetreatGoClient interface {
	Test(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error)
	AddInstagramAccount(ctx context.Context, in *AddInstagramAccountRequest, opts ...grpc.CallOption) (*AddInstagramAccountResponse, error)
	DeleteInstagramAccount(ctx context.Context, in *DeleteInstagramAccountRequest, opts ...grpc.CallOption) (*DeleteInstagramAccountResponse, error)
}

type socialMineGoClient struct {
	cc grpc.ClientConnInterface
}

func NewAiRetreatGoClient(cc grpc.ClientConnInterface) AiRetreatGoClient {
	return &socialMineGoClient{cc}
}

func (c *socialMineGoClient) Test(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error) {
	out := new(TestResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/Test", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialMineGoClient) AddInstagramAccount(ctx context.Context, in *AddInstagramAccountRequest, opts ...grpc.CallOption) (*AddInstagramAccountResponse, error) {
	out := new(AddInstagramAccountResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/AddInstagramAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *socialMineGoClient) DeleteInstagramAccount(ctx context.Context, in *DeleteInstagramAccountRequest, opts ...grpc.CallOption) (*DeleteInstagramAccountResponse, error) {
	out := new(DeleteInstagramAccountResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/DeleteInstagramAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiRetreatGoServer is the server API for AiRetreatGo service.
// All implementations must embed UnimplementedAiRetreatGoServer
// for forward compatibility
type AiRetreatGoServer interface {
	Test(context.Context, *TestRequest) (*TestResponse, error)
	AddInstagramAccount(context.Context, *AddInstagramAccountRequest) (*AddInstagramAccountResponse, error)
	DeleteInstagramAccount(context.Context, *DeleteInstagramAccountRequest) (*DeleteInstagramAccountResponse, error)
	mustEmbedUnimplementedAiRetreatGoServer()
}

// UnimplementedAiRetreatGoServer must be embedded to have forward compatible implementations.
type UnimplementedAiRetreatGoServer struct {
}

func (UnimplementedAiRetreatGoServer) Test(context.Context, *TestRequest) (*TestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (UnimplementedAiRetreatGoServer) AddInstagramAccount(context.Context, *AddInstagramAccountRequest) (*AddInstagramAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddInstagramAccount not implemented")
}
func (UnimplementedAiRetreatGoServer) DeleteInstagramAccount(context.Context, *DeleteInstagramAccountRequest) (*DeleteInstagramAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteInstagramAccount not implemented")
}
func (UnimplementedAiRetreatGoServer) mustEmbedUnimplementedAiRetreatGoServer() {}

// UnsafeAiRetreatGoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AiRetreatGoServer will
// result in compilation errors.
type UnsafeAiRetreatGoServer interface {
	mustEmbedUnimplementedAiRetreatGoServer()
}

func RegisterAiRetreatGoServer(s grpc.ServiceRegistrar, srv AiRetreatGoServer) {
	s.RegisterService(&AiRetreatGo_ServiceDesc, srv)
}

func _AiRetreatGo_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/Test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).Test(ctx, req.(*TestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_AddInstagramAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddInstagramAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).AddInstagramAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/AddInstagramAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).AddInstagramAccount(ctx, req.(*AddInstagramAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_DeleteInstagramAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInstagramAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).DeleteInstagramAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/DeleteInstagramAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).DeleteInstagramAccount(ctx, req.(*DeleteInstagramAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AiRetreatGo_ServiceDesc is the grpc.ServiceDesc for AiRetreatGo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AiRetreatGo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.AiRetreatGo",
	HandlerType: (*AiRetreatGoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Test",
			Handler:    _AiRetreatGo_Test_Handler,
		},
		{
			MethodName: "AddInstagramAccount",
			Handler:    _AiRetreatGo_AddInstagramAccount_Handler,
		},
		{
			MethodName: "DeleteInstagramAccount",
			Handler:    _AiRetreatGo_DeleteInstagramAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/server.proto",
}

// AiRetreatGoHealthClient is the client API for AiRetreatGoHealth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AiRetreatGoHealthClient interface {
	Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
}

type socialMineGoHealthClient struct {
	cc grpc.ClientConnInterface
}

func NewAiRetreatGoHealthClient(cc grpc.ClientConnInterface) AiRetreatGoHealthClient {
	return &socialMineGoHealthClient{cc}
}

func (c *socialMineGoHealthClient) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGoHealth/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiRetreatGoHealthServer is the server API for AiRetreatGoHealth service.
// All implementations must embed UnimplementedAiRetreatGoHealthServer
// for forward compatibility
type AiRetreatGoHealthServer interface {
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
	mustEmbedUnimplementedAiRetreatGoHealthServer()
}

// UnimplementedAiRetreatGoHealthServer must be embedded to have forward compatible implementations.
type UnimplementedAiRetreatGoHealthServer struct {
}

func (UnimplementedAiRetreatGoHealthServer) Check(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAiRetreatGoHealthServer) mustEmbedUnimplementedAiRetreatGoHealthServer() {}

// UnsafeAiRetreatGoHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AiRetreatGoHealthServer will
// result in compilation errors.
type UnsafeAiRetreatGoHealthServer interface {
	mustEmbedUnimplementedAiRetreatGoHealthServer()
}

func RegisterAiRetreatGoHealthServer(s grpc.ServiceRegistrar, srv AiRetreatGoHealthServer) {
	s.RegisterService(&AiRetreatGoHealth_ServiceDesc, srv)
}

func _AiRetreatGoHealth_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoHealthServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGoHealth/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoHealthServer).Check(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AiRetreatGoHealth_ServiceDesc is the grpc.ServiceDesc for AiRetreatGoHealth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AiRetreatGoHealth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.AiRetreatGoHealth",
	HandlerType: (*AiRetreatGoHealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _AiRetreatGoHealth_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/server.proto",
}
