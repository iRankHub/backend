package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/services/notification"
	notificationModels "github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/utils"
)

type SignUpService struct {
	db                  *sql.DB
	notificationService *notification.Service
}

func NewSignUpService(db *sql.DB, ns *notification.Service) *SignUpService {
	return &SignUpService{
		db:                  db,
		notificationService: ns,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, firstName, lastName, email, password, userRole, gender string, nationalID string, safeguardingCertificateURL string, additionalInfo map[string]interface{}) error {
	if firstName == "" || lastName == "" || email == "" || password == "" || userRole == "" {
		return fmt.Errorf("missing required fields")
	}

	if gender != "male" && gender != "female" && gender != "non-binary" {
		return fmt.Errorf("invalid gender. Must be 'male', 'female', or 'non-binary'")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Check if email is unique
	_, err = queries.GetUserByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("email already in use")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking email uniqueness: %v", err)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	status := sql.NullString{String: "pending", Valid: true}
	if userRole == "admin" {
		status = sql.NullString{String: "approved", Valid: true}
	}

	user, err := queries.CreateUser(ctx, models.CreateUserParams{
		Name:     firstName + " " + lastName,
		Email:    email,
		Password: hashedPassword,
		Userrole: userRole,
		Status:   status,
		Gender:   sql.NullString{String: gender, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Create role-specific record
	switch userRole {
	case "student":
		err = s.createStudentRecord(ctx, queries, user.Userid, firstName, lastName, email, gender, hashedPassword, additionalInfo)
	case "school":
		err = s.createSchoolRecord(ctx, queries, user.Userid, email, nationalID, additionalInfo)
	case "volunteer":
		err = s.createVolunteerRecord(ctx, queries, user.Userid, firstName, lastName, gender, hashedPassword, nationalID, safeguardingCertificateURL, additionalInfo)
	case "admin":
		err = s.createAdminProfile(ctx, queries, user)
	default:
		return fmt.Errorf("invalid user role")
	}

	if err != nil {
		return fmt.Errorf("failed to create user-specific record: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Get client metadata for notifications
	clientMeta := utils.FromContext(ctx)

	// Send appropriate notifications based on role
	if userRole == "admin" {
		metadata := notificationModels.AuthMetadata{
			DeviceInfo:   clientMeta.DeviceInfo,
			Location:     "Admin Account Creation",
			LastAttempt:  time.Now(),
			AttemptCount: 0,
			IPAddress:    clientMeta.IP,
		}

		err = s.notificationService.SendAccountCreation(
			ctx,
			email,
			notificationModels.AdminRole,
			metadata,
		)
	} else {
		// Send welcome email to user
		metadata := notificationModels.AuthMetadata{
			DeviceInfo:   clientMeta.DeviceInfo,
			Location:     "Account Creation",
			LastAttempt:  time.Now(),
			AttemptCount: 0,
			IPAddress:    clientMeta.IP,
		}

		err = s.notificationService.SendAccountCreation(
			ctx,
			email,
			s.getUserRole(userRole),
			metadata,
		)

		if err != nil {
			return fmt.Errorf("failed to send welcome notification: %v", err)
		}

		// Notify admins about new signup
		err = s.notifyAdminsOfNewSignup(ctx, user.Userid, userRole, clientMeta)
	}

	if err != nil {
		return fmt.Errorf("failed to send notifications: %v", err)
	}

	return nil
}

func (s *SignUpService) notifyAdminsOfNewSignup(ctx context.Context, userID int32, userRole string, clientMeta utils.ClientMetadata) error {
	metadata := notificationModels.AuthMetadata{
		DeviceInfo:   clientMeta.DeviceInfo,
		Location:     "New User Registration",
		LastAttempt:  time.Now(),
		AttemptCount: 0,
		IPAddress:    clientMeta.IP,
	}

	// Query for admin users
	queries := models.New(s.db)
	adminUsers, err := queries.GetVolunteersAndAdmins(ctx, models.GetVolunteersAndAdminsParams{
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to get admin users: %v", err)
	}

	// Notify each admin
	for _, admin := range adminUsers {
		if admin.Userrole == "admin" {
			err = s.notificationService.SendSecurityAlert(
				ctx,
				admin.Email,
				notificationModels.AdminRole,
				fmt.Sprintf("New %s user registration requires approval. IP: %s", userRole, clientMeta.IP),
				metadata,
			)
			if err != nil {
				return fmt.Errorf("failed to notify admin: %v", err)
			}
		}
	}

	return nil
}

func (s *SignUpService) createStudentRecord(ctx context.Context, queries *models.Queries, userID int32, firstName, lastName, email, gender, hashedPassword string, additionalInfo map[string]interface{}) error {
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
		Gender:      sql.NullString{String: gender, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create student record: %v", err)
	}
	return nil
}

func (s *SignUpService) createSchoolRecord(ctx context.Context, queries *models.Queries, userID int32, email, nationalID string, additionalInfo map[string]interface{}) error {
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
		Schoolname:              schoolName,
		Address:                 address,
		Country:                 sql.NullString{String: country, Valid: true},
		Province:                sql.NullString{String: province, Valid: true},
		District:                sql.NullString{String: district, Valid: true},
		Contactpersonid:         userID,
		Contactemail:            contactEmail,
		Schoolemail:             email,
		Schooltype:              schoolType,
		Contactpersonnationalid: sql.NullString{String: nationalID, Valid: true},
	})
	return err
}

func (s *SignUpService) createVolunteerRecord(ctx context.Context, queries *models.Queries, userID int32, firstName, lastName, gender, hashedPassword, nationalID string, safeguardingCertificateURL string, additionalInfo map[string]interface{}) error {
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

	hasInternship, ok := additionalInfo["hasInternship"].(bool)
	if !ok {
		return fmt.Errorf("internship information is missing or invalid")
	}

	isEnrolledInUniversity, ok := additionalInfo["isEnrolledInUniversity"].(bool)
	if !ok {
		return fmt.Errorf("university enrollment information is missing or invalid")
	}

	_, err = queries.CreateVolunteer(ctx, models.CreateVolunteerParams{
		Firstname:              firstName,
		Lastname:               lastName,
		Dateofbirth:            sql.NullTime{Time: dateOfBirth, Valid: true},
		Role:                   roleInterestedIn,
		Graduateyear:           sql.NullInt32{Int32: graduationYear, Valid: true},
		Password:               hashedPassword,
		Safeguardcertificate:   sql.NullString{String: safeguardingCertificateURL, Valid: safeguardingCertificateURL != ""},
		Hasinternship:          sql.NullBool{Bool: hasInternship, Valid: true},
		Userid:                 userID,
		Isenrolledinuniversity: sql.NullBool{Bool: isEnrolledInUniversity, Valid: true},
		Gender:                 sql.NullString{String: gender, Valid: true},
		Nationalid:             sql.NullString{String: nationalID, Valid: true},
	})
	return err
}

func (s *SignUpService) createAdminProfile(ctx context.Context, queries *models.Queries, user models.User) error {
	_, err := queries.CreateUserProfile(ctx, models.CreateUserProfileParams{
		Userid:             user.Userid,
		Name:               user.Name,
		Userrole:           user.Userrole,
		Email:              user.Email,
		Password:           user.Password,
		Verificationstatus: user.Verificationstatus,
		Address:            sql.NullString{},
		Phone:              sql.NullString{},
		Bio:                sql.NullString{},
		Profilepicture:     sql.NullString{},
		Gender:             user.Gender,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin user profile: %v", err)
	}
	return nil
}

func (s *SignUpService) getUserRole(role string) notificationModels.UserRole {
	switch role {
	case "admin":
		return notificationModels.AdminRole
	case "school":
		return notificationModels.SchoolRole
	case "student":
		return notificationModels.StudentRole
	case "volunteer":
		return notificationModels.VolunteerRole
	default:
		return notificationModels.UnspecifiedRole
	}
}
