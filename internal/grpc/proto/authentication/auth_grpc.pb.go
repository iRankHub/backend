// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: internal/grpc/proto/authentication/auth.proto

package authentication

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
	AuthService_SignUp_FullMethodName                     = "/auth.AuthService/SignUp"
	AuthService_BatchImportUsers_FullMethodName           = "/auth.AuthService/BatchImportUsers"
	AuthService_AdminLogin_FullMethodName                 = "/auth.AuthService/AdminLogin"
	AuthService_StudentLogin_FullMethodName               = "/auth.AuthService/StudentLogin"
	AuthService_VolunteerLogin_FullMethodName             = "/auth.AuthService/VolunteerLogin"
	AuthService_SchoolLogin_FullMethodName                = "/auth.AuthService/SchoolLogin"
	AuthService_EnableTwoFactor_FullMethodName            = "/auth.AuthService/EnableTwoFactor"
	AuthService_DisableTwoFactor_FullMethodName           = "/auth.AuthService/DisableTwoFactor"
	AuthService_GenerateTwoFactorOTP_FullMethodName       = "/auth.AuthService/GenerateTwoFactorOTP"
	AuthService_VerifyTwoFactor_FullMethodName            = "/auth.AuthService/VerifyTwoFactor"
	AuthService_RequestPasswordReset_FullMethodName       = "/auth.AuthService/RequestPasswordReset"
	AuthService_ResetPassword_FullMethodName              = "/auth.AuthService/ResetPassword"
	AuthService_BeginWebAuthnRegistration_FullMethodName  = "/auth.AuthService/BeginWebAuthnRegistration"
	AuthService_FinishWebAuthnRegistration_FullMethodName = "/auth.AuthService/FinishWebAuthnRegistration"
	AuthService_BeginWebAuthnLogin_FullMethodName         = "/auth.AuthService/BeginWebAuthnLogin"
	AuthService_FinishWebAuthnLogin_FullMethodName        = "/auth.AuthService/FinishWebAuthnLogin"
	AuthService_Logout_FullMethodName                     = "/auth.AuthService/Logout"
)

// AuthServiceClient is the client API for AuthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthServiceClient interface {
	SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error)
	BatchImportUsers(ctx context.Context, in *BatchImportUsersRequest, opts ...grpc.CallOption) (*BatchImportUsersResponse, error)
	AdminLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	StudentLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	VolunteerLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	SchoolLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	EnableTwoFactor(ctx context.Context, in *EnableTwoFactorRequest, opts ...grpc.CallOption) (*EnableTwoFactorResponse, error)
	DisableTwoFactor(ctx context.Context, in *DisableTwoFactorRequest, opts ...grpc.CallOption) (*DisableTwoFactorResponse, error)
	GenerateTwoFactorOTP(ctx context.Context, in *GenerateTwoFactorOTPRequest, opts ...grpc.CallOption) (*GenerateTwoFactorOTPResponse, error)
	VerifyTwoFactor(ctx context.Context, in *VerifyTwoFactorRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	RequestPasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetResponse, error)
	ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error)
	BeginWebAuthnRegistration(ctx context.Context, in *BeginWebAuthnRegistrationRequest, opts ...grpc.CallOption) (*BeginWebAuthnRegistrationResponse, error)
	FinishWebAuthnRegistration(ctx context.Context, in *FinishWebAuthnRegistrationRequest, opts ...grpc.CallOption) (*FinishWebAuthnRegistrationResponse, error)
	BeginWebAuthnLogin(ctx context.Context, in *BeginWebAuthnLoginRequest, opts ...grpc.CallOption) (*BeginWebAuthnLoginResponse, error)
	FinishWebAuthnLogin(ctx context.Context, in *FinishWebAuthnLoginRequest, opts ...grpc.CallOption) (*FinishWebAuthnLoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
}

type authServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthServiceClient(cc grpc.ClientConnInterface) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) SignUp(ctx context.Context, in *SignUpRequest, opts ...grpc.CallOption) (*SignUpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignUpResponse)
	err := c.cc.Invoke(ctx, AuthService_SignUp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) BatchImportUsers(ctx context.Context, in *BatchImportUsersRequest, opts ...grpc.CallOption) (*BatchImportUsersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchImportUsersResponse)
	err := c.cc.Invoke(ctx, AuthService_BatchImportUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) AdminLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_AdminLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) StudentLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_StudentLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) VolunteerLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_VolunteerLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) SchoolLogin(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_SchoolLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) EnableTwoFactor(ctx context.Context, in *EnableTwoFactorRequest, opts ...grpc.CallOption) (*EnableTwoFactorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EnableTwoFactorResponse)
	err := c.cc.Invoke(ctx, AuthService_EnableTwoFactor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) DisableTwoFactor(ctx context.Context, in *DisableTwoFactorRequest, opts ...grpc.CallOption) (*DisableTwoFactorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DisableTwoFactorResponse)
	err := c.cc.Invoke(ctx, AuthService_DisableTwoFactor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) GenerateTwoFactorOTP(ctx context.Context, in *GenerateTwoFactorOTPRequest, opts ...grpc.CallOption) (*GenerateTwoFactorOTPResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateTwoFactorOTPResponse)
	err := c.cc.Invoke(ctx, AuthService_GenerateTwoFactorOTP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) VerifyTwoFactor(ctx context.Context, in *VerifyTwoFactorRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, AuthService_VerifyTwoFactor_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) RequestPasswordReset(ctx context.Context, in *PasswordResetRequest, opts ...grpc.CallOption) (*PasswordResetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PasswordResetResponse)
	err := c.cc.Invoke(ctx, AuthService_RequestPasswordReset_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResetPasswordResponse)
	err := c.cc.Invoke(ctx, AuthService_ResetPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) BeginWebAuthnRegistration(ctx context.Context, in *BeginWebAuthnRegistrationRequest, opts ...grpc.CallOption) (*BeginWebAuthnRegistrationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BeginWebAuthnRegistrationResponse)
	err := c.cc.Invoke(ctx, AuthService_BeginWebAuthnRegistration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) FinishWebAuthnRegistration(ctx context.Context, in *FinishWebAuthnRegistrationRequest, opts ...grpc.CallOption) (*FinishWebAuthnRegistrationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FinishWebAuthnRegistrationResponse)
	err := c.cc.Invoke(ctx, AuthService_FinishWebAuthnRegistration_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) BeginWebAuthnLogin(ctx context.Context, in *BeginWebAuthnLoginRequest, opts ...grpc.CallOption) (*BeginWebAuthnLoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BeginWebAuthnLoginResponse)
	err := c.cc.Invoke(ctx, AuthService_BeginWebAuthnLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) FinishWebAuthnLogin(ctx context.Context, in *FinishWebAuthnLoginRequest, opts ...grpc.CallOption) (*FinishWebAuthnLoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FinishWebAuthnLoginResponse)
	err := c.cc.Invoke(ctx, AuthService_FinishWebAuthnLogin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, AuthService_Logout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServiceServer is the server API for AuthService service.
// All implementations must embed UnimplementedAuthServiceServer
// for forward compatibility.
type AuthServiceServer interface {
	SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error)
	BatchImportUsers(context.Context, *BatchImportUsersRequest) (*BatchImportUsersResponse, error)
	AdminLogin(context.Context, *LoginRequest) (*LoginResponse, error)
	StudentLogin(context.Context, *LoginRequest) (*LoginResponse, error)
	VolunteerLogin(context.Context, *LoginRequest) (*LoginResponse, error)
	SchoolLogin(context.Context, *LoginRequest) (*LoginResponse, error)
	EnableTwoFactor(context.Context, *EnableTwoFactorRequest) (*EnableTwoFactorResponse, error)
	DisableTwoFactor(context.Context, *DisableTwoFactorRequest) (*DisableTwoFactorResponse, error)
	GenerateTwoFactorOTP(context.Context, *GenerateTwoFactorOTPRequest) (*GenerateTwoFactorOTPResponse, error)
	VerifyTwoFactor(context.Context, *VerifyTwoFactorRequest) (*LoginResponse, error)
	RequestPasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetResponse, error)
	ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error)
	BeginWebAuthnRegistration(context.Context, *BeginWebAuthnRegistrationRequest) (*BeginWebAuthnRegistrationResponse, error)
	FinishWebAuthnRegistration(context.Context, *FinishWebAuthnRegistrationRequest) (*FinishWebAuthnRegistrationResponse, error)
	BeginWebAuthnLogin(context.Context, *BeginWebAuthnLoginRequest) (*BeginWebAuthnLoginResponse, error)
	FinishWebAuthnLogin(context.Context, *FinishWebAuthnLoginRequest) (*FinishWebAuthnLoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	mustEmbedUnimplementedAuthServiceServer()
}

// UnimplementedAuthServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthServiceServer struct{}

func (UnimplementedAuthServiceServer) SignUp(context.Context, *SignUpRequest) (*SignUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedAuthServiceServer) BatchImportUsers(context.Context, *BatchImportUsersRequest) (*BatchImportUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchImportUsers not implemented")
}
func (UnimplementedAuthServiceServer) AdminLogin(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdminLogin not implemented")
}
func (UnimplementedAuthServiceServer) StudentLogin(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StudentLogin not implemented")
}
func (UnimplementedAuthServiceServer) VolunteerLogin(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VolunteerLogin not implemented")
}
func (UnimplementedAuthServiceServer) SchoolLogin(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SchoolLogin not implemented")
}
func (UnimplementedAuthServiceServer) EnableTwoFactor(context.Context, *EnableTwoFactorRequest) (*EnableTwoFactorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableTwoFactor not implemented")
}
func (UnimplementedAuthServiceServer) DisableTwoFactor(context.Context, *DisableTwoFactorRequest) (*DisableTwoFactorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DisableTwoFactor not implemented")
}
func (UnimplementedAuthServiceServer) GenerateTwoFactorOTP(context.Context, *GenerateTwoFactorOTPRequest) (*GenerateTwoFactorOTPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateTwoFactorOTP not implemented")
}
func (UnimplementedAuthServiceServer) VerifyTwoFactor(context.Context, *VerifyTwoFactorRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyTwoFactor not implemented")
}
func (UnimplementedAuthServiceServer) RequestPasswordReset(context.Context, *PasswordResetRequest) (*PasswordResetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestPasswordReset not implemented")
}
func (UnimplementedAuthServiceServer) ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPassword not implemented")
}
func (UnimplementedAuthServiceServer) BeginWebAuthnRegistration(context.Context, *BeginWebAuthnRegistrationRequest) (*BeginWebAuthnRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginWebAuthnRegistration not implemented")
}
func (UnimplementedAuthServiceServer) FinishWebAuthnRegistration(context.Context, *FinishWebAuthnRegistrationRequest) (*FinishWebAuthnRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishWebAuthnRegistration not implemented")
}
func (UnimplementedAuthServiceServer) BeginWebAuthnLogin(context.Context, *BeginWebAuthnLoginRequest) (*BeginWebAuthnLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginWebAuthnLogin not implemented")
}
func (UnimplementedAuthServiceServer) FinishWebAuthnLogin(context.Context, *FinishWebAuthnLoginRequest) (*FinishWebAuthnLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinishWebAuthnLogin not implemented")
}
func (UnimplementedAuthServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAuthServiceServer) mustEmbedUnimplementedAuthServiceServer() {}
func (UnimplementedAuthServiceServer) testEmbeddedByValue()                     {}

// UnsafeAuthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServiceServer will
// result in compilation errors.
type UnsafeAuthServiceServer interface {
	mustEmbedUnimplementedAuthServiceServer()
}

func RegisterAuthServiceServer(s grpc.ServiceRegistrar, srv AuthServiceServer) {
	// If the following call pancis, it indicates UnimplementedAuthServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AuthService_ServiceDesc, srv)
}

func _AuthService_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_SignUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).SignUp(ctx, req.(*SignUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_BatchImportUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchImportUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).BatchImportUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_BatchImportUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).BatchImportUsers(ctx, req.(*BatchImportUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_AdminLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).AdminLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_AdminLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).AdminLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_StudentLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).StudentLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_StudentLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).StudentLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_VolunteerLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).VolunteerLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_VolunteerLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).VolunteerLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_SchoolLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).SchoolLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_SchoolLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).SchoolLogin(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_EnableTwoFactor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnableTwoFactorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).EnableTwoFactor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_EnableTwoFactor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).EnableTwoFactor(ctx, req.(*EnableTwoFactorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_DisableTwoFactor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DisableTwoFactorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).DisableTwoFactor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_DisableTwoFactor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).DisableTwoFactor(ctx, req.(*DisableTwoFactorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_GenerateTwoFactorOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateTwoFactorOTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).GenerateTwoFactorOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_GenerateTwoFactorOTP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).GenerateTwoFactorOTP(ctx, req.(*GenerateTwoFactorOTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_VerifyTwoFactor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyTwoFactorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).VerifyTwoFactor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_VerifyTwoFactor_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).VerifyTwoFactor(ctx, req.(*VerifyTwoFactorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_RequestPasswordReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PasswordResetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).RequestPasswordReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_RequestPasswordReset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).RequestPasswordReset(ctx, req.(*PasswordResetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_ResetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).ResetPassword(ctx, req.(*ResetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_BeginWebAuthnRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginWebAuthnRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).BeginWebAuthnRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_BeginWebAuthnRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).BeginWebAuthnRegistration(ctx, req.(*BeginWebAuthnRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_FinishWebAuthnRegistration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishWebAuthnRegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).FinishWebAuthnRegistration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_FinishWebAuthnRegistration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).FinishWebAuthnRegistration(ctx, req.(*FinishWebAuthnRegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_BeginWebAuthnLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BeginWebAuthnLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).BeginWebAuthnLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_BeginWebAuthnLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).BeginWebAuthnLogin(ctx, req.(*BeginWebAuthnLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_FinishWebAuthnLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinishWebAuthnLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).FinishWebAuthnLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_FinishWebAuthnLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).FinishWebAuthnLogin(ctx, req.(*FinishWebAuthnLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AuthService_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthService_ServiceDesc is the grpc.ServiceDesc for AuthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _AuthService_SignUp_Handler,
		},
		{
			MethodName: "BatchImportUsers",
			Handler:    _AuthService_BatchImportUsers_Handler,
		},
		{
			MethodName: "AdminLogin",
			Handler:    _AuthService_AdminLogin_Handler,
		},
		{
			MethodName: "StudentLogin",
			Handler:    _AuthService_StudentLogin_Handler,
		},
		{
			MethodName: "VolunteerLogin",
			Handler:    _AuthService_VolunteerLogin_Handler,
		},
		{
			MethodName: "SchoolLogin",
			Handler:    _AuthService_SchoolLogin_Handler,
		},
		{
			MethodName: "EnableTwoFactor",
			Handler:    _AuthService_EnableTwoFactor_Handler,
		},
		{
			MethodName: "DisableTwoFactor",
			Handler:    _AuthService_DisableTwoFactor_Handler,
		},
		{
			MethodName: "GenerateTwoFactorOTP",
			Handler:    _AuthService_GenerateTwoFactorOTP_Handler,
		},
		{
			MethodName: "VerifyTwoFactor",
			Handler:    _AuthService_VerifyTwoFactor_Handler,
		},
		{
			MethodName: "RequestPasswordReset",
			Handler:    _AuthService_RequestPasswordReset_Handler,
		},
		{
			MethodName: "ResetPassword",
			Handler:    _AuthService_ResetPassword_Handler,
		},
		{
			MethodName: "BeginWebAuthnRegistration",
			Handler:    _AuthService_BeginWebAuthnRegistration_Handler,
		},
		{
			MethodName: "FinishWebAuthnRegistration",
			Handler:    _AuthService_FinishWebAuthnRegistration_Handler,
		},
		{
			MethodName: "BeginWebAuthnLogin",
			Handler:    _AuthService_BeginWebAuthnLogin_Handler,
		},
		{
			MethodName: "FinishWebAuthnLogin",
			Handler:    _AuthService_FinishWebAuthnLogin_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AuthService_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/authentication/auth.proto",
}
