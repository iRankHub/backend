package server

import (
	"context"
	"database/sql"

	"github.com/iRankHub/backend/internal/grpc/proto/analytics"
	services "github.com/iRankHub/backend/internal/services/analytics"
)

type AnalyticsServer struct {
	analytics.UnimplementedAnalyticsServiceServer
	service *services.AnalyticsService
}

func NewAnalyticsServer(db *sql.DB) (*AnalyticsServer, error) {
	service := services.NewAnalyticsService(db)
	return &AnalyticsServer{
		service: service,
	}, nil
}

func (s *AnalyticsServer) GetFinancialReports(ctx context.Context, req *analytics.FinancialReportRequest) (*analytics.FinancialReportResponse, error) {
	return s.service.GetFinancialReports(ctx, req)
}

func (s *AnalyticsServer) GetAttendanceReports(ctx context.Context, req *analytics.AttendanceReportRequest) (*analytics.AttendanceReportResponse, error) {
	return s.service.GetAttendanceReports(ctx, req)
}
