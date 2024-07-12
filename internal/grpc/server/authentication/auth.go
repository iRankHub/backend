package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	services "github.com/iRankHub/backend/internal/services/authentication"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	authentication.UnimplementedAuthServiceServer
	db               *sql.DB
	loginService     *services.LoginService
	signUpService    *services.SignUpService
	twoFactorService *services.TwoFactorService
	recoveryService  *services.RecoveryService
	biometricService *services.BiometricService
}

func NewAuthServer(db *sql.DB) (authentication.AuthServiceServer, error) {
	twoFactorService := services.NewTwoFactorService(db)
	recoveryService := services.NewRecoveryService(db)
	biometricService := services.NewBiometricService(db)
	loginService := services.NewLoginService(db, twoFactorService, recoveryService)
	signUpService := services.NewSignUpService(db)

	return &authServer{
		db:               db,
		loginService:     loginService,
		signUpService:    signUpService,
		twoFactorService: twoFactorService,
		recoveryService:  recoveryService,
		biometricService: biometricService,
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
	}

	err := s.signUpService.SignUp(ctx, req.FirstName, req.LastName, req.Email, req.Password, req.UserRole, additionalInfo)
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

func (s *authServer) Login(ctx context.Context, req *authentication.LoginRequest) (*authentication.LoginResponse, error) {
	user, err := s.loginService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err.Error() == "two factor authentication required" {
			return &authentication.LoginResponse{Success: false, RequireTwoFactor: true}, nil
		}
		if err.Error() == "password reset required" {
			return &authentication.LoginResponse{Success: false, RequirePasswordReset: true, Message: "A password reset email has been sent to your account."}, nil
		}
		return &authentication.LoginResponse{Success: false, Message: "Invalid email or password"}, nil
	}

	if user.Status.Valid && user.Status.String == "pending" {
		token, err := utils.GenerateToken(user.Userid, user.Userrole, user.Email)
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

	token, err := utils.GenerateToken(user.Userid, user.Userrole, user.Email)
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

func (s *authServer) EnableTwoFactor(ctx context.Context, req *authentication.EnableTwoFactorRequest) (*authentication.EnableTwoFactorResponse, error) {
    // Verify the token
    claims, err := utils.ValidateToken(req.Token)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %v", err)
    }

    userID := int32(claims["user_id"].(float64))
    if userID != req.UserID {
        return nil, fmt.Errorf("unauthorized: token does not match user ID")
    }

    secret, qrCode, err := s.twoFactorService.EnableTwoFactor(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to enable two-factor authentication: %v", err)
    }

    return &authentication.EnableTwoFactorResponse{
        Secret: secret,
        QrCode: qrCode,
    }, nil
}

func (s *authServer) VerifyTwoFactor(ctx context.Context, req *authentication.VerifyTwoFactorRequest) (*authentication.VerifyTwoFactorResponse, error) {
	// Validate the token
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	// Check if the token belongs to the user
	userID := int32(claims["user_id"].(float64))
	if userID != req.UserID {
		return nil, fmt.Errorf("token does not belong to the user")
	}

	success, err := s.twoFactorService.VerifyAndEnableTwoFactor(ctx, req.UserID, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to verify and enable two-factor authentication: %v", err)
	}

	return &authentication.VerifyTwoFactorResponse{Success: success}, nil
}

func (s *authServer) DisableTwoFactor(ctx context.Context, req *authentication.DisableTwoFactorRequest) (*authentication.DisableTwoFactorResponse, error) {
    // Verify the token
    claims, err := utils.ValidateToken(req.Token)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %v", err)
    }

    userID := int32(claims["user_id"].(float64))
    if userID != req.UserID {
        return nil, fmt.Errorf("unauthorized: token does not match user ID")
    }

    err = s.twoFactorService.DisableTwoFactor(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to disable two-factor authentication: %v", err)
    }

    return &authentication.DisableTwoFactorResponse{Success: true}, nil
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

func (s *authServer) EnableBiometricLogin(ctx context.Context, req *authentication.EnableBiometricLoginRequest) (*authentication.EnableBiometricLoginResponse, error) {
	token, err := s.biometricService.EnableBiometricLogin(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to enable biometric login: %v", err)
	}

	return &authentication.EnableBiometricLoginResponse{BiometricToken: token}, nil
}

func (s *authServer) BiometricLogin(ctx context.Context, req *authentication.BiometricLoginRequest) (*authentication.LoginResponse, error) {
	user, err := s.biometricService.VerifyBiometricToken(ctx, req.BiometricToken)
	if err != nil {
		// If VerifyBiometricToken fails, we don't have a user to pass to HandleFailedLoginAttempt
		// We should handle this case differently
		return nil, fmt.Errorf("failed to verify biometric token: %v", err)
	}

	err = s.loginService.HandleFailedLoginAttempt(ctx, user)
	if err != nil {
		if err.Error() == "two factor authentication required" {
			return &authentication.LoginResponse{Success: false, RequireTwoFactor: true}, nil
		}
		if err.Error() == "password reset required" {
			return &authentication.LoginResponse{Success: false, RequirePasswordReset: true, Message: "A password reset email has been sent to your account."}, nil
		}
		return nil, fmt.Errorf("login attempt failed: %v", err)
	}

	err = s.loginService.HandleSuccessfulLogin(ctx, user.Userid)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.Userid, user.Userrole, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &authentication.LoginResponse{
		Success:  true,
		Token:    token,
		UserRole: user.Userrole,
		UserID:   user.Userid,
		Message:  "Login Successful",
	}, nil
}
