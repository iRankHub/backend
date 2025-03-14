// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: internal/grpc/proto/notification/notification.proto

package notification

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
	NotificationService_SendNotification_FullMethodName         = "/notification.NotificationService/SendNotification"
	NotificationService_GetUnreadNotifications_FullMethodName   = "/notification.NotificationService/GetUnreadNotifications"
	NotificationService_MarkNotificationsAsRead_FullMethodName  = "/notification.NotificationService/MarkNotificationsAsRead"
	NotificationService_SubscribeToNotifications_FullMethodName = "/notification.NotificationService/SubscribeToNotifications"
)

// NotificationServiceClient is the client API for NotificationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationServiceClient interface {
	SendNotification(ctx context.Context, in *SendNotificationRequest, opts ...grpc.CallOption) (*SendNotificationResponse, error)
	GetUnreadNotifications(ctx context.Context, in *GetUnreadNotificationsRequest, opts ...grpc.CallOption) (*GetUnreadNotificationsResponse, error)
	MarkNotificationsAsRead(ctx context.Context, in *MarkNotificationsAsReadRequest, opts ...grpc.CallOption) (*MarkNotificationsAsReadResponse, error)
	SubscribeToNotifications(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[NotificationEvent], error)
}

type notificationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationServiceClient(cc grpc.ClientConnInterface) NotificationServiceClient {
	return &notificationServiceClient{cc}
}

func (c *notificationServiceClient) SendNotification(ctx context.Context, in *SendNotificationRequest, opts ...grpc.CallOption) (*SendNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendNotificationResponse)
	err := c.cc.Invoke(ctx, NotificationService_SendNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) GetUnreadNotifications(ctx context.Context, in *GetUnreadNotificationsRequest, opts ...grpc.CallOption) (*GetUnreadNotificationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUnreadNotificationsResponse)
	err := c.cc.Invoke(ctx, NotificationService_GetUnreadNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) MarkNotificationsAsRead(ctx context.Context, in *MarkNotificationsAsReadRequest, opts ...grpc.CallOption) (*MarkNotificationsAsReadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MarkNotificationsAsReadResponse)
	err := c.cc.Invoke(ctx, NotificationService_MarkNotificationsAsRead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationServiceClient) SubscribeToNotifications(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[NotificationEvent], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &NotificationService_ServiceDesc.Streams[0], NotificationService_SubscribeToNotifications_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeRequest, NotificationEvent]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type NotificationService_SubscribeToNotificationsClient = grpc.ServerStreamingClient[NotificationEvent]

// NotificationServiceServer is the server API for NotificationService service.
// All implementations must embed UnimplementedNotificationServiceServer
// for forward compatibility.
type NotificationServiceServer interface {
	SendNotification(context.Context, *SendNotificationRequest) (*SendNotificationResponse, error)
	GetUnreadNotifications(context.Context, *GetUnreadNotificationsRequest) (*GetUnreadNotificationsResponse, error)
	MarkNotificationsAsRead(context.Context, *MarkNotificationsAsReadRequest) (*MarkNotificationsAsReadResponse, error)
	SubscribeToNotifications(*SubscribeRequest, grpc.ServerStreamingServer[NotificationEvent]) error
	mustEmbedUnimplementedNotificationServiceServer()
}

// UnimplementedNotificationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServiceServer struct{}

func (UnimplementedNotificationServiceServer) SendNotification(context.Context, *SendNotificationRequest) (*SendNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendNotification not implemented")
}
func (UnimplementedNotificationServiceServer) GetUnreadNotifications(context.Context, *GetUnreadNotificationsRequest) (*GetUnreadNotificationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) MarkNotificationsAsRead(context.Context, *MarkNotificationsAsReadRequest) (*MarkNotificationsAsReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkNotificationsAsRead not implemented")
}
func (UnimplementedNotificationServiceServer) SubscribeToNotifications(*SubscribeRequest, grpc.ServerStreamingServer[NotificationEvent]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToNotifications not implemented")
}
func (UnimplementedNotificationServiceServer) mustEmbedUnimplementedNotificationServiceServer() {}
func (UnimplementedNotificationServiceServer) testEmbeddedByValue()                             {}

// UnsafeNotificationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServiceServer will
// result in compilation errors.
type UnsafeNotificationServiceServer interface {
	mustEmbedUnimplementedNotificationServiceServer()
}

func RegisterNotificationServiceServer(s grpc.ServiceRegistrar, srv NotificationServiceServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&NotificationService_ServiceDesc, srv)
}

func _NotificationService_SendNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).SendNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_SendNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).SendNotification(ctx, req.(*SendNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_GetUnreadNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).GetUnreadNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_GetUnreadNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).GetUnreadNotifications(ctx, req.(*GetUnreadNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_MarkNotificationsAsRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkNotificationsAsReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServiceServer).MarkNotificationsAsRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotificationService_MarkNotificationsAsRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServiceServer).MarkNotificationsAsRead(ctx, req.(*MarkNotificationsAsReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotificationService_SubscribeToNotifications_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NotificationServiceServer).SubscribeToNotifications(m, &grpc.GenericServerStream[SubscribeRequest, NotificationEvent]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type NotificationService_SubscribeToNotificationsServer = grpc.ServerStreamingServer[NotificationEvent]

// NotificationService_ServiceDesc is the grpc.ServiceDesc for NotificationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotificationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.NotificationService",
	HandlerType: (*NotificationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendNotification",
			Handler:    _NotificationService_SendNotification_Handler,
		},
		{
			MethodName: "GetUnreadNotifications",
			Handler:    _NotificationService_GetUnreadNotifications_Handler,
		},
		{
			MethodName: "MarkNotificationsAsRead",
			Handler:    _NotificationService_MarkNotificationsAsRead_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToNotifications",
			Handler:       _NotificationService_SubscribeToNotifications_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/grpc/proto/notification/notification.proto",
}
