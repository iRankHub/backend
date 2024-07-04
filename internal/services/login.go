package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type AuthService struct {
	queries *models.Queries
}

func NewAuthService(queries *models.Queries) *AuthService {
	return &AuthService{queries: queries}
}

func (s *AuthService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate an Ed25519 key pair
	privateKey, _, err := utils.GeneratePasetoKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PASETO key pair: %v", err)
	}

	// Generate a PASETO token
	token, err := utils.GenerateToken(user.Userid, user.Userrole, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &proto.LoginResponse{Success: true, Token: token, UserRole: user.Userrole, UserID: user.Userid}, nil
}