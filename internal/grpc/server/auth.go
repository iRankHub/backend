package server

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type authServer struct {
	proto.UnimplementedAuthServiceServer
	queries    *models.Queries
	privateKey ed25519.PrivateKey
}

func NewAuthServer(queries *models.Queries) (proto.AuthServiceServer, error) {
	privateKey, _, err := utils.GeneratePasetoKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PASETO key pair: %v", err)
	}
	return &authServer{queries: queries, privateKey: privateKey}, nil
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

	return &proto.SignUpResponse{Success: true, Message: "Sign-up successful"}, nil
}

func (s *authServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
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

	// Generate a PASETO token
	token, err := utils.GenerateToken(user.Userid, user.Userrole, s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &proto.LoginResponse{Success: true, Token: token, UserRole: user.Userrole, UserID: user.Userid}, nil
}