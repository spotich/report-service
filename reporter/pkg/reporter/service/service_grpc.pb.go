// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.4
// source: api/reporter/service.proto

package service

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

// ReporterClient is the client API for Reporter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReporterClient interface {
	GetReport(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Reporter_GetReportClient, error)
}

type reporterClient struct {
	cc grpc.ClientConnInterface
}

func NewReporterClient(cc grpc.ClientConnInterface) ReporterClient {
	return &reporterClient{cc}
}

func (c *reporterClient) GetReport(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Reporter_GetReportClient, error) {
	stream, err := c.cc.NewStream(ctx, &Reporter_ServiceDesc.Streams[0], "/reporter.Reporter/GetReport", opts...)
	if err != nil {
		return nil, err
	}
	x := &reporterGetReportClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Reporter_GetReportClient interface {
	Recv() (*GetResponse, error)
	grpc.ClientStream
}

type reporterGetReportClient struct {
	grpc.ClientStream
}

func (x *reporterGetReportClient) Recv() (*GetResponse, error) {
	m := new(GetResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReporterServer is the server API for Reporter service.
// All implementations must embed UnimplementedReporterServer
// for forward compatibility
type ReporterServer interface {
	GetReport(*GetRequest, Reporter_GetReportServer) error
	mustEmbedUnimplementedReporterServer()
}

// UnimplementedReporterServer must be embedded to have forward compatible implementations.
type UnimplementedReporterServer struct {
}

func (UnimplementedReporterServer) GetReport(*GetRequest, Reporter_GetReportServer) error {
	return status.Errorf(codes.Unimplemented, "method GetReport not implemented")
}
func (UnimplementedReporterServer) mustEmbedUnimplementedReporterServer() {}

// UnsafeReporterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReporterServer will
// result in compilation errors.
type UnsafeReporterServer interface {
	mustEmbedUnimplementedReporterServer()
}

func RegisterReporterServer(s grpc.ServiceRegistrar, srv ReporterServer) {
	s.RegisterService(&Reporter_ServiceDesc, srv)
}

func _Reporter_GetReport_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReporterServer).GetReport(m, &reporterGetReportServer{stream})
}

type Reporter_GetReportServer interface {
	Send(*GetResponse) error
	grpc.ServerStream
}

type reporterGetReportServer struct {
	grpc.ServerStream
}

func (x *reporterGetReportServer) Send(m *GetResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Reporter_ServiceDesc is the grpc.ServiceDesc for Reporter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Reporter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "reporter.Reporter",
	HandlerType: (*ReporterServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetReport",
			Handler:       _Reporter_GetReport_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/reporter/service.proto",
}
