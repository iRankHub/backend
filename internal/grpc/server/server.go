package server

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
)

type authServer struct {
  db *sql.DB
  proto.UnimplementedAuthServiceServer
}

func NewAuthServer(db *sql.DB) proto.AuthServiceServer {
  return &authServer{db: db}
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

  // Start a transaction
  tx, err := s.db.Begin()
  if err != nil {
    return nil, fmt.Errorf("failed to start transaction: %v", err)
  }
  defer tx.Rollback()

  // Create a new user
  userID, err := createUser(tx, req.FirstName, req.LastName, req.Email, string(hashedPassword), req.UserRole)
  if err != nil {
    return nil, fmt.Errorf("failed to create user: %v", err)
  }

  // Create a student, school, or volunteer record based on the user role
  switch req.UserRole {
  case "student":
    err = createStudent(tx, userID, req.DateOfBirth, req.SchoolID)
  case "school":
    err = createSchool(tx, userID, req.SchoolName, req.Country, req.Province, req.District, req.SchoolType, req.ContactPersonName, req.ContactPersonNumber, req.ContactEmail)
  case "volunteer":
    err = createVolunteer(tx, userID, req.DateOfBirth, req.NationalID, req.SchoolAttended, req.GraduationYear, req.RoleInterestedIn, req.SafeguardingCertificate)
  default:
    return nil, fmt.Errorf("invalid user role")
  }
  if err != nil {
    return nil, fmt.Errorf("failed to create user-specific record: %v", err)
  }

  // Commit the transaction
  err = tx.Commit()
  if err != nil {
    return nil, fmt.Errorf("failed to commit transaction: %v", err)
  }

  return &proto.SignUpResponse{Success: true, Message: "Sign-up successful"}, nil
}

func (s *authServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// Retrieve the user based on the provided username
	user, err := getUser(s.db, req.Username)
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

func createUser(tx *sql.Tx, firstName, lastName, email, hashedPassword, userRole string) (int32, error) {
  // Insert a new user into the Users table
  var userID int32
  err := tx.QueryRow(`
    INSERT INTO Users (FirstName, LastName, Email, Password, UserRole)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING UserID
  `, firstName, lastName, email, hashedPassword, userRole).Scan(&userID)
  if err != nil {
    return 0, fmt.Errorf("failed to insert user: %v", err)
  }

  return userID, nil
}

func createStudent(tx *sql.Tx, userID int32, dateOfBirth string, schoolID int32) error {
  // Insert a new student into the Students table
  _, err := tx.Exec(`
    INSERT INTO Students (UserID, DateOfBirth, SchoolID, UniqueStudentID)
    VALUES ($1, $2, $3, $4)
  `, userID, dateOfBirth, schoolID, generateUniqueID())
  if err != nil {
    return fmt.Errorf("failed to insert student: %v", err)
  }

  return nil
}

func createSchool(tx *sql.Tx, userID int32, name, country, province, district, schoolType, contactPersonName, contactPersonNumber, contactEmail string) error {
  // Insert a new school into the Schools table
  _, err := tx.Exec(`
    INSERT INTO Schools (UserID, Name, Country, Province, District, SchoolType, ContactPersonName, ContactPersonNumber, ContactEmail, UniqueSchoolID)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
  `, userID, name, country, province, district, schoolType, contactPersonName, contactPersonNumber, contactEmail, generateUniqueID())
  if err != nil {
    return fmt.Errorf("failed to insert school: %v", err)
  }

  return nil
}

func createVolunteer(tx *sql.Tx, userID int32, dateOfBirth, nationalID, schoolAttended string, graduationYear int32, roleInterestedIn string, safeguardingCertificate []byte) error {
  // Insert a new volunteer into the Volunteers table
  _, err := tx.Exec(`
    INSERT INTO Volunteers (UserID, DateOfBirth, NationalID, SchoolAttended, GraduationYear, RoleInterestedIn, SafeguardingCertificate, UniqueVolunteerID)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `, userID, dateOfBirth, nationalID, schoolAttended, graduationYear, roleInterestedIn, safeguardingCertificate, generateUniqueID())
  if err != nil {
    return fmt.Errorf("failed to insert volunteer: %v", err)
  }

  return nil
}

func getUser(db *sql.DB, username string) (*models.User, error) {
  // Retrieve the user from the Users table based on the provided username
  var user models.User
  err := db.QueryRow(`
    SELECT UserID, FirstName, LastName, Email, Password, UserRole
    FROM Users
    WHERE Email = $1 OR UserID = (
      SELECT UserID FROM Students WHERE UniqueStudentID = $1
      UNION
      SELECT UserID FROM Schools WHERE UniqueSchoolID = $1
      UNION
      SELECT UserID FROM Volunteers WHERE UniqueVolunteerID = $1
    )
  `, username).Scan(&user.Userid, &user.Name, &user.Email, &user.Password, &user.Userrole)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, fmt.Errorf("user not found")
    }
    return nil, fmt.Errorf("failed to retrieve user: %v", err)
  }

  return &user, nil
}

func generateUniqueID() string {
  // Generate a unique 8-digit ID
  // TODO: Implement a proper unique ID generation logic
  return "12345678"
}