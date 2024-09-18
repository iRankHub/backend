// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: internal/grpc/proto/debate_management/debate.proto

package debate_management

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
	DebateService_GetRooms_FullMethodName                    = "/debate_management.DebateService/GetRooms"
	DebateService_GetRoom_FullMethodName                     = "/debate_management.DebateService/GetRoom"
	DebateService_UpdateRoom_FullMethodName                  = "/debate_management.DebateService/UpdateRoom"
	DebateService_GetJudges_FullMethodName                   = "/debate_management.DebateService/GetJudges"
	DebateService_GetJudge_FullMethodName                    = "/debate_management.DebateService/GetJudge"
	DebateService_UpdateJudge_FullMethodName                 = "/debate_management.DebateService/UpdateJudge"
	DebateService_GetPairings_FullMethodName                 = "/debate_management.DebateService/GetPairings"
	DebateService_UpdatePairings_FullMethodName              = "/debate_management.DebateService/UpdatePairings"
	DebateService_GetBallots_FullMethodName                  = "/debate_management.DebateService/GetBallots"
	DebateService_GetBallot_FullMethodName                   = "/debate_management.DebateService/GetBallot"
	DebateService_UpdateBallot_FullMethodName                = "/debate_management.DebateService/UpdateBallot"
	DebateService_GetBallotByJudgeID_FullMethodName          = "/debate_management.DebateService/GetBallotByJudgeID"
	DebateService_GeneratePreliminaryPairings_FullMethodName = "/debate_management.DebateService/GeneratePreliminaryPairings"
	DebateService_GenerateEliminationPairings_FullMethodName = "/debate_management.DebateService/GenerateEliminationPairings"
	DebateService_CreateTeam_FullMethodName                  = "/debate_management.DebateService/CreateTeam"
	DebateService_GetTeam_FullMethodName                     = "/debate_management.DebateService/GetTeam"
	DebateService_UpdateTeam_FullMethodName                  = "/debate_management.DebateService/UpdateTeam"
	DebateService_GetTeamsByTournament_FullMethodName        = "/debate_management.DebateService/GetTeamsByTournament"
	DebateService_DeleteTeam_FullMethodName                  = "/debate_management.DebateService/DeleteTeam"
)

// DebateServiceClient is the client API for DebateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DebateServiceClient interface {
	// Room operations
	GetRooms(ctx context.Context, in *GetRoomsRequest, opts ...grpc.CallOption) (*GetRoomsResponse, error)
	GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error)
	UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*UpdateRoomResponse, error)
	// Judge operations
	GetJudges(ctx context.Context, in *GetJudgesRequest, opts ...grpc.CallOption) (*GetJudgesResponse, error)
	GetJudge(ctx context.Context, in *GetJudgeRequest, opts ...grpc.CallOption) (*GetJudgeResponse, error)
	UpdateJudge(ctx context.Context, in *UpdateJudgeRequest, opts ...grpc.CallOption) (*UpdateJudgeResponse, error)
	// Pairing operations
	GetPairings(ctx context.Context, in *GetPairingsRequest, opts ...grpc.CallOption) (*GetPairingsResponse, error)
	UpdatePairings(ctx context.Context, in *UpdatePairingsRequest, opts ...grpc.CallOption) (*UpdatePairingsResponse, error)
	// Ballot operations
	GetBallots(ctx context.Context, in *GetBallotsRequest, opts ...grpc.CallOption) (*GetBallotsResponse, error)
	GetBallot(ctx context.Context, in *GetBallotRequest, opts ...grpc.CallOption) (*GetBallotResponse, error)
	UpdateBallot(ctx context.Context, in *UpdateBallotRequest, opts ...grpc.CallOption) (*UpdateBallotResponse, error)
	GetBallotByJudgeID(ctx context.Context, in *GetBallotByJudgeIDRequest, opts ...grpc.CallOption) (*GetBallotByJudgeIDResponse, error)
	// Algorithm integration
	GeneratePreliminaryPairings(ctx context.Context, in *GeneratePreliminaryPairingsRequest, opts ...grpc.CallOption) (*GeneratePairingsResponse, error)
	GenerateEliminationPairings(ctx context.Context, in *GenerateEliminationPairingsRequest, opts ...grpc.CallOption) (*GeneratePairingsResponse, error)
	// Team operations
	CreateTeam(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*Team, error)
	GetTeam(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*Team, error)
	UpdateTeam(ctx context.Context, in *UpdateTeamRequest, opts ...grpc.CallOption) (*Team, error)
	GetTeamsByTournament(ctx context.Context, in *GetTeamsByTournamentRequest, opts ...grpc.CallOption) (*GetTeamsByTournamentResponse, error)
	DeleteTeam(ctx context.Context, in *DeleteTeamRequest, opts ...grpc.CallOption) (*DeleteTeamResponse, error)
}

type debateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDebateServiceClient(cc grpc.ClientConnInterface) DebateServiceClient {
	return &debateServiceClient{cc}
}

func (c *debateServiceClient) GetRooms(ctx context.Context, in *GetRoomsRequest, opts ...grpc.CallOption) (*GetRoomsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomsResponse)
	err := c.cc.Invoke(ctx, DebateService_GetRooms_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetRoom(ctx context.Context, in *GetRoomRequest, opts ...grpc.CallOption) (*GetRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomResponse)
	err := c.cc.Invoke(ctx, DebateService_GetRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*UpdateRoomResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateRoomResponse)
	err := c.cc.Invoke(ctx, DebateService_UpdateRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetJudges(ctx context.Context, in *GetJudgesRequest, opts ...grpc.CallOption) (*GetJudgesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetJudgesResponse)
	err := c.cc.Invoke(ctx, DebateService_GetJudges_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetJudge(ctx context.Context, in *GetJudgeRequest, opts ...grpc.CallOption) (*GetJudgeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetJudgeResponse)
	err := c.cc.Invoke(ctx, DebateService_GetJudge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) UpdateJudge(ctx context.Context, in *UpdateJudgeRequest, opts ...grpc.CallOption) (*UpdateJudgeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateJudgeResponse)
	err := c.cc.Invoke(ctx, DebateService_UpdateJudge_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetPairings(ctx context.Context, in *GetPairingsRequest, opts ...grpc.CallOption) (*GetPairingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPairingsResponse)
	err := c.cc.Invoke(ctx, DebateService_GetPairings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) UpdatePairings(ctx context.Context, in *UpdatePairingsRequest, opts ...grpc.CallOption) (*UpdatePairingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePairingsResponse)
	err := c.cc.Invoke(ctx, DebateService_UpdatePairings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetBallots(ctx context.Context, in *GetBallotsRequest, opts ...grpc.CallOption) (*GetBallotsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBallotsResponse)
	err := c.cc.Invoke(ctx, DebateService_GetBallots_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetBallot(ctx context.Context, in *GetBallotRequest, opts ...grpc.CallOption) (*GetBallotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBallotResponse)
	err := c.cc.Invoke(ctx, DebateService_GetBallot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) UpdateBallot(ctx context.Context, in *UpdateBallotRequest, opts ...grpc.CallOption) (*UpdateBallotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateBallotResponse)
	err := c.cc.Invoke(ctx, DebateService_UpdateBallot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetBallotByJudgeID(ctx context.Context, in *GetBallotByJudgeIDRequest, opts ...grpc.CallOption) (*GetBallotByJudgeIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBallotByJudgeIDResponse)
	err := c.cc.Invoke(ctx, DebateService_GetBallotByJudgeID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GeneratePreliminaryPairings(ctx context.Context, in *GeneratePreliminaryPairingsRequest, opts ...grpc.CallOption) (*GeneratePairingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GeneratePairingsResponse)
	err := c.cc.Invoke(ctx, DebateService_GeneratePreliminaryPairings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GenerateEliminationPairings(ctx context.Context, in *GenerateEliminationPairingsRequest, opts ...grpc.CallOption) (*GeneratePairingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GeneratePairingsResponse)
	err := c.cc.Invoke(ctx, DebateService_GenerateEliminationPairings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) CreateTeam(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*Team, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Team)
	err := c.cc.Invoke(ctx, DebateService_CreateTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetTeam(ctx context.Context, in *GetTeamRequest, opts ...grpc.CallOption) (*Team, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Team)
	err := c.cc.Invoke(ctx, DebateService_GetTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) UpdateTeam(ctx context.Context, in *UpdateTeamRequest, opts ...grpc.CallOption) (*Team, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Team)
	err := c.cc.Invoke(ctx, DebateService_UpdateTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) GetTeamsByTournament(ctx context.Context, in *GetTeamsByTournamentRequest, opts ...grpc.CallOption) (*GetTeamsByTournamentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTeamsByTournamentResponse)
	err := c.cc.Invoke(ctx, DebateService_GetTeamsByTournament_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debateServiceClient) DeleteTeam(ctx context.Context, in *DeleteTeamRequest, opts ...grpc.CallOption) (*DeleteTeamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteTeamResponse)
	err := c.cc.Invoke(ctx, DebateService_DeleteTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebateServiceServer is the server API for DebateService service.
// All implementations must embed UnimplementedDebateServiceServer
// for forward compatibility.
type DebateServiceServer interface {
	// Room operations
	GetRooms(context.Context, *GetRoomsRequest) (*GetRoomsResponse, error)
	GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error)
	UpdateRoom(context.Context, *UpdateRoomRequest) (*UpdateRoomResponse, error)
	// Judge operations
	GetJudges(context.Context, *GetJudgesRequest) (*GetJudgesResponse, error)
	GetJudge(context.Context, *GetJudgeRequest) (*GetJudgeResponse, error)
	UpdateJudge(context.Context, *UpdateJudgeRequest) (*UpdateJudgeResponse, error)
	// Pairing operations
	GetPairings(context.Context, *GetPairingsRequest) (*GetPairingsResponse, error)
	UpdatePairings(context.Context, *UpdatePairingsRequest) (*UpdatePairingsResponse, error)
	// Ballot operations
	GetBallots(context.Context, *GetBallotsRequest) (*GetBallotsResponse, error)
	GetBallot(context.Context, *GetBallotRequest) (*GetBallotResponse, error)
	UpdateBallot(context.Context, *UpdateBallotRequest) (*UpdateBallotResponse, error)
	GetBallotByJudgeID(context.Context, *GetBallotByJudgeIDRequest) (*GetBallotByJudgeIDResponse, error)
	// Algorithm integration
	GeneratePreliminaryPairings(context.Context, *GeneratePreliminaryPairingsRequest) (*GeneratePairingsResponse, error)
	GenerateEliminationPairings(context.Context, *GenerateEliminationPairingsRequest) (*GeneratePairingsResponse, error)
	// Team operations
	CreateTeam(context.Context, *CreateTeamRequest) (*Team, error)
	GetTeam(context.Context, *GetTeamRequest) (*Team, error)
	UpdateTeam(context.Context, *UpdateTeamRequest) (*Team, error)
	GetTeamsByTournament(context.Context, *GetTeamsByTournamentRequest) (*GetTeamsByTournamentResponse, error)
	DeleteTeam(context.Context, *DeleteTeamRequest) (*DeleteTeamResponse, error)
	mustEmbedUnimplementedDebateServiceServer()
}

// UnimplementedDebateServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDebateServiceServer struct{}

func (UnimplementedDebateServiceServer) GetRooms(context.Context, *GetRoomsRequest) (*GetRoomsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRooms not implemented")
}
func (UnimplementedDebateServiceServer) GetRoom(context.Context, *GetRoomRequest) (*GetRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoom not implemented")
}
func (UnimplementedDebateServiceServer) UpdateRoom(context.Context, *UpdateRoomRequest) (*UpdateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRoom not implemented")
}
func (UnimplementedDebateServiceServer) GetJudges(context.Context, *GetJudgesRequest) (*GetJudgesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJudges not implemented")
}
func (UnimplementedDebateServiceServer) GetJudge(context.Context, *GetJudgeRequest) (*GetJudgeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJudge not implemented")
}
func (UnimplementedDebateServiceServer) UpdateJudge(context.Context, *UpdateJudgeRequest) (*UpdateJudgeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJudge not implemented")
}
func (UnimplementedDebateServiceServer) GetPairings(context.Context, *GetPairingsRequest) (*GetPairingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPairings not implemented")
}
func (UnimplementedDebateServiceServer) UpdatePairings(context.Context, *UpdatePairingsRequest) (*UpdatePairingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePairings not implemented")
}
func (UnimplementedDebateServiceServer) GetBallots(context.Context, *GetBallotsRequest) (*GetBallotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBallots not implemented")
}
func (UnimplementedDebateServiceServer) GetBallot(context.Context, *GetBallotRequest) (*GetBallotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBallot not implemented")
}
func (UnimplementedDebateServiceServer) UpdateBallot(context.Context, *UpdateBallotRequest) (*UpdateBallotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBallot not implemented")
}
func (UnimplementedDebateServiceServer) GetBallotByJudgeID(context.Context, *GetBallotByJudgeIDRequest) (*GetBallotByJudgeIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBallotByJudgeID not implemented")
}
func (UnimplementedDebateServiceServer) GeneratePreliminaryPairings(context.Context, *GeneratePreliminaryPairingsRequest) (*GeneratePairingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratePreliminaryPairings not implemented")
}
func (UnimplementedDebateServiceServer) GenerateEliminationPairings(context.Context, *GenerateEliminationPairingsRequest) (*GeneratePairingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateEliminationPairings not implemented")
}
func (UnimplementedDebateServiceServer) CreateTeam(context.Context, *CreateTeamRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTeam not implemented")
}
func (UnimplementedDebateServiceServer) GetTeam(context.Context, *GetTeamRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTeam not implemented")
}
func (UnimplementedDebateServiceServer) UpdateTeam(context.Context, *UpdateTeamRequest) (*Team, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTeam not implemented")
}
func (UnimplementedDebateServiceServer) GetTeamsByTournament(context.Context, *GetTeamsByTournamentRequest) (*GetTeamsByTournamentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTeamsByTournament not implemented")
}
func (UnimplementedDebateServiceServer) DeleteTeam(context.Context, *DeleteTeamRequest) (*DeleteTeamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTeam not implemented")
}
func (UnimplementedDebateServiceServer) mustEmbedUnimplementedDebateServiceServer() {}
func (UnimplementedDebateServiceServer) testEmbeddedByValue()                       {}

// UnsafeDebateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DebateServiceServer will
// result in compilation errors.
type UnsafeDebateServiceServer interface {
	mustEmbedUnimplementedDebateServiceServer()
}

func RegisterDebateServiceServer(s grpc.ServiceRegistrar, srv DebateServiceServer) {
	// If the following call pancis, it indicates UnimplementedDebateServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DebateService_ServiceDesc, srv)
}

func _DebateService_GetRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetRooms(ctx, req.(*GetRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetRoom(ctx, req.(*GetRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_UpdateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).UpdateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_UpdateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).UpdateRoom(ctx, req.(*UpdateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetJudges_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJudgesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetJudges(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetJudges_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetJudges(ctx, req.(*GetJudgesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetJudge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJudgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetJudge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetJudge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetJudge(ctx, req.(*GetJudgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_UpdateJudge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateJudgeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).UpdateJudge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_UpdateJudge_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).UpdateJudge(ctx, req.(*UpdateJudgeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetPairings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPairingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetPairings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetPairings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetPairings(ctx, req.(*GetPairingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_UpdatePairings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePairingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).UpdatePairings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_UpdatePairings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).UpdatePairings(ctx, req.(*UpdatePairingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetBallots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBallotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetBallots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetBallots_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetBallots(ctx, req.(*GetBallotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetBallot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBallotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetBallot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetBallot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetBallot(ctx, req.(*GetBallotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_UpdateBallot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBallotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).UpdateBallot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_UpdateBallot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).UpdateBallot(ctx, req.(*UpdateBallotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetBallotByJudgeID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBallotByJudgeIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetBallotByJudgeID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetBallotByJudgeID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetBallotByJudgeID(ctx, req.(*GetBallotByJudgeIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GeneratePreliminaryPairings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneratePreliminaryPairingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GeneratePreliminaryPairings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GeneratePreliminaryPairings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GeneratePreliminaryPairings(ctx, req.(*GeneratePreliminaryPairingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GenerateEliminationPairings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateEliminationPairingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GenerateEliminationPairings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GenerateEliminationPairings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GenerateEliminationPairings(ctx, req.(*GenerateEliminationPairingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_CreateTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).CreateTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_CreateTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).CreateTeam(ctx, req.(*CreateTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetTeam(ctx, req.(*GetTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_UpdateTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).UpdateTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_UpdateTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).UpdateTeam(ctx, req.(*UpdateTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_GetTeamsByTournament_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTeamsByTournamentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).GetTeamsByTournament(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_GetTeamsByTournament_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).GetTeamsByTournament(ctx, req.(*GetTeamsByTournamentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebateService_DeleteTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebateServiceServer).DeleteTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebateService_DeleteTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebateServiceServer).DeleteTeam(ctx, req.(*DeleteTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DebateService_ServiceDesc is the grpc.ServiceDesc for DebateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DebateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "debate_management.DebateService",
	HandlerType: (*DebateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRooms",
			Handler:    _DebateService_GetRooms_Handler,
		},
		{
			MethodName: "GetRoom",
			Handler:    _DebateService_GetRoom_Handler,
		},
		{
			MethodName: "UpdateRoom",
			Handler:    _DebateService_UpdateRoom_Handler,
		},
		{
			MethodName: "GetJudges",
			Handler:    _DebateService_GetJudges_Handler,
		},
		{
			MethodName: "GetJudge",
			Handler:    _DebateService_GetJudge_Handler,
		},
		{
			MethodName: "UpdateJudge",
			Handler:    _DebateService_UpdateJudge_Handler,
		},
		{
			MethodName: "GetPairings",
			Handler:    _DebateService_GetPairings_Handler,
		},
		{
			MethodName: "UpdatePairings",
			Handler:    _DebateService_UpdatePairings_Handler,
		},
		{
			MethodName: "GetBallots",
			Handler:    _DebateService_GetBallots_Handler,
		},
		{
			MethodName: "GetBallot",
			Handler:    _DebateService_GetBallot_Handler,
		},
		{
			MethodName: "UpdateBallot",
			Handler:    _DebateService_UpdateBallot_Handler,
		},
		{
			MethodName: "GetBallotByJudgeID",
			Handler:    _DebateService_GetBallotByJudgeID_Handler,
		},
		{
			MethodName: "GeneratePreliminaryPairings",
			Handler:    _DebateService_GeneratePreliminaryPairings_Handler,
		},
		{
			MethodName: "GenerateEliminationPairings",
			Handler:    _DebateService_GenerateEliminationPairings_Handler,
		},
		{
			MethodName: "CreateTeam",
			Handler:    _DebateService_CreateTeam_Handler,
		},
		{
			MethodName: "GetTeam",
			Handler:    _DebateService_GetTeam_Handler,
		},
		{
			MethodName: "UpdateTeam",
			Handler:    _DebateService_UpdateTeam_Handler,
		},
		{
			MethodName: "GetTeamsByTournament",
			Handler:    _DebateService_GetTeamsByTournament_Handler,
		},
		{
			MethodName: "DeleteTeam",
			Handler:    _DebateService_DeleteTeam_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/debate_management/debate.proto",
}
