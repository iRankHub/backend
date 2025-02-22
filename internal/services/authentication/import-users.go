package services

import (
	"context"
	"log"
	"sync"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	notificationService "github.com/iRankHub/backend/internal/services/notification"
	"github.com/iRankHub/backend/internal/utils"
	notification "github.com/iRankHub/backend/internal/utils/notifications"
)

type ImportUsersService struct {
	signUpService       *SignUpService
	notificationService *notificationService.NotificationService
}

func NewImportUsersService(signUpService *SignUpService, ns *notificationService.NotificationService) *ImportUsersService {
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

				// Send email with temporary password
				go func() {
					if err := notification.SendTemporaryPasswordEmail(s.notificationService, userData.Email, userData.FirstName, password); err != nil {
						log.Printf("Failed to send temporary password email to %s: %v", userData.Email, err)
					}
				}()
			}
		}(userData)
	}

	wg.Wait()

	return importedCount, failedEmails
}
