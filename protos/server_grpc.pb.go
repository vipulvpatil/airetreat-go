// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
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
	CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error)
	JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (*JoinGameResponse, error)
	AutoJoinGame(ctx context.Context, in *AutoJoinGameRequest, opts ...grpc.CallOption) (*AutoJoinGameResponse, error)
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	Tag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*TagResponse, error)
	Help(ctx context.Context, in *HelpRequest, opts ...grpc.CallOption) (*HelpResponse, error)
	GetGameForPlayer(ctx context.Context, in *GetGameForPlayerRequest, opts ...grpc.CallOption) (*GetGameForPlayerResponse, error)
	GetGamesForPlayer(ctx context.Context, in *GetGamesForPlayerRequest, opts ...grpc.CallOption) (*GetGamesForPlayerResponse, error)
	SyncPlayerData(ctx context.Context, in *SyncPlayerDataRequest, opts ...grpc.CallOption) (*SyncPlayerDataResponse, error)
}

type aiRetreatGoClient struct {
	cc grpc.ClientConnInterface
}

func NewAiRetreatGoClient(cc grpc.ClientConnInterface) AiRetreatGoClient {
	return &aiRetreatGoClient{cc}
}

func (c *aiRetreatGoClient) CreateGame(ctx context.Context, in *CreateGameRequest, opts ...grpc.CallOption) (*CreateGameResponse, error) {
	out := new(CreateGameResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/CreateGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) JoinGame(ctx context.Context, in *JoinGameRequest, opts ...grpc.CallOption) (*JoinGameResponse, error) {
	out := new(JoinGameResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/JoinGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) AutoJoinGame(ctx context.Context, in *AutoJoinGameRequest, opts ...grpc.CallOption) (*AutoJoinGameResponse, error) {
	out := new(AutoJoinGameResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/AutoJoinGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) Tag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*TagResponse, error) {
	out := new(TagResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/Tag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) Help(ctx context.Context, in *HelpRequest, opts ...grpc.CallOption) (*HelpResponse, error) {
	out := new(HelpResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/Help", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) GetGameForPlayer(ctx context.Context, in *GetGameForPlayerRequest, opts ...grpc.CallOption) (*GetGameForPlayerResponse, error) {
	out := new(GetGameForPlayerResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/GetGameForPlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) GetGamesForPlayer(ctx context.Context, in *GetGamesForPlayerRequest, opts ...grpc.CallOption) (*GetGamesForPlayerResponse, error) {
	out := new(GetGamesForPlayerResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/GetGamesForPlayer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiRetreatGoClient) SyncPlayerData(ctx context.Context, in *SyncPlayerDataRequest, opts ...grpc.CallOption) (*SyncPlayerDataResponse, error) {
	out := new(SyncPlayerDataResponse)
	err := c.cc.Invoke(ctx, "/protos.AiRetreatGo/SyncPlayerData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiRetreatGoServer is the server API for AiRetreatGo service.
// All implementations must embed UnimplementedAiRetreatGoServer
// for forward compatibility
type AiRetreatGoServer interface {
	CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error)
	JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error)
	AutoJoinGame(context.Context, *AutoJoinGameRequest) (*AutoJoinGameResponse, error)
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	Tag(context.Context, *TagRequest) (*TagResponse, error)
	Help(context.Context, *HelpRequest) (*HelpResponse, error)
	GetGameForPlayer(context.Context, *GetGameForPlayerRequest) (*GetGameForPlayerResponse, error)
	GetGamesForPlayer(context.Context, *GetGamesForPlayerRequest) (*GetGamesForPlayerResponse, error)
	SyncPlayerData(context.Context, *SyncPlayerDataRequest) (*SyncPlayerDataResponse, error)
	mustEmbedUnimplementedAiRetreatGoServer()
}

// UnimplementedAiRetreatGoServer must be embedded to have forward compatible implementations.
type UnimplementedAiRetreatGoServer struct {
}

func (UnimplementedAiRetreatGoServer) CreateGame(context.Context, *CreateGameRequest) (*CreateGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGame not implemented")
}
func (UnimplementedAiRetreatGoServer) JoinGame(context.Context, *JoinGameRequest) (*JoinGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGame not implemented")
}
func (UnimplementedAiRetreatGoServer) AutoJoinGame(context.Context, *AutoJoinGameRequest) (*AutoJoinGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutoJoinGame not implemented")
}
func (UnimplementedAiRetreatGoServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedAiRetreatGoServer) Tag(context.Context, *TagRequest) (*TagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tag not implemented")
}
func (UnimplementedAiRetreatGoServer) Help(context.Context, *HelpRequest) (*HelpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Help not implemented")
}
func (UnimplementedAiRetreatGoServer) GetGameForPlayer(context.Context, *GetGameForPlayerRequest) (*GetGameForPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameForPlayer not implemented")
}
func (UnimplementedAiRetreatGoServer) GetGamesForPlayer(context.Context, *GetGamesForPlayerRequest) (*GetGamesForPlayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGamesForPlayer not implemented")
}
func (UnimplementedAiRetreatGoServer) SyncPlayerData(context.Context, *SyncPlayerDataRequest) (*SyncPlayerDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncPlayerData not implemented")
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

func _AiRetreatGo_CreateGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).CreateGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/CreateGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).CreateGame(ctx, req.(*CreateGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_JoinGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).JoinGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/JoinGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).JoinGame(ctx, req.(*JoinGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_AutoJoinGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutoJoinGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).AutoJoinGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/AutoJoinGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).AutoJoinGame(ctx, req.(*AutoJoinGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_Tag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).Tag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/Tag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).Tag(ctx, req.(*TagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_Help_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).Help(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/Help",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).Help(ctx, req.(*HelpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_GetGameForPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGameForPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).GetGameForPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/GetGameForPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).GetGameForPlayer(ctx, req.(*GetGameForPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_GetGamesForPlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGamesForPlayerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).GetGamesForPlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/GetGamesForPlayer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).GetGamesForPlayer(ctx, req.(*GetGamesForPlayerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiRetreatGo_SyncPlayerData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncPlayerDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiRetreatGoServer).SyncPlayerData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.AiRetreatGo/SyncPlayerData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiRetreatGoServer).SyncPlayerData(ctx, req.(*SyncPlayerDataRequest))
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
			MethodName: "CreateGame",
			Handler:    _AiRetreatGo_CreateGame_Handler,
		},
		{
			MethodName: "JoinGame",
			Handler:    _AiRetreatGo_JoinGame_Handler,
		},
		{
			MethodName: "AutoJoinGame",
			Handler:    _AiRetreatGo_AutoJoinGame_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _AiRetreatGo_SendMessage_Handler,
		},
		{
			MethodName: "Tag",
			Handler:    _AiRetreatGo_Tag_Handler,
		},
		{
			MethodName: "Help",
			Handler:    _AiRetreatGo_Help_Handler,
		},
		{
			MethodName: "GetGameForPlayer",
			Handler:    _AiRetreatGo_GetGameForPlayer_Handler,
		},
		{
			MethodName: "GetGamesForPlayer",
			Handler:    _AiRetreatGo_GetGamesForPlayer_Handler,
		},
		{
			MethodName: "SyncPlayerData",
			Handler:    _AiRetreatGo_SyncPlayerData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/server.proto",
}
