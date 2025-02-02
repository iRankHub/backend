package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/pquerna/otp/totp"

	"github.com/iRankHub/backend/internal/grpc/proto/user_management"
	"github.com/iRankHub/backend/internal/models"
	notificationService "github.com/iRankHub/backend/internal/services/notification"
	"github.com/iRankHub/backend/internal/utils"
	notifications "github.com/iRankHub/backend/internal/utils/notifications"
)

type UserManagementService struct {
	db                  *sql.DB
	notificationService *notificationService.NotificationService
	s3Client            *utils.S3Client
}

func NewUserManagementService(db *sql.DB) (*UserManagementService, error) {
	ns, err := notificationService.NewNotificationService(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification service: %v", err)
	}

	s3Client, err := utils.NewS3Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %v", err)
	}

	return &UserManagementService{
		db:                  db,
		notificationService: ns,
		s3Client:            s3Client,
	}, nil
}

func (s *UserManagementService) GetPendingUsers(ctx context.Context, token string) ([]models.User, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("only admins can get pending users")
	}

	queries := models.New(s.db)
	return queries.GetUsersByStatus(ctx, sql.NullString{String: "pending", Valid: true})
}

func (s *UserManagementService) GetAllUsers(ctx context.Context, token string, page, pageSize int32, searchQuery string) ([]models.GetAllUsersRow, int32, int32, int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, 0, 0, 0, fmt.Errorf("only admins can get all users")
	}

	queries := models.New(s.db)
	users, err := queries.GetAllUsers(ctx, models.GetAllUsersParams{
		Limit:       pageSize,
		Offset:      (page - 1) * pageSize,
		SearchQuery: searchQuery,
	})
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get users: %v", err)
	}

	totalCount, err := queries.GetTotalUserCountWithSearch(ctx, searchQuery)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("failed to get total user count: %v", err)
	}

	if totalCount > math.MaxInt32 {
		return nil, 0, 0, 0, fmt.Errorf("total user count exceeds maximum value for int32")
	}

	var approvedUsersCount, recentSignupsCount int32
	if len(users) > 0 {
		approvedUsersCount = int32(users[0].ApprovedUsersCount)
		recentSignupsCount = int32(users[0].RecentSignupsCount)
	}

	return users, int32(totalCount), approvedUsersCount, recentSignupsCount, nil
}

func (s *UserManagementService) GetUserStatistics(ctx context.Context, token string) (*models.GetUserStatisticsRow, string, string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, "", "", fmt.Errorf("only admins can get user statistics")
	}

	queries := models.New(s.db)
	stats, err := queries.GetUserStatistics(ctx)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get user statistics: %v", err)
	}

	newRegistrationsPercentageChange := calculatePercentageChange(stats.LastMonthNewUsersCount, stats.NewRegistrationsCount)
	approvedUsersPercentageChange := calculatePercentageChange(int64(stats.YesterdayApprovedCount.Int32), stats.ApprovedCount)

	return &stats, newRegistrationsPercentageChange, approvedUsersPercentageChange, nil
}

func calculatePercentageChange(oldValue, newValue int64) string {
	if oldValue == 0 && newValue == 0 {
		return "0.00%"
	}
	if oldValue == 0 {
		return "+∞%"
	}
	change := float64(newValue-oldValue) / float64(oldValue) * 100
	sign := "+"
	if change < 0 {
		sign = "-"
		change = math.Abs(change)
	}
	return fmt.Sprintf("%s%.2f%%", sign, change)
}

func (s *UserManagementService) ApproveUser(ctx context.Context, token string, userID int32) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return fmt.Errorf("only admins can approve users")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Update user status
	err = queries.UpdateUserStatus(ctx, models.UpdateUserStatusParams{
		Userid: userID,
		Status: sql.NullString{String: "approved", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update user status: %v", err)
	}

	// Get user details
	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	var address sql.NullString
	if user.Userrole == "school" {
		school, err := queries.GetSchoolByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get school details: %v", err)
		}
		address = sql.NullString{String: school.Address, Valid: true}
	}

	_, err = queries.CreateUserProfile(ctx, models.CreateUserProfileParams{
		Userid:             user.Userid,
		Name:               user.Name,
		Userrole:           user.Userrole,
		Email:              user.Email,
		Password:           user.Password,
		Verificationstatus: user.Verificationstatus,
		Address:            address,
		Phone:              sql.NullString{},
		Bio:                sql.NullString{},
		Profilepicture:     sql.NullString{},
		Gender:             user.Gender,
	})
	if err != nil {
		return fmt.Errorf("failed to create user profile: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send approval notification (consider moving this to a background job)
	go func() {
		if err := notifications.SendApprovalNotification(s.notificationService, user.Email, user.Name); err != nil {
			log.Printf("Failed to send approval notification: %v", err)
		}
	}()

	return nil
}

func (s *UserManagementService) RejectUser(ctx context.Context, token string, userID int32) error {
	// Validate token and check if the user is an admin
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return fmt.Errorf("only admins can reject users")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Get user details before rejection
	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	// Reject and delete user
	err = queries.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to reject user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send rejection notification (consider moving this to a background job)
	go func() {
		if err := notifications.SendRejectionNotification(s.notificationService, user.Email, user.Name); err != nil {
			log.Printf("Failed to send rejection notification: %v", err)
		}
	}()

	return nil
}

func (s *UserManagementService) ApproveUsers(ctx context.Context, token string, userIDs []int32) ([]int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("only admins can approve users")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	failedUserIDs := []int32{}

	for _, userID := range userIDs {
		err := queries.UpdateUserStatus(ctx, models.UpdateUserStatusParams{
			Userid: userID,
			Status: sql.NullString{String: "approved", Valid: true},
		})
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		user, err := queries.GetUserByID(ctx, userID)
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		var address sql.NullString
		if user.Userrole == "school" {
			school, err := queries.GetSchoolByUserID(ctx, userID)
			if err != nil {
				failedUserIDs = append(failedUserIDs, userID)
				continue
			}
			address = sql.NullString{String: school.Address, Valid: true}
		}

		_, err = queries.CreateUserProfile(ctx, models.CreateUserProfileParams{
			Userid:             user.Userid,
			Name:               user.Name,
			Userrole:           user.Userrole,
			Email:              user.Email,
			Password:           user.Password,
			Verificationstatus: user.Verificationstatus,
			Address:            address,
			Phone:              sql.NullString{},
			Bio:                sql.NullString{},
			Profilepicture:     sql.NullString{},
			Gender:             user.Gender,
		})
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		// Launch goroutine to send approval notification
		go func(userEmail, userName string, userId int32) {
			if err := notifications.SendApprovalNotification(s.notificationService, userEmail, userName); err != nil {
				fmt.Printf("Failed to send approval notification to user %d: %v\n", userId, err)
			}
		}(user.Email, user.Name, user.Userid)

	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return failedUserIDs, nil
}

func (s *UserManagementService) RejectUsers(ctx context.Context, token string, userIDs []int32) ([]int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("only admins can reject users")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	failedUserIDs := []int32{}

	for _, userID := range userIDs {
		user, err := queries.RejectAndGetUser(ctx, userID)
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		go func(userEmail, userName string, userId int32) {
			err := notifications.SendRejectionNotification(s.notificationService, userEmail, userName)
			if err != nil {
				fmt.Printf("Failed to send rejection notification to user %d: %v\n", userId, err)
			}
		}(user.Email, user.Name, userID)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return failedUserIDs, nil
}

func (s *UserManagementService) DeleteUsers(ctx context.Context, token string, userIDs []int32) ([]int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("only admins can delete users")
	}

	// Extract the admin's user ID from the token
	adminUserID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Check if the admin's user ID is in the list of user IDs to be deleted
	for _, userID := range userIDs {
		if int32(adminUserID) == userID {
			return nil, fmt.Errorf("admins cannot delete themselves")
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)
	failedUserIDs := []int32{}

	for _, userID := range userIDs {
		err := queries.DeleteUser(ctx, userID)
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return failedUserIDs, nil
}

func (s *UserManagementService) GetUserProfile(ctx context.Context, token string, userID int32) (*models.GetUserProfileRow, string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return nil, "", fmt.Errorf("you can only access your own profile unless you're an admin")
	}

	queries := models.New(s.db)
	profile, err := queries.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user profile: %v", err)
	}

	// Generate presigned URL for profile picture
	var profilePicturePresignedURL string
	if profile.Profilepicture.Valid && profile.Profilepicture.String != "" {
		key := utils.ExtractKeyFromURL(profile.Profilepicture.String)
		profilePicturePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return nil, "", fmt.Errorf("failed to generate presigned URL for profile picture: %v", err)
		}
	}

	return &profile, profilePicturePresignedURL, nil
}

func (s *UserManagementService) UpdateAdminProfile(ctx context.Context, token string, req *user_management.UpdateAdminProfileRequest) (string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != req.UserID || userRole != "admin" {
		return "", fmt.Errorf("unauthorized: only admins can update their own profile")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.UpdateAdminProfile(ctx, models.UpdateAdminProfileParams{
		Userid:         req.UserID,
		Name:           req.Name,
		Gender:         sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:          req.Email,
		Address:        sql.NullString{String: req.Address, Valid: req.Address != ""},
		Bio:            sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Profilepicture: sql.NullString{String: req.ProfilePictureUrl, Valid: req.ProfilePictureUrl != ""},
	})
	if err != nil {
		return "", fmt.Errorf("failed to update admin profile: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate presigned URL for profile picture
	var profilePicturePresignedURL string
	if req.ProfilePictureUrl != "" {
		key := utils.ExtractKeyFromURL(req.ProfilePictureUrl)
		profilePicturePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL for profile picture: %v", err)
		}
	}

	return profilePicturePresignedURL, nil
}

func (s *UserManagementService) UpdateSchoolProfile(ctx context.Context, token string, req *user_management.UpdateSchoolProfileRequest) (string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != req.UserID || userRole != "school" {
		return "", fmt.Errorf("unauthorized: only school contact persons can update their school profile")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Update Users table
	err = queries.UpdateSchoolUser(ctx, models.UpdateSchoolUserParams{
		Userid: req.UserID,
		Name:   req.ContactPersonName,
		Gender: sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:  req.ContactEmail,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	// Update UserProfiles table
	err = queries.UpdateSchoolUserProfile(ctx, models.UpdateSchoolUserProfileParams{
		Userid:         req.UserID,
		Name:           req.ContactPersonName,
		Email:          req.ContactEmail,
		Gender:         sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Address:        sql.NullString{String: req.Address, Valid: req.Address != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Bio:            sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		Profilepicture: sql.NullString{String: req.ProfilePictureUrl, Valid: req.ProfilePictureUrl != ""},
	})
	if err != nil {
		return "", fmt.Errorf("failed to update user profile: %v", err)
	}

	// Update Schools table
	err = queries.UpdateSchoolDetails(ctx, models.UpdateSchoolDetailsParams{
		Contactpersonid:         req.UserID,
		Contactpersonnationalid: sql.NullString{String: req.ContactPersonNationalId, Valid: req.ContactPersonNationalId != ""},
		Schoolname:              req.SchoolName,
		Address:                 req.Address,
		Schoolemail:             req.SchoolEmail,
		Schooltype:              req.SchoolType,
		Contactemail:            req.ContactEmail,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update school details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate presigned URL for profile picture
	var profilePicturePresignedURL string
	if req.ProfilePictureUrl != "" {
		key := utils.ExtractKeyFromURL(req.ProfilePictureUrl)
		profilePicturePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL for profile picture: %v", err)
		}
	}

	return profilePicturePresignedURL, nil
}

func (s *UserManagementService) UpdateStudentProfile(ctx context.Context, token string, req *user_management.UpdateStudentProfileRequest) (string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != req.UserID || userRole != "student" {
		return "", fmt.Errorf("unauthorized: only students can update their own profile")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Update Users table
	err = queries.UpdateStudentUser(ctx, models.UpdateStudentUserParams{
		Userid:  req.UserID,
		Column2: sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		Column3: sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		Gender:  sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:   req.Email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update user: %v", err)
	}

	// Update UserProfiles table
	err = queries.UpdateStudentUserProfile(ctx, models.UpdateStudentUserProfileParams{
		Userid:         req.UserID,
		Column2:        sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		Column3:        sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		Gender:         sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:          req.Email,
		Address:        sql.NullString{String: req.Address, Valid: req.Address != ""},
		Bio:            sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		Profilepicture: sql.NullString{String: req.ProfilePictureUrl, Valid: req.ProfilePictureUrl != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
	})
	if err != nil {
		return "", fmt.Errorf("failed to update user profile: %v", err)
	}

	// Update Students table
	dateOfBirth, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return "", fmt.Errorf("invalid date format for date of birth: %v", err)
	}

	err = queries.UpdateStudentDetails(ctx, models.UpdateStudentDetailsParams{
		Userid:    req.UserID,
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Gender:    sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Grade:     req.Grade,
		Column7:   dateOfBirth,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update student details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate presigned URL for profile picture
	var profilePicturePresignedURL string
	if req.ProfilePictureUrl != "" {
		key := utils.ExtractKeyFromURL(req.ProfilePictureUrl)
		profilePicturePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return "", fmt.Errorf("failed to generate presigned URL for profile picture: %v", err)
		}
	}

	return profilePicturePresignedURL, nil
}

func (s *UserManagementService) UpdateVolunteerProfile(ctx context.Context, token string, req *user_management.UpdateVolunteerProfileRequest) (string, string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != req.UserID || userRole != "volunteer" {
		return "", "", fmt.Errorf("unauthorized: only volunteers can update their own profile")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Update Users table
	err = queries.UpdateVolunteerUser(ctx, models.UpdateVolunteerUserParams{
		Userid:  req.UserID,
		Column2: sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		Column3: sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		Gender:  sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:   req.Email,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to update user: %v", err)
	}

	// Update UserProfiles table
	err = queries.UpdateVolunteerUserProfile(ctx, models.UpdateVolunteerUserProfileParams{
		Userid:         req.UserID,
		Column2:        sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		Column3:        sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		Gender:         sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Email:          req.Email,
		Address:        sql.NullString{String: req.Address, Valid: req.Address != ""},
		Bio:            sql.NullString{String: req.Bio, Valid: req.Bio != ""},
		Profilepicture: sql.NullString{String: req.ProfilePictureUrl, Valid: req.ProfilePictureUrl != ""},
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to update user profile: %v", err)
	}

	// Update Volunteers table
	err = queries.UpdateVolunteerDetails(ctx, models.UpdateVolunteerDetailsParams{
		Userid:                 req.UserID,
		Firstname:              req.FirstName,
		Lastname:               req.LastName,
		Gender:                 sql.NullString{String: req.Gender, Valid: req.Gender != ""},
		Nationalid:             sql.NullString{String: req.NationalID, Valid: req.NationalID != ""},
		Graduateyear:           sql.NullInt32{Int32: req.GraduateYear, Valid: req.GraduateYear != 0},
		Isenrolledinuniversity: sql.NullBool{Bool: req.IsEnrolledInUniversity, Valid: true},
		Hasinternship:          sql.NullBool{Bool: req.HasInternship, Valid: true},
		Role:                   req.Role,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to update volunteer details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate presigned URLs
	var profilePicturePresignedURL, safeguardCertificatePresignedURL string
	if req.ProfilePictureUrl != "" {
		key := utils.ExtractKeyFromURL(req.ProfilePictureUrl)
		profilePicturePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate presigned URL for profile picture: %v", err)
		}
	}
	if req.SafeguardCertificateUrl != "" {
		key := utils.ExtractKeyFromURL(req.SafeguardCertificateUrl)
		safeguardCertificatePresignedURL, err = s.s3Client.GetSignedURL(ctx, key, time.Hour)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate presigned URL for safeguard certificate: %v", err)
		}
	}

	return profilePicturePresignedURL, safeguardCertificatePresignedURL, nil
}

func (s *UserManagementService) DeleteUserProfile(ctx context.Context, token string, userID int32) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return fmt.Errorf("you can only delete your own profile unless you're an admin")
	}

	queries := models.New(s.db)

	err = queries.SoftDeleteUserProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to soft delete user profile: %v", err)
	}

	return nil
}

func (s *UserManagementService) DeactivateAccount(ctx context.Context, token string, userID int32) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return fmt.Errorf("you can only deactivate your own account unless you're an admin")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.DeactivateAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate account: %v", err)
	}

	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send account deactivation notification in a goroutine
	go func() {
		err = notifications.SendAccountDeactivationNotification(s.notificationService, user.Email, user.Name)
		if err != nil {
			fmt.Printf("Failed to send account deactivation notification: %v\n", err)
		}
	}()

	return nil
}

func (s *UserManagementService) ReactivateAccount(ctx context.Context, token string, userID int32) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return fmt.Errorf("you can only reactivate your own account unless you're an admin")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.ReactivateAccount(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to reactivate account: %v", err)
	}

	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	go func() {
		err = notifications.SendAccountReactivationNotification(s.notificationService, user.Email, user.Name)
		if err != nil {
			fmt.Printf("Failed to send account reactivation notification: %v\n", err)
		}
	}()

	return nil
}

func (s *UserManagementService) GetAccountStatus(ctx context.Context, token string, userID int32) (string, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return "", fmt.Errorf("you can only get your own account status unless you're an admin")
	}

	queries := models.New(s.db)

	status, err := queries.GetAccountStatus(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get account status: %v", err)
	}

	return status, nil
}

func (s *UserManagementService) GetVolunteersAndAdmins(ctx context.Context, token string, page, pageSize int32) ([]models.User, int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, 0, fmt.Errorf("only admins can access this information")
	}

	queries := models.New(s.db)
	users, err := queries.GetVolunteersAndAdmins(ctx, models.GetVolunteersAndAdminsParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch volunteers and admins: %v", err)
	}

	totalCount, err := queries.GetTotalVolunteersAndAdminsCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %v", err)
	}

	return users, int32(totalCount), nil
}

func (s *UserManagementService) generateSecret() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}

func (s *UserManagementService) InitiatePasswordUpdate(ctx context.Context, token string, userID int32) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	if tokenUserID != userID {
		return fmt.Errorf("unauthorized: token does not match user ID")
	}

	// Generate secret for TOTP
	secret, err := s.generateSecret()
	if err != nil {
		return fmt.Errorf("failed to generate secret: %v", err)
	}

	// Set expiration time (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Store the secret and expiration time, and get user details
	user, err := queries.SetPasswordResetCodeAndGetUser(ctx, models.SetPasswordResetCodeAndGetUserParams{
		Userid:            userID,
		ResetToken:        sql.NullString{String: secret, Valid: true},
		ResetTokenExpires: sql.NullTime{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to set password reset code: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Generate verification code
	verificationCode, err := totp.GenerateCodeCustom(secret, time.Now(), totp.ValidateOpts{
		Period: 900, // 15 minutes in seconds
		Skew:   1,   // Allow 1 period before and after
		Digits: 6,
	})
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %v", err)
	}

	// Send verification code via email in a goroutine
	go func() {
		err := notifications.SendPasswordUpdateVerificationEmail(s.notificationService, user.Email, user.Name, verificationCode)
		if err != nil {
			fmt.Printf("Failed to send verification email: %v\n", err)
		}
	}()

	return nil
}

func (s *UserManagementService) VerifyAndUpdatePassword(ctx context.Context, token string, userID int32, verificationCode string, newPassword string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	if tokenUserID != userID {
		return fmt.Errorf("unauthorized: token does not match user ID")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Get user details and reset code
	user, err := queries.ValidateResetCodeAndGetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no reset code found or reset code expired")
		}
		return fmt.Errorf("failed to get user details: %v", err)
	}

	// Validate the verification code
	valid, err := totp.ValidateCustom(verificationCode, user.ResetToken.String, time.Now(), totp.ValidateOpts{
		Period: 900, // 15 minutes in seconds
		Skew:   1,   // Allow 1 period before and after
		Digits: 6,
	})
	if err != nil {
		return fmt.Errorf("failed to validate verification code: %v", err)
	}
	if !valid {
		return fmt.Errorf("invalid verification code")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update the password in both tables and clear the reset code
	err = queries.UpdatePasswordAndClearResetCode(ctx, models.UpdatePasswordAndClearResetCodeParams{
		Userid:   userID,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Send password update confirmation email in a goroutine
	go func() {
		err := notifications.SendPasswordUpdateConfirmationEmail(s.notificationService, user.Email, user.Name)
		if err != nil {
			fmt.Printf("Failed to send password update confirmation email: %v\n", err)
		}
	}()

	return nil
}
