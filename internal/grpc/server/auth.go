package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	queries *models.Queries
	proto.UnimplementedAuthServiceServer
}

func NewAuthServer(queries *models.Queries) proto.AuthServiceServer {
	return &authServer{queries: queries}
}

func (s *authServer) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
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
		// Create a student record
		// NOTE: You'll need to update this based on the available fields in the proto message and the database schema
		// _, err = s.queries.CreateStudent(ctx, models.CreateStudentParams{
		//   Userid:   user.Userid,
		//   Name:     req.FirstName + " " + req.LastName,
		//   Grade:    "", // Set the appropriate grade value
		//   Schoolid: req.SchoolID,
		// })
	case "school":
		_, err = s.queries.CreateSchool(ctx, models.CreateSchoolParams{
			Name:            req.SchoolName,
			Address:         "", // Set the appropriate address value
			Contactpersonid: user.Userid,
			Contactemail:    req.ContactEmail,
			Category:        req.SchoolType,
		})
	case "volunteer":
		_, err = s.queries.CreateVolunteer(ctx, models.CreateVolunteerParams{
			Name:   req.FirstName + " " + req.LastName,
			Role:   req.RoleInterestedIn,
			Userid: user.Userid,
		})
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

	return &proto.SignUpResponse{Success: true, Message: "Sign-up successful"}, nil
}

func (s *authServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// Retrieve the user based on the provided username (email)
	user, err := s.queries.GetUserByEmail(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate a PASETO token
	token, err := generateToken(user.Userid, user.Userrole)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &proto.LoginResponse{Success: true, Token: token}, nil
}

func generateToken(userID int32, userRole string) (string, error) {
	// Generate a PASETO secret key
	secretKey, err := generatePasetoSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate PASETO secret key: %v", err)
	}

	// Create a new PASETO maker with the generated secret key
	maker := paseto.NewV2()

	// Set the token claims
	claims := map[string]interface{}{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Generate and return the token
	token, err := maker.Sign([]byte(secretKey), claims, nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}

func generatePasetoSecretKey() (string, error) {
	// Generate a 32-byte random key
	keyBytes := make([]byte, 32)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random key: %v", err)
	}

	// Encode the key using base64
	keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)

	return keyBase64, nil
}
