package sserver

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
		EphemeralStoragePercentage: metrics.EphemeralStoragePercentage,
		PvcStoragePercentage:       metrics.PVCStoragePercentage,
		NodeCount:                  int32(metrics.NodeCount),
		PodCount:                   int32(metrics.PodCount),
		PvcCount:                   int32(metrics.PVCCount),
	}, nil
}
