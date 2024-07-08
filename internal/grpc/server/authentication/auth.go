package server

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/services"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	authentication.UnimplementedAuthServiceServer
	queries           *models.Queries
	twoFactorService  *services.TwoFactorService
	recoveryService   *services.RecoveryService
	biometricService  *services.BiometricService
	privateKey        ed25519.PrivateKey
}

func NewAuthServer(queries *models.Queries) (authentication.AuthServiceServer, error) {
    privateKey, publicKey, err := utils.GeneratePasetoKeyPair()
    if err != nil {
        return nil, fmt.Errorf("failed to generate PASETO key pair: %v", err)
    }

    // Set the public key for token validation
    utils.SetPublicKey(publicKey)

    return &authServer{
        queries:           queries,
        twoFactorService:  services.NewTwoFactorService(queries),
        recoveryService:   services.NewRecoveryService(queries),
        biometricService:  services.NewBiometricService(queries),
        privateKey:        privateKey,
    }, nil
}


func (s *authServer) SignUp(ctx context.Context, req *authentication.SignUpRequest) (*authentication.SignUpResponse, error) {
	// Validate input
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.UserRole == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create a new user
	user, err := s.queries.CreateUser(ctx, models.CreateUserParams{
		Name:     req.FirstName + " " + req.LastName,
		Email:    req.Email,
		Password: string(hashedPassword),
		Userrole: req.UserRole,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Create a student, school, or volunteer record based on the user role
	switch req.UserRole {
	case "student":
		// Parse date of birth
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %v", err)
		}

		_, err = s.queries.CreateStudent(ctx, models.CreateStudentParams{
			Firstname:   req.FirstName,
			Lastname:    req.LastName,
			Grade:       "",
			Dateofbirth: sql.NullTime{Time: dateOfBirth, Valid: !dateOfBirth.IsZero()},
			Email:       sql.NullString{String: req.Email, Valid: req.Email != ""},
			Password:    string(hashedPassword),
			Schoolid:    req.SchoolID,
			Userid:      user.Userid,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create student record: %v", err)
		}
	case "school":
		_, err = s.queries.CreateSchool(ctx, models.CreateSchoolParams{
			Schoolname:      req.SchoolName,
			Address:         "",
			Country:         sql.NullString{String: req.Country, Valid: req.Country != ""},
			Province:        sql.NullString{String: req.Province, Valid: req.Province != ""},
			District:        sql.NullString{String: req.District, Valid: req.District != ""},
			Contactpersonid: user.Userid,
			Contactemail:    req.ContactEmail,
			Schooltype:      req.SchoolType,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create school record: %v", err)
		}
	case "volunteer":
		// Parse date of birth
		dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %v", err)
		}

		_, err = s.queries.CreateVolunteer(ctx, models.CreateVolunteerParams{
			Firstname:            req.FirstName,
			Lastname:             req.LastName,
			Dateofbirth:          sql.NullTime{Time: dateOfBirth, Valid: !dateOfBirth.IsZero()},
			Role:                 req.RoleInterestedIn,
			Graduateyear:         sql.NullInt32{Int32: int32(req.GraduationYear), Valid: req.GraduationYear != 0},
			Password:             string(hashedPassword),
			Safeguardcertificate: sql.NullBool{Bool: len(req.SafeguardingCertificate) > 0, Valid: true},
			Userid:               user.Userid,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create volunteer record: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid user role")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create user-specific record: %v", err)
	}

	// Send welcome email
	err = utils.SendWelcomeEmail(req.Email, req.FirstName)
	if err != nil {
		// Log the error, but don't fail the sign-up process
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}

	return &authentication.SignUpResponse{Success: true, Message: "Sign-up successful"}, nil
}

func (s *authServer) sendForcedPasswordReset(ctx context.Context, email string) error {
    err := s.recoveryService.RequestPasswordReset(ctx, email)
    if err != nil {
        return fmt.Errorf("failed to initiate forced password reset: %v", err)
    }
    return nil
}

func (s *authServer) Login(ctx context.Context, req *authentication.LoginRequest) (*authentication.LoginResponse, error) {
    // Validate input
    if req.Email == "" || req.Password == "" {
        return nil, fmt.Errorf("missing required fields")
    }

    // Retrieve the user based on the provided email
    user, err := s.queries.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve user: %v", err)
    }

    // Update last login attempt
    err = s.queries.UpdateLastLoginAttempt(ctx, user.Userid)
    if err != nil {
        return nil, fmt.Errorf("failed to update last login attempt: %v", err)
    }

    // Verify the password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        // Increment failed login attempts
        err = s.queries.IncrementFailedLoginAttempts(ctx, user.Userid)
        if err != nil {
            return nil, fmt.Errorf("failed to update login attempts: %v", err)
        }

        // Check if we need to enforce 2FA or password reset
        if user.FailedLoginAttempts.Int32 >= 4 { // Now 5 total attempts including this one
            if user.TwoFactorEnabled.Valid && user.TwoFactorEnabled.Bool {
                return &authentication.LoginResponse{Success: false, RequireTwoFactor: true}, nil
            } else {
                // Automatically send password reset email
                err = s.sendForcedPasswordReset(ctx, user.Email)
                if err != nil {
                    // Log the error, but don't expose it to the user
                    fmt.Printf("Failed to send forced password reset email: %v\n", err)
                }
                return &authentication.LoginResponse{Success: false, RequirePasswordReset: true, Message: "A password reset email has been sent to your account."}, nil
            }
        }

        return nil, fmt.Errorf("invalid password")
    }

    // Reset failed login attempts on successful login
    err = s.queries.ResetFailedLoginAttempts(ctx, user.Userid)
    if err != nil {
        return nil, fmt.Errorf("failed to reset login attempts: %v", err)
    }

    // Generate a PASETO token
    token, err := utils.GenerateToken(user.Userid, user.Userrole, s.privateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %v", err)
    }

    return &authentication.LoginResponse{Success: true, Token: token, UserRole: user.Userrole, UserID: user.Userid}, nil
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
    user, err := s.queries.GetUserByID(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %v", err)
    }

    if !user.TwoFactorSecret.Valid || !s.twoFactorService.ValidateCode(user.TwoFactorSecret.String, req.Code) {
        return nil, fmt.Errorf("invalid two-factor code")
    }

    // Actually enable two-factor authentication
    err = s.queries.EnableTwoFactor(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to enable two-factor authentication: %v", err)
    }

    return &authentication.VerifyTwoFactorResponse{Success: true}, nil
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

    err = s.queries.DisableTwoFactor(ctx, req.UserID)
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
        // Update last login attempt
        updateErr := s.queries.UpdateLastLoginAttempt(ctx, user.Userid)
        if updateErr != nil {
            return nil, fmt.Errorf("failed to update last login attempt: %v", updateErr)
        }

        // Increment failed login attempts
        incrementErr := s.queries.IncrementFailedLoginAttempts(ctx, user.Userid)
        if incrementErr != nil {
            return nil, fmt.Errorf("failed to update login attempts: %v", incrementErr)
        }

         // Check if we need to enforce 2FA or password reset
		 if user.FailedLoginAttempts.Int32 >= 4 { // Now 5 total attempts including this one
            if user.TwoFactorEnabled.Valid && user.TwoFactorEnabled.Bool {
                return &authentication.LoginResponse{Success: false, RequireTwoFactor: true}, nil
            } else {
                // Automatically send password reset email
                err = s.sendForcedPasswordReset(ctx, user.Email)
                if err != nil {
                    // Log the error, but don't expose it to the user
                    fmt.Printf("Failed to send forced password reset email: %v\n", err)
                }
                return &authentication.LoginResponse{Success: false, RequirePasswordReset: true, Message: "A password reset email has been sent to your account."}, nil
            }
        }

        return nil, fmt.Errorf("failed to verify biometric token: %v", err)
    }

    // Update last login attempt for successful login
    err = s.queries.UpdateLastLoginAttempt(ctx, user.Userid)
    if err != nil {
        return nil, fmt.Errorf("failed to update last login attempt: %v", err)
    }

    // Reset failed login attempts on successful login
    err = s.queries.ResetFailedLoginAttempts(ctx, user.Userid)
    if err != nil {
        return nil, fmt.Errorf("failed to reset login attempts: %v", err)
    }

    token, err := utils.GenerateToken(user.Userid, user.Userrole, s.privateKey)
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %v", err)
    }

    return &authentication.LoginResponse{
        Success:  true,
        Token:    token,
        UserRole: user.Userrole,
        UserID:   user.Userid,
    }, nil
}