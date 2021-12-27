// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.7.1
// source: grpc_service.proto

package sonyflake

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

// SonyflakeServiceClient is the client API for SonyflakeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SonyflakeServiceClient interface {
	NextID(ctx context.Context, in *SonyFlakeRequest, opts ...grpc.CallOption) (*SonyFlakeResponse, error)
}

type sonyflakeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSonyflakeServiceClient(cc grpc.ClientConnInterface) SonyflakeServiceClient {
	return &sonyflakeServiceClient{cc}
}

func (c *sonyflakeServiceClient) NextID(ctx context.Context, in *SonyFlakeRequest, opts ...grpc.CallOption) (*SonyFlakeResponse, error) {
	out := new(SonyFlakeResponse)
	err := c.cc.Invoke(ctx, "/sonyflake.SonyflakeService/NextID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SonyflakeServiceServer is the server API for SonyflakeService service.
// All implementations must embed UnimplementedSonyflakeServiceServer
// for forward compatibility
type SonyflakeServiceServer interface {
	NextID(context.Context, *SonyFlakeRequest) (*SonyFlakeResponse, error)
	mustEmbedUnimplementedSonyflakeServiceServer()
}

// UnimplementedSonyflakeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSonyflakeServiceServer struct {
}

func (UnimplementedSonyflakeServiceServer) NextID(context.Context, *SonyFlakeRequest) (*SonyFlakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextID not implemented")
}
func (UnimplementedSonyflakeServiceServer) mustEmbedUnimplementedSonyflakeServiceServer() {}

// UnsafeSonyflakeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SonyflakeServiceServer will
// result in compilation errors.
type UnsafeSonyflakeServiceServer interface {
	mustEmbedUnimplementedSonyflakeServiceServer()
}

func RegisterSonyflakeServiceServer(s grpc.ServiceRegistrar, srv SonyflakeServiceServer) {
	s.RegisterService(&SonyflakeService_ServiceDesc, srv)
}

func _SonyflakeService_NextID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SonyFlakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SonyflakeServiceServer).NextID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonyflake.SonyflakeService/NextID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SonyflakeServiceServer).NextID(ctx, req.(*SonyFlakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SonyflakeService_ServiceDesc is the grpc.ServiceDesc for SonyflakeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SonyflakeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sonyflake.SonyflakeService",
	HandlerType: (*SonyflakeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NextID",
			Handler:    _SonyflakeService_NextID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc_service.proto",
}