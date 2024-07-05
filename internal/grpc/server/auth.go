package server

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/services"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	proto.UnimplementedAuthServiceServer
	queries           *models.Queries
	authService       *services.AuthService
	twoFactorService  *services.TwoFactorService
	recoveryService   *services.RecoveryService
	biometricService  *services.BiometricService
	privateKey        ed25519.PrivateKey
}

func NewAuthServer(queries *models.Queries, privateKey ed25519.PrivateKey) (proto.AuthServiceServer, error) {
	return &authServer{
		queries:           queries,
		authService:       services.NewAuthService(queries, privateKey),
		twoFactorService:  services.NewTwoFactorService(queries),
		recoveryService:   services.NewRecoveryService(queries),
		biometricService:  services.NewBiometricService(queries),
		privateKey:        privateKey,
	}, nil
}

func (s *authServer) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	return s.authService.SignUp(ctx, req)
}

func (s *authServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return s.authService.Login(ctx, req)
}

func (s *authServer) EnableTwoFactor(ctx context.Context, req *proto.EnableTwoFactorRequest) (*proto.EnableTwoFactorResponse, error) {
	secret, qrCode, err := s.twoFactorService.EnableTwoFactor(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	return &proto.EnableTwoFactorResponse{
		Secret: secret,
		QrCode: qrCode,
	}, nil
}

func (s *authServer) VerifyTwoFactor(ctx context.Context, req *proto.VerifyTwoFactorRequest) (*proto.VerifyTwoFactorResponse, error) {
	user, err := s.queries.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if !user.TwoFactorSecret.Valid || !s.twoFactorService.ValidateCode(user.TwoFactorSecret.String, req.Code) {
		return nil, fmt.Errorf("invalid two-factor code")
	}

	err = s.queries.EnableTwoFactor(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to enable two-factor authentication: %v", err)
	}

	return &proto.VerifyTwoFactorResponse{Success: true}, nil
}

func (s *authServer) RequestPasswordReset(ctx context.Context, req *proto.PasswordResetRequest) (*proto.PasswordResetResponse, error) {
	err := s.recoveryService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to request password reset: %v", err)
	}

	return &proto.PasswordResetResponse{Success: true}, nil
}

func (s *authServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	err := s.recoveryService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to reset password: %v", err)
	}

	return &proto.ResetPasswordResponse{Success: true}, nil
}

func (s *authServer) EnableBiometricLogin(ctx context.Context, req *proto.EnableBiometricLoginRequest) (*proto.EnableBiometricLoginResponse, error) {
	token, err := s.biometricService.EnableBiometricLogin(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to enable biometric login: %v", err)
	}

	return &proto.EnableBiometricLoginResponse{BiometricToken: token}, nil
}

func (s *authServer) BiometricLogin(ctx context.Context, req *proto.BiometricLoginRequest) (*proto.LoginResponse, error) {
	user, err := s.biometricService.VerifyBiometricToken(ctx, req.BiometricToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify biometric token: %v", err)
	}

	token, err := utils.GenerateToken(user.Userid, user.Userrole, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &proto.LoginResponse{
		Success:  true,
		Token:    token,
		UserRole: user.Userrole,
		UserID:   user.Userid,
	}, nil
}