package server

import (
	"context"

	"github.com/iRankHub/backend/internal/grpc/proto/system_health"
	services "github.com/iRankHub/backend/internal/services/system_health"
)

type SystemHealthServer struct {
	system_health.UnimplementedSystemHealthServiceServer
	service *services.SystemHealthService
}

func NewSystemHealthServer() (*SystemHealthServer, error) {
	service, err := services.NewSystemHealthService()
	if err != nil {
		return nil, err
	}
	return &SystemHealthServer{service: service}, nil
}

func (s *SystemHealthServer) GetSystemHealth(ctx context.Context, req *system_health.GetSystemHealthRequest) (*system_health.GetSystemHealthResponse, error) {
	metrics, err := s.service.GetSystemHealth(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &system_health.GetSystemHealthResponse{
		CpuUsagePercentage:         metrics.CPUUsagePercentage,
		MemoryUsagePercentage:      metrics.MemoryUsagePercentage,
		EphemeralStorageUsed:       metrics.EphemeralStorageUsed,
		EphemeralStorageTotal:      metrics.EphemeralStorageTotal,
		EphemeralStoragePercentage: metrics.EphemeralStoragePercentage,
		PvcStorageUsed:             metrics.PVCStorageUsed,
		PvcStorageTotal:            metrics.PVCStorageTotal,
		PvcStoragePercentage:       metrics.PVCStoragePercentage,
		NodeCount:                  int32(metrics.NodeCount),
		PodCount:                   int32(metrics.PodCount),
		PvcCount:                   int32(metrics.PVCCount),
	}, nil
}

func (s *SystemHealthServer) Check(ctx context.Context, req *system_health.HealthCheckRequest) (*system_health.HealthCheckResponse, error) {
	// Simple health check that always returns SERVING if the service is running
	return &system_health.HealthCheckResponse{
		Status: system_health.HealthCheckResponse_SERVING,
	}, nil
}
