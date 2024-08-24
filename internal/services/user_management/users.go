package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
	email "github.com/iRankHub/backend/internal/utils/emails"
)

type UserManagementService struct {
	db *sql.DB
}

func NewUserManagementService(db *sql.DB) *UserManagementService {
	return &UserManagementService{
		db: db,
	}
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

func (s *UserManagementService) GetAllUsers(ctx context.Context, token string, page, pageSize int32) ([]models.User, int32, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, 0, fmt.Errorf("only admins can get all users")
	}

	queries := models.New(s.db)
	users, err := queries.GetAllUsers(ctx, models.GetAllUsersParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %v", err)
	}

	totalCount, err := queries.GetTotalUserCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total user count: %v", err)
	}

	// Convert int64 to int32, checking for overflow
	if totalCount > math.MaxInt32 {
		return nil, 0, fmt.Errorf("total user count exceeds maximum value for int32")
	}

	return users, int32(totalCount), nil
}

func (s *UserManagementService) GetUserDetails(ctx context.Context, token string, userID int32) (models.User, models.Userprofile, error) {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return models.User{}, models.Userprofile{}, fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	userRole := claims["user_role"].(string)

	if tokenUserID != userID && userRole != "admin" {
		return models.User{}, models.Userprofile{}, fmt.Errorf("you can only access your own details unless you're an admin")
	}

	queries := models.New(s.db)
	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, models.Userprofile{}, fmt.Errorf("failed to get user details: %v", err)
	}

	profile, err := queries.GetUserProfile(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, models.Userprofile{}, nil
		}
		return models.User{}, models.Userprofile{}, fmt.Errorf("failed to get user profile: %v", err)
	}

	return user, profile, nil
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
		Verificationstatus: user.Verificationstatus,
		Address:            address,
		Phone:              sql.NullString{},
		Bio:                sql.NullString{},
		Profilepicture:     nil,
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
		if err := email.SendApprovalNotification(user.Email, user.Name); err != nil {
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
		if err := email.SendRejectionNotification(user.Email, user.Name); err != nil {
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
			Profilepicture:     nil,
			Gender:             user.Gender,
		})
		if err != nil {
			failedUserIDs = append(failedUserIDs, userID)
			continue
		}

		// Launch goroutine to send approval notification
		go func(userEmail, userName string, userId int32) {
			if err := email.SendApprovalNotification(userEmail, userName); err != nil {
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

		err = email.SendRejectionNotification(user.Email, user.Name)
		if err != nil {
			fmt.Printf("Failed to send rejection notification to user %d: %v\n", userID, err)
		}
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

func (s *UserManagementService) UpdateUserProfile(ctx context.Context, token string, userID int32, name, email, address, phone, bio string, profilePicture []byte, gender string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %v", err)
	}

	tokenUserID := int32(claims["user_id"].(float64))
	if tokenUserID != userID {
		return fmt.Errorf("you can only update your own profile")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	// Get the user's role
	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	// Update UserProfiles table
	_, err = queries.UpdateUserProfile(ctx, models.UpdateUserProfileParams{
		Userid:         userID,
		Name:           name,
		Email:          email,
		Address:        sql.NullString{String: address, Valid: address != ""},
		Phone:          sql.NullString{String: phone, Valid: phone != ""},
		Bio:            sql.NullString{String: bio, Valid: bio != ""},
		Profilepicture: profilePicture,
		Gender:         sql.NullString{String: gender, Valid: gender != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	// Update Users table
	_, err = queries.UpdateUser(ctx, models.UpdateUserParams{
		Userid: userID,
		Name:   name,
		Email:  email,
		Gender: sql.NullString{String: gender, Valid: gender != ""},
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// If the user is a school, update the school's address
	if user.Userrole == "school" && address != "" {
		_, err = queries.UpdateSchoolAddress(ctx, models.UpdateSchoolAddressParams{
			Contactpersonid: userID,
			Address:         address,
		})
		if err != nil {
			return fmt.Errorf("failed to update school address: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	queries := models.New(tx)

	err = queries.DeleteUserProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user profile: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
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

	err = email.SendAccountDeactivationNotification(user.Email, user.Name)
	if err != nil {
		fmt.Printf("Failed to send account deactivation notification: %v\n", err)
	}

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

	err = email.SendAccountReactivationNotification(user.Email, user.Name)
	if err != nil {
		fmt.Printf("Failed to send account reactivation notification: %v\n", err)
	}

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
