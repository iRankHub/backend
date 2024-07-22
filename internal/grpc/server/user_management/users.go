package server

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/user_management"
	services "github.com/iRankHub/backend/internal/services/user_management"
)

type userManagementServer struct {
	user_management.UnimplementedUserManagementServiceServer
	db *sql.DB
	userManagementService *services.UserManagementService
}

func NewUserManagementServer(db *sql.DB) (user_management.UserManagementServiceServer, error) {
	userManagementService := services.NewUserManagementService(db)

	return &userManagementServer{
		db: db,
		userManagementService: userManagementService,
	}, nil
}

func (s *userManagementServer) GetPendingUsers(ctx context.Context, req *user_management.GetPendingUsersRequest) (*user_management.GetPendingUsersResponse, error) {
	users, err := s.userManagementService.GetPendingUsers(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get pending users: %v", err)
	}

	var userSummaries []*user_management.UserSummary
	for _, user := range users {
		signUpDate := ""
		if user.CreatedAt.Valid {
			signUpDate = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		userSummaries = append(userSummaries, &user_management.UserSummary{
			UserID:     user.Userid,
			Name:       user.Name,
			Email:      user.Email,
			UserRole:   user.Userrole,
			SignUpDate: signUpDate,
		})
	}

	return &user_management.GetPendingUsersResponse{
		Users: userSummaries,
	}, nil
}

func (s *userManagementServer) GetUserDetails(ctx context.Context, req *user_management.GetUserDetailsRequest) (*user_management.GetUserDetailsResponse, error) {
	user, profile, err := s.userManagementService.GetUserDetails(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get user details: %v", err)
	}

	signUpDate := ""
	if user.CreatedAt.Valid {
		signUpDate = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	userDetails := &user_management.UserDetails{
		UserID:     user.Userid,
		Name:       user.Name,
		Email:      user.Email,
		UserRole:   user.Userrole,
		SignUpDate: signUpDate,
		Profile: &user_management.UserProfile{
			Address:        profile.Address.String,
			Phone:          profile.Phone.String,
			Bio:            profile.Bio.String,
			ProfilePicture: profile.Profilepicture,
		},
	}

	return &user_management.GetUserDetailsResponse{
		User: userDetails,
	}, nil
}

func (s *userManagementServer) ApproveUser(ctx context.Context, req *user_management.ApproveUserRequest) (*user_management.ApproveUserResponse, error) {
	err := s.userManagementService.ApproveUser(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to approve user: %v", err)
	}

	return &user_management.ApproveUserResponse{
		Success: true,
		Message: "User approved successfully",
	}, nil
}

func (s *userManagementServer) RejectUser(ctx context.Context, req *user_management.RejectUserRequest) (*user_management.RejectUserResponse, error) {
	err := s.userManagementService.RejectUser(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to reject user: %v", err)
	}

	return &user_management.RejectUserResponse{
		Success: true,
		Message: "User rejected successfully",
	}, nil
}

func (s *userManagementServer) UpdateUserProfile(ctx context.Context, req *user_management.UpdateUserProfileRequest) (*user_management.UpdateUserProfileResponse, error) {
	err := s.userManagementService.UpdateUserProfile(ctx, req.Token, req.UserID, req.Name, req.Email, req.Address, req.Phone, req.Bio, req.ProfilePicture)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update user profile: %v", err)
	}

	return &user_management.UpdateUserProfileResponse{
		Success: true,
		Message: "User profile updated successfully",
	}, nil
}

func (s *userManagementServer) DeleteUserProfile(ctx context.Context, req *user_management.DeleteUserProfileRequest) (*user_management.DeleteUserProfileResponse, error) {
	err := s.userManagementService.DeleteUserProfile(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete user profile: %v", err)
	}

	return &user_management.DeleteUserProfileResponse{
		Success: true,
		Message: "User profile deleted successfully",
	}, nil
}

func (s *userManagementServer) DeactivateAccount(ctx context.Context, req *user_management.DeactivateAccountRequest) (*user_management.DeactivateAccountResponse, error) {
	err := s.userManagementService.DeactivateAccount(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to deactivate account: %v", err)
	}

	return &user_management.DeactivateAccountResponse{
		Success: true,
		Message: "Account deactivated successfully",
	}, nil
}

func (s *userManagementServer) ReactivateAccount(ctx context.Context, req *user_management.ReactivateAccountRequest) (*user_management.ReactivateAccountResponse, error) {
	err := s.userManagementService.ReactivateAccount(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to reactivate account: %v", err)
	}

	return &user_management.ReactivateAccountResponse{
		Success: true,
		Message: "Account reactivated successfully",
	}, nil
}

func (s *userManagementServer) GetAccountStatus(ctx context.Context, req *user_management.GetAccountStatusRequest) (*user_management.GetAccountStatusResponse, error) {
	accountStatus, err := s.userManagementService.GetAccountStatus(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get account status: %v", err)
	}

	return &user_management.GetAccountStatusResponse{
		Status: accountStatus,
	}, nil
}