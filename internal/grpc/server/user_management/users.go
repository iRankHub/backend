package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iRankHub/backend/internal/grpc/proto/user_management"
	"github.com/iRankHub/backend/internal/models"
	services "github.com/iRankHub/backend/internal/services/user_management"

)

type userManagementServer struct {
	user_management.UnimplementedUserManagementServiceServer
	db                          *sql.DB
	userManagementService       *services.UserManagementService
	countryManagementService    *services.CountryService
	schoolsManagementService    *services.SchoolService
	studentsManagementService   *services.StudentService
	volunteersManagementService *services.VolunteerService
}

func NewUserManagementServer(db *sql.DB) (user_management.UserManagementServiceServer, error) {
	userManagementService, err := services.NewUserManagementService(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create user management service: %v", err)
	}

	return &userManagementServer{
		db:                          db,
		userManagementService:       userManagementService,
		countryManagementService:    services.NewCountryManagementService(db),
		schoolsManagementService:    services.NewSchoolsManagementService(db),
		studentsManagementService:   services.NewStudentsManagementService(db),
		volunteersManagementService: services.NewVolunteersManagementService(db),
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
			Gender:     user.Gender.String,
			Status:     user.Status.String,
		})
	}

	return &user_management.GetPendingUsersResponse{
		Users: userSummaries,
	}, nil
}

func (s *userManagementServer) ApproveUser(ctx context.Context, req *user_management.ApproveUserRequest) (*user_management.ApproveUserResponse, error) {
	err := s.userManagementService.ApproveUser(ctx, req.Token, req.UserID)
	if err != nil {
		log.Printf("Failed to approve user: %v", err)
		return &user_management.ApproveUserResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to approve user: %v", err),
		}, status.Errorf(codes.Internal, "Failed to approve user: %v", err)
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

func (s *userManagementServer) ApproveUsers(ctx context.Context, req *user_management.ApproveUsersRequest) (*user_management.ApproveUsersResponse, error) {
	failedUserIDs, err := s.userManagementService.ApproveUsers(ctx, req.Token, req.UserIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to approve users: %v", err)
	}

	message := "All users approved successfully"
	if len(failedUserIDs) > 0 {
		message = "Some users could not be approved"
	}

	return &user_management.ApproveUsersResponse{
		Success:       len(failedUserIDs) < len(req.UserIDs),
		Message:       message,
		FailedUserIDs: failedUserIDs,
	}, nil
}

func (s *userManagementServer) RejectUsers(ctx context.Context, req *user_management.RejectUsersRequest) (*user_management.RejectUsersResponse, error) {
	failedUserIDs, err := s.userManagementService.RejectUsers(ctx, req.Token, req.UserIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to reject users: %v", err)
	}

	message := "All users rejected successfully"
	if len(failedUserIDs) > 0 {
		message = "Some users could not be rejected"
	}

	return &user_management.RejectUsersResponse{
		Success:       len(failedUserIDs) < len(req.UserIDs),
		Message:       message,
		FailedUserIDs: failedUserIDs,
	}, nil
}

func (s *userManagementServer) DeleteUsers(ctx context.Context, req *user_management.DeleteUsersRequest) (*user_management.DeleteUsersResponse, error) {
	failedUserIDs, err := s.userManagementService.DeleteUsers(ctx, req.Token, req.UserIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete users: %v", err)
	}

	message := "All users deleted successfully"
	if len(failedUserIDs) > 0 {
		message = "Some users could not be deleted"
	}

	return &user_management.DeleteUsersResponse{
		Success:       len(failedUserIDs) < len(req.UserIDs),
		Message:       message,
		FailedUserIDs: failedUserIDs,
	}, nil
}


func (s *userManagementServer) GetAllUsers(ctx context.Context, req *user_management.GetAllUsersRequest) (*user_management.GetAllUsersResponse, error) {
    users, totalCount, approvedUsersCount, recentSignupsCount, err := s.userManagementService.GetAllUsers(ctx, req.Token, req.Page, req.PageSize)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to get all users: %v", err)
    }

    var userSummaries []*user_management.UserSummary
    for _, user := range users {
        signUpDate := ""
        if user.CreatedAt.Valid {
            signUpDate = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
        }

        var idebateID string
        if user.Idebateid != nil {
            switch v := user.Idebateid.(type) {
            case string:
                idebateID = v
            case []byte:
                idebateID = string(v)
            default:
                idebateID = fmt.Sprintf("%v", v)
            }
        }

        userSummaries = append(userSummaries, &user_management.UserSummary{
            UserID:     user.Userid,
            Name:       user.Displayname.(string),
            Email:      user.Email,
            UserRole:   user.Userrole,
            SignUpDate: signUpDate,
            Gender:     user.Gender.String,
            Status:     user.Status.String,
            IdebateID:  idebateID,
        })
    }

    return &user_management.GetAllUsersResponse{
        Users:              userSummaries,
        TotalCount:         totalCount,
        ApprovedUsersCount: approvedUsersCount,
        RecentSignupsCount: recentSignupsCount,
    }, nil
}

func (s *userManagementServer) GetUserStatistics(ctx context.Context, req *user_management.GetUserStatisticsRequest) (*user_management.GetUserStatisticsResponse, error) {
    stats, newRegistrationsPercentageChange, approvedUsersPercentageChange, err := s.userManagementService.GetUserStatistics(ctx, req.Token)
    if err != nil {
        log.Printf("Error in GetUserStatistics: %v", err)
        return nil, status.Errorf(codes.Internal, "Failed to get user statistics: %v", err)
    }

    return &user_management.GetUserStatisticsResponse{
        AdminCount:                       stats.AdminCount,
        SchoolCount:                      stats.SchoolCount,
        StudentCount:                     stats.StudentCount,
        VolunteerCount:                   stats.VolunteerCount,
        ApprovedCount:                    stats.ApprovedCount,
        NewRegistrationsCount:            stats.NewRegistrationsCount,
        NewRegistrationsPercentageChange: newRegistrationsPercentageChange,
        ApprovedUsersPercentageChange:    approvedUsersPercentageChange,
    }, nil
}

func (s *userManagementServer) GetUserProfile(ctx context.Context, req *user_management.GetUserProfileRequest) (*user_management.GetUserProfileResponse, error) {
	profile, err := s.userManagementService.GetUserProfile(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get user profile: %v", err)
	}

	return &user_management.GetUserProfileResponse{
		Profile: convertModelProfileToProto(profile),
	}, nil
}

func (s *userManagementServer) UpdateAdminProfile(ctx context.Context, req *user_management.UpdateAdminProfileRequest) (*user_management.UpdateAdminProfileResponse, error) {
    err := s.userManagementService.UpdateAdminProfile(ctx, req.Token, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update admin profile: %v", err)
    }

    return &user_management.UpdateAdminProfileResponse{
        Success: true,
        Message: "Admin profile updated successfully",
    }, nil
}

func (s *userManagementServer) UpdateSchoolProfile(ctx context.Context, req *user_management.UpdateSchoolProfileRequest) (*user_management.UpdateSchoolProfileResponse, error) {
    err := s.userManagementService.UpdateSchoolProfile(ctx, req.Token, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update school profile: %v", err)
    }

    return &user_management.UpdateSchoolProfileResponse{
        Success: true,
        Message: "School profile updated successfully",
    }, nil
}

func (s *userManagementServer) UpdateStudentProfile(ctx context.Context, req *user_management.UpdateStudentProfileRequest) (*user_management.UpdateStudentProfileResponse, error) {
    err := s.userManagementService.UpdateStudentProfile(ctx, req.Token, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update student profile: %v", err)
    }

    return &user_management.UpdateStudentProfileResponse{
        Success: true,
        Message: "Student profile updated successfully",
    }, nil
}

func (s *userManagementServer) UpdateVolunteerProfile(ctx context.Context, req *user_management.UpdateVolunteerProfileRequest) (*user_management.UpdateVolunteerProfileResponse, error) {
    err := s.userManagementService.UpdateVolunteerProfile(ctx, req.Token, req)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to update volunteer profile: %v", err)
    }

    return &user_management.UpdateVolunteerProfileResponse{
        Success: true,
        Message: "Volunteer profile updated successfully",
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

func (s *userManagementServer) GetStudents(ctx context.Context, req *user_management.GetStudentsRequest) (*user_management.GetStudentsResponse, error) {
	students, totalCount, err := s.studentsManagementService.GetStudents(ctx, req.Token, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get students: %v", err)
	}

	var protoStudents []*user_management.Student
	for _, student := range students {
		dateOfBirth := ""
		if student.Dateofbirth.Valid {
			dateOfBirth = student.Dateofbirth.Time.Format("2006-01-02")
		}
		email := ""
		if student.Email.Valid {
			email = student.Email.String
		}
		protoStudents = append(protoStudents, &user_management.Student{
			StudentID:   student.Studentid,
			FirstName:   student.Firstname,
			LastName:    student.Lastname,
			Grade:       student.Grade,
			DateOfBirth: dateOfBirth,
			Email:       email,
			SchoolID:    student.Schoolid,
		})
	}

	return &user_management.GetStudentsResponse{
		Students:   protoStudents,
		TotalCount: totalCount,
	}, nil
}

func (s *userManagementServer) GetVolunteers(ctx context.Context, req *user_management.GetVolunteersRequest) (*user_management.GetVolunteersResponse, error) {
	volunteers, totalCount, err := s.volunteersManagementService.GetVolunteers(ctx, req.Token, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volunteers: %v", err)
	}

	var protoVolunteers []*user_management.Volunteer
	for _, volunteer := range volunteers {
		dateOfBirth := ""
		if volunteer.Dateofbirth.Valid {
			dateOfBirth = volunteer.Dateofbirth.Time.Format("2006-01-02")
		}
		graduateYear := int32(0)
		if volunteer.Graduateyear.Valid {
			graduateYear = volunteer.Graduateyear.Int32
		}
		protoVolunteers = append(protoVolunteers, &user_management.Volunteer{
			VolunteerID:          volunteer.Volunteerid,
			FirstName:            volunteer.Firstname,
			LastName:             volunteer.Lastname,
			DateOfBirth:          dateOfBirth,
			Role:                 volunteer.Role,
			GraduateYear:         graduateYear,
			SafeGuardCertificate: volunteer.Safeguardcertificate,
		})
	}

	return &user_management.GetVolunteersResponse{
		Volunteers: protoVolunteers,
		TotalCount: totalCount,
	}, nil
}

func (s *userManagementServer) GetSchools(ctx context.Context, req *user_management.GetSchoolsRequest) (*user_management.GetSchoolsResponse, error) {
	schools, totalCount, err := s.schoolsManagementService.GetSchools(ctx, req.Token, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get schools: %v", err)
	}

	var protoSchools []*user_management.School
	for _, school := range schools {
		protoSchools = append(protoSchools, &user_management.School{
			SchoolID:     school.Schoolid,
			Name:         school.Schoolname,
			Address:      school.Address,
			Country:      school.Country.String,
			Province:     school.Province.String,
			District:     school.District.String,
			SchoolType:   school.Schooltype,
			ContactEmail: school.Contactemail,
			SchoolEmail:  school.Schoolemail,
		})
	}

	return &user_management.GetSchoolsResponse{
		Schools:    protoSchools,
		TotalCount: totalCount,
	}, nil
}

func (s *userManagementServer) GetCountries(ctx context.Context, req *user_management.GetCountriesRequest) (*user_management.GetCountriesResponse, error) {
	countries, err := s.countryManagementService.GetCountries(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get countries: %v", err)
	}

	var protoCountries []*user_management.Country
	for _, country := range countries {
		protoCountries = append(protoCountries, &user_management.Country{
			Name: country.Countryname,
			Code: country.Isocode,
		})
	}

	return &user_management.GetCountriesResponse{
		Countries: protoCountries,
	}, nil
}

func (s *userManagementServer) GetCountriesNoAuth(ctx context.Context, req *user_management.GetCountriesNoAuthRequest) (*user_management.GetCountriesNoAuthResponse, error) {
	countries, err := s.countryManagementService.GetCountriesNoAuth(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get countries: %v", err)
	}

	var protoCountries []*user_management.Country
	for _, country := range countries {
		protoCountries = append(protoCountries, &user_management.Country{
			Name: country.Countryname,
			Code: country.Isocode,
		})
	}

	return &user_management.GetCountriesNoAuthResponse{
		Countries: protoCountries,
	}, nil
}

func (s *userManagementServer) GetVolunteersAndAdmins(ctx context.Context, req *user_management.GetVolunteersAndAdminsRequest) (*user_management.GetVolunteersAndAdminsResponse, error) {
	users, totalCount, err := s.userManagementService.GetVolunteersAndAdmins(ctx, req.Token, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get volunteers and admins: %v", err)
	}

	var userSummaries []*user_management.UserSummary
	for _, user := range users {
		userSummaries = append(userSummaries, &user_management.UserSummary{
			UserID:     user.Userid,
			Name:       user.Name,
			Email:      user.Email,
			UserRole:   user.Userrole,
			SignUpDate: user.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			Gender:     user.Gender.String,
			Status:     user.Status.String,
		})
	}

	return &user_management.GetVolunteersAndAdminsResponse{
		Users:      userSummaries,
		TotalCount: totalCount,
	}, nil
}

func (s *userManagementServer) GetSchoolsNoAuth(ctx context.Context, req *user_management.GetSchoolsNoAuthRequest) (*user_management.GetSchoolsNoAuthResponse, error) {
	schools, totalCount, err := s.schoolsManagementService.GetSchoolsNoAuth(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get schools: %v", err)
	}

	var protoSchools []*user_management.School
	for _, school := range schools {
		protoSchools = append(protoSchools, &user_management.School{
			SchoolID:     school.Schoolid,
			Name:         school.Schoolname,
			Address:      school.Address,
			Country:      school.Country.String,
			Province:     school.Province.String,
			District:     school.District.String,
			SchoolType:   school.Schooltype,
			ContactEmail: school.Contactemail,
			SchoolEmail:  school.Schoolemail,
		})
	}

	return &user_management.GetSchoolsNoAuthResponse{
		Schools:    protoSchools,
		TotalCount: totalCount,
	}, nil
}

func (s *userManagementServer) InitiatePasswordUpdate(ctx context.Context, req *user_management.InitiatePasswordUpdateRequest) (*user_management.InitiatePasswordUpdateResponse, error) {
	err := s.userManagementService.InitiatePasswordUpdate(ctx, req.Token, req.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to initiate password update: %v", err)
	}

	return &user_management.InitiatePasswordUpdateResponse{
		Success: true,
		Message: "Password update initiated. Please check your email for the verification code.",
	}, nil
}

func (s *userManagementServer) VerifyAndUpdatePassword(ctx context.Context, req *user_management.VerifyAndUpdatePasswordRequest) (*user_management.VerifyAndUpdatePasswordResponse, error) {
	err := s.userManagementService.VerifyAndUpdatePassword(ctx, req.Token, req.UserID, req.VerificationCode, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to verify and update password: %v", err)
	}

	return &user_management.VerifyAndUpdatePasswordResponse{
		Success: true,
		Message: "Password updated successfully.",
	}, nil
}

func (s *userManagementServer) GetSchoolIDsByNames(ctx context.Context, req *user_management.GetSchoolIDsByNamesRequest) (*user_management.GetSchoolIDsByNamesResponse, error) {
	schoolIDs, err := s.schoolsManagementService.GetSchoolIDsByNames(ctx, req.Token, req.SchoolNames)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get school IDs: %v", err)
	}

	return &user_management.GetSchoolIDsByNamesResponse{
		SchoolIds: schoolIDs,
	}, nil
}

func convertModelProfileToProto(profile *models.GetUserProfileRow) *user_management.UserProfile {
	protoProfile := &user_management.UserProfile{
		UserID:               profile.Userid,
		Name:                 profile.Name,
		Email:                profile.Email,
		UserRole:             profile.Userrole,
		Gender:               profile.Gender.String,
		Address:              profile.Address.String,
		Phone:                profile.Phone.String,
		Bio:                  profile.Bio.String,
		ProfilePicture:       profile.Profilepicture,
		VerificationStatus:   profile.Verificationstatus.Bool,
		SignUpDate:           profile.Signupdate.Time.Format("2006-01-02 15:04:05"),
		TwoFactorEnabled:     profile.TwoFactorEnabled.Bool,
		BiometricAuthEnabled: profile.BiometricAuthEnabled,
	}

	switch profile.Userrole {
	case "student":
		protoProfile.RoleSpecificDetails = &user_management.UserProfile_StudentDetails{
			StudentDetails: &user_management.StudentDetails{
				Grade:       profile.Grade.String,
				DateOfBirth: profile.Dateofbirth.Time.Format("2006-01-02"),
				SchoolID:    profile.Schoolid.Int32,
			},
		}
	case "school":
		protoProfile.RoleSpecificDetails = &user_management.UserProfile_SchoolDetails{
			SchoolDetails: &user_management.SchoolDetails{
				SchoolName: profile.Schoolname.String,
				Address:    profile.Schooladdress.String,
				Country:    profile.Country.String,
				Province:   profile.Province.String,
				District:   profile.District.String,
				SchoolType: profile.Schooltype.String,
			},
		}
	case "volunteer":
		protoProfile.RoleSpecificDetails = &user_management.UserProfile_VolunteerDetails{
			VolunteerDetails: &user_management.VolunteerDetails{
				Role:                   profile.Volunteerrole.String,
				GraduateYear:           profile.Graduateyear.Int32,
				SafeGuardCertificate:   profile.Safeguardcertificate,
				HasInternship:          profile.Hasinternship.Bool,
				IsEnrolledInUniversity: profile.Isenrolledinuniversity.Bool,
			},
		}
	}

	return protoProfile
}
