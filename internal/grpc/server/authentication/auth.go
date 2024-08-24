package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/models"
	services "github.com/iRankHub/backend/internal/services/authentication"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	authentication.UnimplementedAuthServiceServer
	db                 *sql.DB
	webauthn           *webauthn.WebAuthn
	loginService       *services.LoginService
	signUpService      *services.SignUpService
	importUsersService *services.ImportUsersService
	twoFactorService   *services.TwoFactorService
	recoveryService    *services.RecoveryService
	biometricService   *services.BiometricService
}

func NewAuthServer(db *sql.DB) (authentication.AuthServiceServer, error) {
	wconfig := &webauthn.Config{
		RPDisplayName: "iRankHub",
		RPID:          "localhost",             // TODO: Change this to actual domain in production
		RPOrigin:      "http://localhost:3000", // TODO: Change this to origin url in production
	}

	w, err := webauthn.New(wconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebAuthn: %v", err)
	}

	twoFactorService := services.NewTwoFactorService(db)
	recoveryService := services.NewRecoveryService(db)
	biometricService := services.NewBiometricService(db, w)
	loginService := services.NewLoginService(db, twoFactorService, recoveryService)
	signUpService := services.NewSignUpService(db)
	importUsersService := services.NewImportUsersService(signUpService)

	return &authServer{
		db:                 db,
		webauthn:           w,
		loginService:       loginService,
		signUpService:      signUpService,
		importUsersService: importUsersService,
		twoFactorService:   twoFactorService,
		recoveryService:    recoveryService,
		biometricService:   biometricService,
	}, nil
}

func (s *authServer) SignUp(ctx context.Context, req *authentication.SignUpRequest) (*authentication.SignUpResponse, error) {
	additionalInfo := map[string]interface{}{
		"dateOfBirth":             req.DateOfBirth,
		"schoolID":                req.SchoolID,
		"schoolName":              req.SchoolName,
		"address":                 req.Address,
		"country":                 req.Country,
		"province":                req.Province,
		"district":                req.District,
		"contactEmail":            req.ContactEmail,
		"schoolType":              req.SchoolType,
		"roleInterestedIn":        req.RoleInterestedIn,
		"graduationYear":          req.GraduationYear,
		"safeguardingCertificate": req.SafeguardingCertificate,
		"grade":                   req.Grade,
		"hasInternship":           req.HasInternship,
		"isEnrolledInUniversity":  req.IsEnrolledInUniversity,
	}

	err := s.signUpService.SignUp(ctx, req.FirstName, req.LastName, req.Email, req.Password, req.UserRole, req.Gender, additionalInfo)
	if err != nil {
		return nil, err
	}

	var message, status string
	if req.UserRole == "admin" {
		message = "Admin account created successfully. You can now log in to the system."
		status = "approved"
	} else {
		message = "Sign-up successful. Please wait for admin approval."
		status = "pending"
	}

	return &authentication.SignUpResponse{Success: true, Message: message, Status: status}, nil
}

func (s *authServer) BatchImportUsers(ctx context.Context, req *authentication.BatchImportUsersRequest) (*authentication.BatchImportUsersResponse, error) {
	importedCount, failedEmails := s.importUsersService.BatchImportUsers(ctx, req.Users)

	return &authentication.BatchImportUsersResponse{
		Success:       importedCount > 0,
		Message:       fmt.Sprintf("Imported %d users successfully", importedCount),
		ImportedCount: importedCount,
		FailedEmails:  failedEmails,
	}, nil
}

func (s *authServer) Login(ctx context.Context, req *authentication.LoginRequest) (*authentication.LoginResponse, error) {
	user, err := s.loginService.Login(ctx, req.EmailOrId, req.Password)
	if err != nil {
		if err.Error() == "two factor authentication required" {
			err := s.twoFactorService.GenerateTwoFactorOTP(ctx, req.EmailOrId)
			if err != nil {
				return nil, fmt.Errorf("failed to generate two-factor OTP: %v", err)
			}
			return &authentication.LoginResponse{
				Success:          false,
				RequireTwoFactor: true,
				Message:          "Two-factor authentication required. An OTP has been sent to your email.",
			}, nil
		}
		if err.Error() == "forced password reset required" {
			return &authentication.LoginResponse{
				Success:              false,
				RequirePasswordReset: true,
				Message:              "A password reset is required for your account. Please check your email for instructions.",
			}, nil
		}
		return &authentication.LoginResponse{Success: false, Message: "Invalid email/ID or password"}, nil
	}

	return s.generateSuccessfulLoginResponse(user)
}

func (s *authServer) Logout(ctx context.Context, req *authentication.LogoutRequest) (*authentication.LogoutResponse, error) {
	// Validate the token
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userID := int32(claims["user_id"].(float64))
	if userID != req.UserID {
		return nil, fmt.Errorf("unauthorized: token does not match user ID")
	}

	// Invalidate the token
	utils.InvalidateToken(req.Token)

	return &authentication.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

func (s *authServer) GenerateTwoFactorOTP(ctx context.Context, req *authentication.GenerateTwoFactorOTPRequest) (*authentication.GenerateTwoFactorOTPResponse, error) {
	err := s.twoFactorService.GenerateTwoFactorOTP(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate two-factor OTP: %v", err)
	}

	return &authentication.GenerateTwoFactorOTPResponse{
		Success: true,
		Message: "Two-factor authentication OTP sent. Please check your email.",
	}, nil
}

func (s *authServer) VerifyTwoFactor(ctx context.Context, req *authentication.VerifyTwoFactorRequest) (*authentication.LoginResponse, error) {
	success, err := s.twoFactorService.VerifyTwoFactor(ctx, req.Email, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to verify two-factor authentication: %v", err)
	}

	if !success {
		return &authentication.LoginResponse{
			Success: false,
			Message: "Failed to verify two-factor authentication code.",
		}, nil
	}

	// If 2FA verification is successful, complete the login process
	user, err := s.loginService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information: %v", err)
	}

	return s.generateSuccessfulLoginResponse(user)
}

func (s *authServer) generateSuccessfulLoginResponse(user *models.User) (*authentication.LoginResponse, error) {
	if user.Status.Valid && user.Status.String == "pending" {
		token, err := utils.GenerateToken(user.Userid, user.Name, user.Userrole, user.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %v", err)
		}

		// Invalidate the token after 20 seconds
		go func() {
			time.Sleep(20 * time.Second)
			utils.InvalidateToken(token)
		}()

		return &authentication.LoginResponse{
			Success:  true,
			Token:    token,
			UserRole: user.Userrole,
			UserID:   user.Userid,
			Message:  "Your account is pending approval. You will be logged out in 20 seconds.",
			Status:   user.Status.String,
		}, nil
	}

	if user.Status.Valid && user.Status.String == "rejected" {
		return &authentication.LoginResponse{Success: false, Message: "Your account has been rejected."}, nil
	}

	token, err := utils.GenerateToken(user.Userid, user.Name, user.Userrole, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &authentication.LoginResponse{
		Success:  true,
		Token:    token,
		UserRole: user.Userrole,
		UserID:   user.Userid,
		Message:  "Login successful",
		Status:   user.Status.String,
	}, nil
}

func (s *authServer) RequestPasswordReset(ctx context.Context, req *authentication.PasswordResetRequest) (*authentication.PasswordResetResponse, error) {
	err := s.recoveryService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to request password reset: %v", err)
	}

	return &authentication.PasswordResetResponse{Success: true}, nil
}

func (s *authServer) ResetPassword(ctx context.Context, req *authentication.ResetPasswordRequest) (*authentication.ResetPasswordResponse, error) {
	err := s.recoveryService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to reset password: %v", err)
	}

	return &authentication.ResetPasswordResponse{Success: true}, nil
}

func (s *authServer) BeginWebAuthnRegistration(ctx context.Context, req *authentication.BeginWebAuthnRegistrationRequest) (*authentication.BeginWebAuthnRegistrationResponse, error) {
	// Verify the token
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userID := int32(claims["user_id"].(float64))
	if userID != req.UserID {
		return nil, fmt.Errorf("unauthorized: token does not match user ID")
	}

	options, err := s.biometricService.BeginRegistration(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to begin WebAuthn registration: %v", err)
	}

	return &authentication.BeginWebAuthnRegistrationResponse{
		Options: options,
	}, nil
}

func (s *authServer) FinishWebAuthnRegistration(ctx context.Context, req *authentication.FinishWebAuthnRegistrationRequest) (*authentication.FinishWebAuthnRegistrationResponse, error) {
	// Verify the token
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userID := int32(claims["user_id"].(float64))
	if userID != req.UserID {
		return nil, fmt.Errorf("unauthorized: token does not match user ID")
	}

	err = s.biometricService.FinishRegistration(ctx, req.UserID, req.Credential)
	if err != nil {
		return nil, fmt.Errorf("failed to finish WebAuthn registration: %v", err)
	}

	return &authentication.FinishWebAuthnRegistrationResponse{
		Success: true,
	}, nil
}

func (s *authServer) BeginWebAuthnLogin(ctx context.Context, req *authentication.BeginWebAuthnLoginRequest) (*authentication.BeginWebAuthnLoginResponse, error) {
	options, err := s.biometricService.BeginLogin(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to begin WebAuthn login: %v", err)
	}

	return &authentication.BeginWebAuthnLoginResponse{
		Options: options,
	}, nil
}

func (s *authServer) FinishWebAuthnLogin(ctx context.Context, req *authentication.FinishWebAuthnLoginRequest) (*authentication.FinishWebAuthnLoginResponse, error) {
	err := s.biometricService.FinishLogin(ctx, req.Email, req.Credential)
	if err != nil {
		return nil, fmt.Errorf("failed to finish WebAuthn login: %v", err)
	}

	// Get user information
	user, err := s.loginService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information: %v", err)
	}

	// Generate token
	token, err := utils.GenerateToken(user.Userid, user.Name, user.Userrole, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &authentication.FinishWebAuthnLoginResponse{
		Success: true,
		Token:   token,
	}, nil
}
