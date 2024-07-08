package services

import (
	"context"
	"crypto/ed25519"
	"fmt"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type AuthService struct {
	queries    *models.Queries
	privateKey ed25519.PrivateKey
}

func NewAuthService(queries *models.Queries, privateKey ed25519.PrivateKey) *AuthService {
	return &AuthService{queries: queries, privateKey: privateKey}
}

func (s *AuthService) Login(ctx context.Context, req *authentication.LoginRequest) (*authentication.LoginResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// Retrieve the user based on the provided email
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	// Verify the password
	err = utils.ComparePasswords(user.Password, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate a PASETO token
	token, err := utils.GenerateToken(user.Userid, user.Userrole, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &authentication.LoginResponse{Success: true, Token: token, UserRole: user.Userrole, UserID: user.Userid}, nil
}