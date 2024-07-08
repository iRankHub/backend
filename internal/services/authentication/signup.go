package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type SignUpService struct {
	db *sql.DB
}

func NewSignUpService(db *sql.DB) *SignUpService {
	return &SignUpService{
		db: db,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, firstName, lastName, email, password, userRole string, additionalInfo map[string]interface{}) error {
	if firstName == "" || lastName == "" || email == "" || password == "" || userRole == "" {
		return fmt.Errorf("missing required fields")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user, err := queries.CreateUser(ctx, models.CreateUserParams{
		Name:     firstName + " " + lastName,
		Email:    email,
		Password: hashedPassword,
		Userrole: userRole,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	switch userRole {
	case "student":
		err = s.createStudentRecord(ctx, queries, user.Userid, firstName, lastName, email, hashedPassword, additionalInfo)
	case "school":
		err = s.createSchoolRecord(ctx, queries, user.Userid, additionalInfo)
	case "volunteer":
		err = s.createVolunteerRecord(ctx, queries, user.Userid, firstName, lastName, hashedPassword, additionalInfo)
	default:
		return fmt.Errorf("invalid user role")
	}

	if err != nil {
		return fmt.Errorf("failed to create user-specific record: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	err = utils.SendWelcomeEmail(email, firstName)
	if err != nil {
		fmt.Printf("Failed to send welcome email: %v\n", err)
	}

	return nil
}

func (s *SignUpService) createStudentRecord(ctx context.Context, queries *models.Queries, userID int32, firstName, lastName, email string, hashedPassword string, additionalInfo map[string]interface{}) error {
	dateOfBirthStr, ok := additionalInfo["dateOfBirth"].(string)
	if !ok || dateOfBirthStr == "" {
		return fmt.Errorf("date of birth is missing or invalid")
	}
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return fmt.Errorf("invalid date of birth format: %v", err)
	}

	schoolID, ok := additionalInfo["schoolID"].(int32)
	if !ok || schoolID == 0 {
		return fmt.Errorf("school ID is missing or invalid")
	}

	grade, ok := additionalInfo["grade"].(string)
	if !ok || grade == "" {
		return fmt.Errorf("grade is missing or invalid")
	}

	_, err = queries.CreateStudent(ctx, models.CreateStudentParams{
		Firstname:   firstName,
		Lastname:    lastName,
		Grade:       grade,
		Dateofbirth: sql.NullTime{Time: dateOfBirth, Valid: true},
		Email:       sql.NullString{String: email, Valid: true},
		Password:    hashedPassword,
		Schoolid:    schoolID,
		Userid:      userID,
	})
	if err != nil {
		return fmt.Errorf("failed to create student record: %v", err)
	}
	return nil
}

func (s *SignUpService) createSchoolRecord(ctx context.Context, queries *models.Queries, userID int32, additionalInfo map[string]interface{}) error {

	schoolName, ok := additionalInfo["schoolName"].(string)
	if !ok || schoolName == "" {
		return fmt.Errorf("school name is missing or invalid")
	}

	address, ok := additionalInfo["address"].(string)
	if !ok || address == "" {
		return fmt.Errorf("address is missing or invalid")
	}

	country, ok := additionalInfo["country"].(string)
	if !ok || country == "" {
		return fmt.Errorf("country is missing or invalid")
	}

	province, ok := additionalInfo["province"].(string)
	if !ok || province == "" {
		return fmt.Errorf("province is missing or invalid")
	}

	district, ok := additionalInfo["district"].(string)
	if !ok || district == "" {
		return fmt.Errorf("district is missing or invalid")
	}

	contactEmail, ok := additionalInfo["contactEmail"].(string)
	if !ok || contactEmail == "" {
		return fmt.Errorf("contact email is missing or invalid")
	}

	schoolType, ok := additionalInfo["schoolType"].(string)
	if !ok || schoolType == "" {
		return fmt.Errorf("school type is missing or invalid")
	}

	_, err := queries.CreateSchool(ctx, models.CreateSchoolParams{
		Schoolname:      schoolName,
		Address:         address,
		Country:         sql.NullString{String: country, Valid: true},
		Province:        sql.NullString{String: province, Valid: true},
		District:        sql.NullString{String: district, Valid: true},
		Contactpersonid: userID,
		Contactemail:    contactEmail,
		Schooltype:      schoolType,
	})
	return err
}

func (s *SignUpService) createVolunteerRecord(ctx context.Context, queries *models.Queries, userID int32, firstName, lastName, hashedPassword string, additionalInfo map[string]interface{}) error {
	dateOfBirthStr, ok := additionalInfo["dateOfBirth"].(string)
	if !ok || dateOfBirthStr == "" {
		return fmt.Errorf("date of birth is missing or invalid")
	}
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return fmt.Errorf("invalid date of birth format: %v", err)
	}

	graduationYear, ok := additionalInfo["graduationYear"].(int32)
	if !ok || graduationYear == 0 {
		return fmt.Errorf("graduation year is missing or invalid")
	}

	roleInterestedIn, ok := additionalInfo["roleInterestedIn"].(string)
	if !ok || roleInterestedIn == "" {
		return fmt.Errorf("role interested in is missing or invalid")
	}

	safeguardingCertificate, ok := additionalInfo["safeguardingCertificate"].(bool)
	if !ok {
		return fmt.Errorf("safeguarding certificate is missing or invalid")
	}

	_, err = queries.CreateVolunteer(ctx, models.CreateVolunteerParams{
		Firstname:            firstName,
		Lastname:             lastName,
		Dateofbirth:          sql.NullTime{Time: dateOfBirth, Valid: true},
		Role:                 roleInterestedIn,
		Graduateyear:         sql.NullInt32{Int32: graduationYear, Valid: true},
		Password:             hashedPassword,
		Safeguardcertificate: sql.NullBool{Bool: safeguardingCertificate, Valid: true},
		Userid:               userID,
	})
	return err
}