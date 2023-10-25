// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: api/service.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TimerClient is the client API for Timer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimerClient interface {
	GetTime(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetResponse, error)
}

type timerClient struct {
	cc grpc.ClientConnInterface
}

func NewTimerClient(cc grpc.ClientConnInterface) TimerClient {
	return &timerClient{cc}
}

func (c *timerClient) GetTime(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/timer.Timer/GetTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimerServer is the server API for Timer service.
// All implementations must embed UnimplementedTimerServer
// for forward compatibility
type TimerServer interface {
	GetTime(context.Context, *emptypb.Empty) (*GetResponse, error)
	mustEmbedUnimplementedTimerServer()
}

// UnimplementedTimerServer must be embedded to have forward compatible implementations.
type UnimplementedTimerServer struct {
}

func (UnimplementedTimerServer) GetTime(context.Context, *emptypb.Empty) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTime not implemented")
}
func (UnimplementedTimerServer) mustEmbedUnimplementedTimerServer() {}

// UnsafeTimerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimerServer will
// result in compilation errors.
type UnsafeTimerServer interface {
	mustEmbedUnimplementedTimerServer()
}

func RegisterTimerServer(s grpc.ServiceRegistrar, srv TimerServer) {
	s.RegisterService(&Timer_ServiceDesc, srv)
}

func _Timer_GetTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimerServer).GetTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/timer.Timer/GetTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimerServer).GetTime(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Timer_ServiceDesc is the grpc.ServiceDesc for Timer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Timer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "timer.Timer",
	HandlerType: (*TimerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTime",
			Handler:    _Timer_GetTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/service.proto",
}
