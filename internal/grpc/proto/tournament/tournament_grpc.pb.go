// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.0
// source: internal/grpc/proto/tournament/tournament.proto

package tournament

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	TournamentService_CreateLeague_FullMethodName           = "/tournament.TournamentService/CreateLeague"
	TournamentService_GetLeague_FullMethodName              = "/tournament.TournamentService/GetLeague"
	TournamentService_ListLeagues_FullMethodName            = "/tournament.TournamentService/ListLeagues"
	TournamentService_UpdateLeague_FullMethodName           = "/tournament.TournamentService/UpdateLeague"
	TournamentService_DeleteLeague_FullMethodName           = "/tournament.TournamentService/DeleteLeague"
	TournamentService_CreateTournamentFormat_FullMethodName = "/tournament.TournamentService/CreateTournamentFormat"
	TournamentService_GetTournamentFormat_FullMethodName    = "/tournament.TournamentService/GetTournamentFormat"
	TournamentService_ListTournamentFormats_FullMethodName  = "/tournament.TournamentService/ListTournamentFormats"
	TournamentService_UpdateTournamentFormat_FullMethodName = "/tournament.TournamentService/UpdateTournamentFormat"
	TournamentService_DeleteTournamentFormat_FullMethodName = "/tournament.TournamentService/DeleteTournamentFormat"
	TournamentService_CreateTournament_FullMethodName       = "/tournament.TournamentService/CreateTournament"
	TournamentService_GetTournament_FullMethodName          = "/tournament.TournamentService/GetTournament"
	TournamentService_ListTournaments_FullMethodName        = "/tournament.TournamentService/ListTournaments"
	TournamentService_UpdateTournament_FullMethodName       = "/tournament.TournamentService/UpdateTournament"
	TournamentService_DeleteTournament_FullMethodName       = "/tournament.TournamentService/DeleteTournament"
)

// TournamentServiceClient is the client API for TournamentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TournamentServiceClient interface {
	// League operations
	CreateLeague(ctx context.Context, in *CreateLeagueRequest, opts ...grpc.CallOption) (*League, error)
	GetLeague(ctx context.Context, in *GetLeagueRequest, opts ...grpc.CallOption) (*League, error)
	ListLeagues(ctx context.Context, in *ListLeaguesRequest, opts ...grpc.CallOption) (*ListLeaguesResponse, error)
	UpdateLeague(ctx context.Context, in *UpdateLeagueRequest, opts ...grpc.CallOption) (*League, error)
	DeleteLeague(ctx context.Context, in *DeleteLeagueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Tournament Format operations
	CreateTournamentFormat(ctx context.Context, in *CreateTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error)
	GetTournamentFormat(ctx context.Context, in *GetTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error)
	ListTournamentFormats(ctx context.Context, in *ListTournamentFormatsRequest, opts ...grpc.CallOption) (*ListTournamentFormatsResponse, error)
	UpdateTournamentFormat(ctx context.Context, in *UpdateTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error)
	DeleteTournamentFormat(ctx context.Context, in *DeleteTournamentFormatRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Tournament operations
	CreateTournament(ctx context.Context, in *CreateTournamentRequest, opts ...grpc.CallOption) (*Tournament, error)
	GetTournament(ctx context.Context, in *GetTournamentRequest, opts ...grpc.CallOption) (*Tournament, error)
	ListTournaments(ctx context.Context, in *ListTournamentsRequest, opts ...grpc.CallOption) (*ListTournamentsResponse, error)
	UpdateTournament(ctx context.Context, in *UpdateTournamentRequest, opts ...grpc.CallOption) (*Tournament, error)
	DeleteTournament(ctx context.Context, in *DeleteTournamentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type tournamentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTournamentServiceClient(cc grpc.ClientConnInterface) TournamentServiceClient {
	return &tournamentServiceClient{cc}
}

func (c *tournamentServiceClient) CreateLeague(ctx context.Context, in *CreateLeagueRequest, opts ...grpc.CallOption) (*League, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(League)
	err := c.cc.Invoke(ctx, TournamentService_CreateLeague_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) GetLeague(ctx context.Context, in *GetLeagueRequest, opts ...grpc.CallOption) (*League, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(League)
	err := c.cc.Invoke(ctx, TournamentService_GetLeague_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) ListLeagues(ctx context.Context, in *ListLeaguesRequest, opts ...grpc.CallOption) (*ListLeaguesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLeaguesResponse)
	err := c.cc.Invoke(ctx, TournamentService_ListLeagues_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) UpdateLeague(ctx context.Context, in *UpdateLeagueRequest, opts ...grpc.CallOption) (*League, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(League)
	err := c.cc.Invoke(ctx, TournamentService_UpdateLeague_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) DeleteLeague(ctx context.Context, in *DeleteLeagueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, TournamentService_DeleteLeague_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) CreateTournamentFormat(ctx context.Context, in *CreateTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TournamentFormat)
	err := c.cc.Invoke(ctx, TournamentService_CreateTournamentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) GetTournamentFormat(ctx context.Context, in *GetTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TournamentFormat)
	err := c.cc.Invoke(ctx, TournamentService_GetTournamentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) ListTournamentFormats(ctx context.Context, in *ListTournamentFormatsRequest, opts ...grpc.CallOption) (*ListTournamentFormatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTournamentFormatsResponse)
	err := c.cc.Invoke(ctx, TournamentService_ListTournamentFormats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) UpdateTournamentFormat(ctx context.Context, in *UpdateTournamentFormatRequest, opts ...grpc.CallOption) (*TournamentFormat, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TournamentFormat)
	err := c.cc.Invoke(ctx, TournamentService_UpdateTournamentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) DeleteTournamentFormat(ctx context.Context, in *DeleteTournamentFormatRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, TournamentService_DeleteTournamentFormat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) CreateTournament(ctx context.Context, in *CreateTournamentRequest, opts ...grpc.CallOption) (*Tournament, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tournament)
	err := c.cc.Invoke(ctx, TournamentService_CreateTournament_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) GetTournament(ctx context.Context, in *GetTournamentRequest, opts ...grpc.CallOption) (*Tournament, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tournament)
	err := c.cc.Invoke(ctx, TournamentService_GetTournament_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) ListTournaments(ctx context.Context, in *ListTournamentsRequest, opts ...grpc.CallOption) (*ListTournamentsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTournamentsResponse)
	err := c.cc.Invoke(ctx, TournamentService_ListTournaments_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) UpdateTournament(ctx context.Context, in *UpdateTournamentRequest, opts ...grpc.CallOption) (*Tournament, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Tournament)
	err := c.cc.Invoke(ctx, TournamentService_UpdateTournament_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tournamentServiceClient) DeleteTournament(ctx context.Context, in *DeleteTournamentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, TournamentService_DeleteTournament_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TournamentServiceServer is the server API for TournamentService service.
// All implementations must embed UnimplementedTournamentServiceServer
// for forward compatibility
type TournamentServiceServer interface {
	// League operations
	CreateLeague(context.Context, *CreateLeagueRequest) (*League, error)
	GetLeague(context.Context, *GetLeagueRequest) (*League, error)
	ListLeagues(context.Context, *ListLeaguesRequest) (*ListLeaguesResponse, error)
	UpdateLeague(context.Context, *UpdateLeagueRequest) (*League, error)
	DeleteLeague(context.Context, *DeleteLeagueRequest) (*emptypb.Empty, error)
	// Tournament Format operations
	CreateTournamentFormat(context.Context, *CreateTournamentFormatRequest) (*TournamentFormat, error)
	GetTournamentFormat(context.Context, *GetTournamentFormatRequest) (*TournamentFormat, error)
	ListTournamentFormats(context.Context, *ListTournamentFormatsRequest) (*ListTournamentFormatsResponse, error)
	UpdateTournamentFormat(context.Context, *UpdateTournamentFormatRequest) (*TournamentFormat, error)
	DeleteTournamentFormat(context.Context, *DeleteTournamentFormatRequest) (*emptypb.Empty, error)
	// Tournament operations
	CreateTournament(context.Context, *CreateTournamentRequest) (*Tournament, error)
	GetTournament(context.Context, *GetTournamentRequest) (*Tournament, error)
	ListTournaments(context.Context, *ListTournamentsRequest) (*ListTournamentsResponse, error)
	UpdateTournament(context.Context, *UpdateTournamentRequest) (*Tournament, error)
	DeleteTournament(context.Context, *DeleteTournamentRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedTournamentServiceServer()
}

// UnimplementedTournamentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTournamentServiceServer struct {
}

func (UnimplementedTournamentServiceServer) CreateLeague(context.Context, *CreateLeagueRequest) (*League, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLeague not implemented")
}
func (UnimplementedTournamentServiceServer) GetLeague(context.Context, *GetLeagueRequest) (*League, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeague not implemented")
}
func (UnimplementedTournamentServiceServer) ListLeagues(context.Context, *ListLeaguesRequest) (*ListLeaguesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLeagues not implemented")
}
func (UnimplementedTournamentServiceServer) UpdateLeague(context.Context, *UpdateLeagueRequest) (*League, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLeague not implemented")
}
func (UnimplementedTournamentServiceServer) DeleteLeague(context.Context, *DeleteLeagueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLeague not implemented")
}
func (UnimplementedTournamentServiceServer) CreateTournamentFormat(context.Context, *CreateTournamentFormatRequest) (*TournamentFormat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTournamentFormat not implemented")
}
func (UnimplementedTournamentServiceServer) GetTournamentFormat(context.Context, *GetTournamentFormatRequest) (*TournamentFormat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTournamentFormat not implemented")
}
func (UnimplementedTournamentServiceServer) ListTournamentFormats(context.Context, *ListTournamentFormatsRequest) (*ListTournamentFormatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTournamentFormats not implemented")
}
func (UnimplementedTournamentServiceServer) UpdateTournamentFormat(context.Context, *UpdateTournamentFormatRequest) (*TournamentFormat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTournamentFormat not implemented")
}
func (UnimplementedTournamentServiceServer) DeleteTournamentFormat(context.Context, *DeleteTournamentFormatRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTournamentFormat not implemented")
}
func (UnimplementedTournamentServiceServer) CreateTournament(context.Context, *CreateTournamentRequest) (*Tournament, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTournament not implemented")
}
func (UnimplementedTournamentServiceServer) GetTournament(context.Context, *GetTournamentRequest) (*Tournament, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTournament not implemented")
}
func (UnimplementedTournamentServiceServer) ListTournaments(context.Context, *ListTournamentsRequest) (*ListTournamentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTournaments not implemented")
}
func (UnimplementedTournamentServiceServer) UpdateTournament(context.Context, *UpdateTournamentRequest) (*Tournament, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTournament not implemented")
}
func (UnimplementedTournamentServiceServer) DeleteTournament(context.Context, *DeleteTournamentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTournament not implemented")
}
func (UnimplementedTournamentServiceServer) mustEmbedUnimplementedTournamentServiceServer() {}

// UnsafeTournamentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TournamentServiceServer will
// result in compilation errors.
type UnsafeTournamentServiceServer interface {
	mustEmbedUnimplementedTournamentServiceServer()
}

func RegisterTournamentServiceServer(s grpc.ServiceRegistrar, srv TournamentServiceServer) {
	s.RegisterService(&TournamentService_ServiceDesc, srv)
}

func _TournamentService_CreateLeague_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLeagueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).CreateLeague(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_CreateLeague_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).CreateLeague(ctx, req.(*CreateLeagueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_GetLeague_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeagueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).GetLeague(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_GetLeague_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).GetLeague(ctx, req.(*GetLeagueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_ListLeagues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLeaguesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).ListLeagues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_ListLeagues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).ListLeagues(ctx, req.(*ListLeaguesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_UpdateLeague_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLeagueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).UpdateLeague(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_UpdateLeague_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).UpdateLeague(ctx, req.(*UpdateLeagueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_DeleteLeague_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLeagueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).DeleteLeague(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_DeleteLeague_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).DeleteLeague(ctx, req.(*DeleteLeagueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_CreateTournamentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTournamentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).CreateTournamentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_CreateTournamentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).CreateTournamentFormat(ctx, req.(*CreateTournamentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_GetTournamentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTournamentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).GetTournamentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_GetTournamentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).GetTournamentFormat(ctx, req.(*GetTournamentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_ListTournamentFormats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTournamentFormatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).ListTournamentFormats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_ListTournamentFormats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).ListTournamentFormats(ctx, req.(*ListTournamentFormatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_UpdateTournamentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTournamentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).UpdateTournamentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_UpdateTournamentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).UpdateTournamentFormat(ctx, req.(*UpdateTournamentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_DeleteTournamentFormat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTournamentFormatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).DeleteTournamentFormat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_DeleteTournamentFormat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).DeleteTournamentFormat(ctx, req.(*DeleteTournamentFormatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_CreateTournament_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).CreateTournament(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_CreateTournament_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).CreateTournament(ctx, req.(*CreateTournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_GetTournament_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).GetTournament(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_GetTournament_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).GetTournament(ctx, req.(*GetTournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_ListTournaments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTournamentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).ListTournaments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_ListTournaments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).ListTournaments(ctx, req.(*ListTournamentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_UpdateTournament_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).UpdateTournament(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_UpdateTournament_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).UpdateTournament(ctx, req.(*UpdateTournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TournamentService_DeleteTournament_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TournamentServiceServer).DeleteTournament(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TournamentService_DeleteTournament_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TournamentServiceServer).DeleteTournament(ctx, req.(*DeleteTournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TournamentService_ServiceDesc is the grpc.ServiceDesc for TournamentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TournamentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tournament.TournamentService",
	HandlerType: (*TournamentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLeague",
			Handler:    _TournamentService_CreateLeague_Handler,
		},
		{
			MethodName: "GetLeague",
			Handler:    _TournamentService_GetLeague_Handler,
		},
		{
			MethodName: "ListLeagues",
			Handler:    _TournamentService_ListLeagues_Handler,
		},
		{
			MethodName: "UpdateLeague",
			Handler:    _TournamentService_UpdateLeague_Handler,
		},
		{
			MethodName: "DeleteLeague",
			Handler:    _TournamentService_DeleteLeague_Handler,
		},
		{
			MethodName: "CreateTournamentFormat",
			Handler:    _TournamentService_CreateTournamentFormat_Handler,
		},
		{
			MethodName: "GetTournamentFormat",
			Handler:    _TournamentService_GetTournamentFormat_Handler,
		},
		{
			MethodName: "ListTournamentFormats",
			Handler:    _TournamentService_ListTournamentFormats_Handler,
		},
		{
			MethodName: "UpdateTournamentFormat",
			Handler:    _TournamentService_UpdateTournamentFormat_Handler,
		},
		{
			MethodName: "DeleteTournamentFormat",
			Handler:    _TournamentService_DeleteTournamentFormat_Handler,
		},
		{
			MethodName: "CreateTournament",
			Handler:    _TournamentService_CreateTournament_Handler,
		},
		{
			MethodName: "GetTournament",
			Handler:    _TournamentService_GetTournament_Handler,
		},
		{
			MethodName: "ListTournaments",
			Handler:    _TournamentService_ListTournaments_Handler,
		},
		{
			MethodName: "UpdateTournament",
			Handler:    _TournamentService_UpdateTournament_Handler,
		},
		{
			MethodName: "DeleteTournament",
			Handler:    _TournamentService_DeleteTournament_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/tournament/tournament.proto",
}
