package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
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

	err = queries.UpdateUserStatus(ctx, models.UpdateUserStatusParams{
		Userid: userID,
		Status: sql.NullString{String: "approved", Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update user status: %v", err)
	}

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
		Userid:         user.Userid,
		Address:        address,
		Phone:          sql.NullString{},
		Bio:            sql.NullString{},
		Profilepicture: nil,
	})
	if err != nil {
		return fmt.Errorf("failed to create user profile: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	err = utils.SendApprovalNotification(user.Email, user.Name)
	if err != nil {
		fmt.Printf("Failed to send approval notification: %v\n", err)
	}

	return nil
}

func (s *UserManagementService) RejectUser(ctx context.Context, token string, userID int32) error {
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

	user, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user details: %v", err)
	}

	err = queries.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	err = utils.SendRejectionNotification(user.Email, user.Name)
	if err != nil {
		fmt.Printf("Failed to send rejection notification: %v\n", err)
	}

	return nil
}

func (s *UserManagementService) UpdateUserProfile(ctx context.Context, token string, userID int32, name, email, address, phone, bio string, profilePicture []byte) error {
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

	_, err = queries.UpdateUserProfile(ctx, models.UpdateUserProfileParams{
		Userid:         userID,
		Name:           name,
		Email:          email,
		Address:        sql.NullString{String: address, Valid: address != ""},
		Phone:          sql.NullString{String: phone, Valid: phone != ""},
		Bio:            sql.NullString{String: bio, Valid: bio != ""},
		Profilepicture: profilePicture,
	})
	if err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
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

	err = utils.SendAccountDeactivationNotification(user.Email, user.Name)
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

	err = utils.SendAccountReactivationNotification(user.Email, user.Name)
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