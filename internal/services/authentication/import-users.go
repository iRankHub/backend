package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/services/notification"
	"github.com/iRankHub/backend/internal/services/notification/models"
	"github.com/iRankHub/backend/internal/utils"
)

type ImportUsersService struct {
	signUpService       *SignUpService
	notificationService *notification.Service
}

func NewImportUsersService(signUpService *SignUpService, ns *notification.Service) *ImportUsersService {
	return &ImportUsersService{
		signUpService:       signUpService,
		notificationService: ns,
	}
}

func (s *ImportUsersService) BatchImportUsers(ctx context.Context, users []*authentication.UserData) (int32, []string) {
	var (
		importedCount int32
		failedEmails  []string
		mu            sync.Mutex
		wg            sync.WaitGroup
	)

	for _, userData := range users {
		wg.Add(1)
		go func(userData *authentication.UserData) {
			defer wg.Done()

			additionalInfo := map[string]interface{}{
				"dateOfBirth":            userData.DateOfBirth,
				"schoolID":               userData.SchoolID,
				"schoolName":             userData.SchoolName,
				"address":                userData.Address,
				"country":                userData.Country,
				"province":               userData.Province,
				"district":               userData.District,
				"contactEmail":           userData.ContactEmail,
				"schoolType":             userData.SchoolType,
				"roleInterestedIn":       userData.RoleInterestedIn,
				"graduationYear":         userData.GraduationYear,
				"grade":                  userData.Grade,
				"hasInternship":          userData.HasInternship,
				"isEnrolledInUniversity": userData.IsEnrolledInUniversity,
			}

			password := utils.GenerateRandomPassword()

			err := s.signUpService.SignUp(
				ctx,
				userData.FirstName,
				userData.LastName,
				userData.Email,
				password,
				userData.UserRole,
				userData.Gender,
				userData.NationalID,
				userData.SafeguardingCertificateUrl,
				additionalInfo,
			)

			if err != nil {
				mu.Lock()
				failedEmails = append(failedEmails, userData.Email)
				mu.Unlock()
			} else {
				mu.Lock()
				importedCount++
				mu.Unlock()

				// Send welcome notification with temp password
				go func() {
					metadata := models.AuthMetadata{
						DeviceInfo:   "Batch Import",
						Location:     "Account Creation",
						LastAttempt:  time.Now(),
						AttemptCount: 0,
						IPAddress:    "system",
					}

					// First, send temp password notification
					if err := s.notificationService.SendAccountCreation(
						context.Background(),
						userData.Email,
						s.mapUserRole(userData.UserRole),
						metadata,
					); err != nil {
						log.Printf("Failed to send account creation notification to %s: %v", userData.Email, err)
					}

					// Then send security alert with temp password
					alertMsg := "Your temporary password is: " + password + ". Please change it upon first login."
					if err := s.notificationService.SendSecurityAlert(
						context.Background(),
						userData.Email,
						s.mapUserRole(userData.UserRole),
						alertMsg,
						metadata,
					); err != nil {
						log.Printf("Failed to send temporary password notification to %s: %v", userData.Email, err)
					}
				}()
			}
		}(userData)
	}

	wg.Wait()

	return importedCount, failedEmails
}

// mapUserRole converts proto user role to models.UserRole
func (s *ImportUsersService) mapUserRole(role string) models.UserRole {
	switch role {
	case "admin":
		return models.AdminRole
	case "school":
		return models.SchoolRole
	case "student":
		return models.StudentRole
	case "volunteer":
		return models.VolunteerRole
	default:
		return models.UnspecifiedRole
	}
}
