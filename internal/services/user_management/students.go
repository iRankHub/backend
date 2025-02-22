package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type StudentService struct {
	db *sql.DB
}

func NewStudentsManagementService(db *sql.DB) *StudentService {
	return &StudentService{db: db}
}

func (s *StudentService) GetStudents(ctx context.Context, token string, page, pageSize int32) ([]models.Student, int32, error) {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	queries := models.New(s.db)
	paginatedStudents, err := queries.GetStudentsPaginated(ctx, models.GetStudentsPaginatedParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch students: %v", err)
	}

	// Convert []models.GetStudentsPaginatedRow to []models.Student
	students := make([]models.Student, len(paginatedStudents))
	for i, s := range paginatedStudents {
		students[i] = models.Student{
			Studentid:   s.Studentid,
			Firstname:   s.Firstname,
			Lastname:    s.Lastname,
			Grade:       s.Grade,
			Dateofbirth: s.Dateofbirth,
			Email:       s.Email,
			Schoolid:    s.Schoolid,
			Userid:      s.Userid,
		}
	}

	totalCount, err := queries.GetTotalStudentCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total student count: %v", err)
	}

	return students, int32(totalCount), nil
}

func (s *StudentService) GetStudentsBySchoolContactID(ctx context.Context, token string, userID int32, page, pageSize int32) ([]models.Student, int32, error) {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid token: %v", err)
	}

	queries := models.New(s.db)

	// First, get the school ID associated with the user ID
	school, err := queries.GetSchoolByContactPersonID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, fmt.Errorf("no school found for the given user ID")
		}
		return nil, 0, fmt.Errorf("failed to fetch school: %v", err)
	}

	// Now, fetch students for this school
	students, err := queries.GetStudentsBySchoolID(ctx, models.GetStudentsBySchoolIDParams{
		Schoolid: school.Schoolid,
		Limit:    pageSize,
		Offset:   (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch students: %v", err)
	}

	// Get total count of students for this school
	totalCount, err := queries.GetStudentCountBySchoolID(ctx, school.Schoolid)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total student count: %v", err)
	}

	return students, int32(totalCount), nil
}
