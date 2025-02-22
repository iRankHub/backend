// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: internal/grpc/proto/analytics/analytics.proto

package analytics

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AnalyticsService_GetFinancialReports_FullMethodName  = "/analytics.AnalyticsService/GetFinancialReports"
	AnalyticsService_GetAttendanceReports_FullMethodName = "/analytics.AnalyticsService/GetAttendanceReports"
)

// AnalyticsServiceClient is the client API for AnalyticsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyticsServiceClient interface {
	GetFinancialReports(ctx context.Context, in *FinancialReportRequest, opts ...grpc.CallOption) (*FinancialReportResponse, error)
	GetAttendanceReports(ctx context.Context, in *AttendanceReportRequest, opts ...grpc.CallOption) (*AttendanceReportResponse, error)
}

type analyticsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyticsServiceClient(cc grpc.ClientConnInterface) AnalyticsServiceClient {
	return &analyticsServiceClient{cc}
}

func (c *analyticsServiceClient) GetFinancialReports(ctx context.Context, in *FinancialReportRequest, opts ...grpc.CallOption) (*FinancialReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FinancialReportResponse)
	err := c.cc.Invoke(ctx, AnalyticsService_GetFinancialReports_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsServiceClient) GetAttendanceReports(ctx context.Context, in *AttendanceReportRequest, opts ...grpc.CallOption) (*AttendanceReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AttendanceReportResponse)
	err := c.cc.Invoke(ctx, AnalyticsService_GetAttendanceReports_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyticsServiceServer is the server API for AnalyticsService service.
// All implementations must embed UnimplementedAnalyticsServiceServer
// for forward compatibility.
type AnalyticsServiceServer interface {
	GetFinancialReports(context.Context, *FinancialReportRequest) (*FinancialReportResponse, error)
	GetAttendanceReports(context.Context, *AttendanceReportRequest) (*AttendanceReportResponse, error)
	mustEmbedUnimplementedAnalyticsServiceServer()
}

// UnimplementedAnalyticsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAnalyticsServiceServer struct{}

func (UnimplementedAnalyticsServiceServer) GetFinancialReports(context.Context, *FinancialReportRequest) (*FinancialReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFinancialReports not implemented")
}
func (UnimplementedAnalyticsServiceServer) GetAttendanceReports(context.Context, *AttendanceReportRequest) (*AttendanceReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAttendanceReports not implemented")
}
func (UnimplementedAnalyticsServiceServer) mustEmbedUnimplementedAnalyticsServiceServer() {}
func (UnimplementedAnalyticsServiceServer) testEmbeddedByValue()                          {}

// UnsafeAnalyticsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyticsServiceServer will
// result in compilation errors.
type UnsafeAnalyticsServiceServer interface {
	mustEmbedUnimplementedAnalyticsServiceServer()
}

func RegisterAnalyticsServiceServer(s grpc.ServiceRegistrar, srv AnalyticsServiceServer) {
	// If the following call pancis, it indicates UnimplementedAnalyticsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AnalyticsService_ServiceDesc, srv)
}

func _AnalyticsService_GetFinancialReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinancialReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServiceServer).GetFinancialReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsService_GetFinancialReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServiceServer).GetFinancialReports(ctx, req.(*FinancialReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsService_GetAttendanceReports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttendanceReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServiceServer).GetAttendanceReports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsService_GetAttendanceReports_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServiceServer).GetAttendanceReports(ctx, req.(*AttendanceReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AnalyticsService_ServiceDesc is the grpc.ServiceDesc for AnalyticsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AnalyticsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "analytics.AnalyticsService",
	HandlerType: (*AnalyticsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFinancialReports",
			Handler:    _AnalyticsService_GetFinancialReports_Handler,
		},
		{
			MethodName: "GetAttendanceReports",
			Handler:    _AnalyticsService_GetAttendanceReports_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/analytics/analytics.proto",
}
