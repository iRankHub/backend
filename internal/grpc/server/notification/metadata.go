package server

import (
	"fmt"
	"time"

	pb "github.com/iRankHub/backend/internal/grpc/proto/notification"
	"github.com/iRankHub/backend/internal/services/notification/models"
)

// Proto to Service conversions
func protoToAuthMetadata(proto *pb.AuthMetadata) *models.AuthMetadata {
	lastAttempt, _ := time.Parse(time.RFC3339, proto.GetLastAttempt())
	return &models.AuthMetadata{
		IPAddress:    proto.GetIpAddress(),
		DeviceInfo:   proto.GetDeviceInfo(),
		Location:     proto.GetLocation(),
		AttemptCount: int(proto.GetAttemptCount()),
		LastAttempt:  lastAttempt,
	}
}

func protoToUserMetadata(proto *pb.UserMetadata) *models.UserMetadata {
	approvedAt, _ := time.Parse(time.RFC3339, proto.GetApprovedAt())
	expirationDate, _ := time.Parse(time.RFC3339, proto.GetExpirationDate())
	return &models.UserMetadata{
		Changes:        proto.GetChanges(),
		PreviousRole:   proto.GetPreviousRole(),
		NewRole:        proto.GetNewRole(),
		Reason:         proto.GetReason(),
		ApprovedBy:     proto.GetApprovedBy(),
		ApprovedAt:     approvedAt,
		ExpirationDate: &expirationDate,
	}
}

func protoToTournamentMetadata(proto *pb.TournamentMetadata) *models.TournamentMetadata {
	startDate, _ := time.Parse(time.RFC3339, proto.GetStartDate())
	endDate, _ := time.Parse(time.RFC3339, proto.GetEndDate())
	return &models.TournamentMetadata{
		TournamentID:   proto.GetTournamentId(),
		TournamentName: proto.GetTournamentName(),
		StartDate:      startDate,
		EndDate:        endDate,
		Location:       proto.GetLocation(),
		Format:         proto.GetFormat(),
		League:         proto.GetLeague(),
		Fee:            proto.GetFee(),
		Currency:       proto.GetCurrency(),
		Coordinator:    proto.GetCoordinator(),
	}
}

func protoToDebateMetadata(proto *pb.DebateMetadata) *models.DebateMetadata {
	startTime, _ := time.Parse(time.RFC3339, proto.GetStartTime())
	endTime, _ := time.Parse(time.RFC3339, proto.GetEndTime())
	return &models.DebateMetadata{
		DebateID:      proto.GetDebateId(),
		TournamentID:  proto.GetTournamentId(),
		RoundNumber:   int(proto.GetRoundNumber()),
		IsElimination: proto.GetIsElimination(),
		Room:          proto.GetRoom(),
		StartTime:     startTime,
		EndTime:       endTime,
		Team1:         proto.GetTeam1(),
		Team2:         proto.GetTeam2(),
		JudgePanel:    proto.GetJudgePanel(),
		HeadJudge:     proto.GetHeadJudge(),
		Motion:        proto.GetMotion(),
	}
}

func protoToReportMetadata(proto *pb.ReportMetadata) *models.ReportMetadata {
	generatedAt, _ := time.Parse(time.RFC3339, proto.GetGeneratedAt())
	expiresAt, _ := time.Parse(time.RFC3339, proto.GetExpiresAt())

	// Convert map[string]string to map[string]any
	keyMetrics := make(map[string]any)
	for k, v := range proto.GetKeyMetrics() {
		keyMetrics[k] = v // Value will be stored as interface{}
	}

	return &models.ReportMetadata{
		ReportID:    proto.GetReportId(),
		ReportType:  proto.GetReportType(),
		Period:      proto.GetPeriod(),
		GeneratedAt: generatedAt,
		GeneratedBy: proto.GetGeneratedBy(),
		Size:        proto.GetSize(),
		DownloadURL: proto.GetDownloadUrl(),
		Summary:     proto.GetSummary(),
		KeyMetrics:  keyMetrics,
		ExpiresAt:   expiresAt,
		FileSize:    proto.GetFileSize(),
	}
}

// Service to Proto conversions
func authMetadataToProto(meta *models.AuthMetadata) *pb.AuthMetadata {
	return &pb.AuthMetadata{
		IpAddress:    meta.IPAddress,
		DeviceInfo:   meta.DeviceInfo,
		Location:     meta.Location,
		AttemptCount: int32(meta.AttemptCount),
		LastAttempt:  meta.LastAttempt.Format(time.RFC3339),
	}
}

func userMetadataToProto(meta *models.UserMetadata) *pb.UserMetadata {
	expirationDate := ""
	if meta.ExpirationDate != nil {
		expirationDate = meta.ExpirationDate.Format(time.RFC3339)
	}

	return &pb.UserMetadata{
		Changes:        meta.Changes,
		PreviousRole:   meta.PreviousRole,
		NewRole:        meta.NewRole,
		Reason:         meta.Reason,
		ApprovedBy:     meta.ApprovedBy,
		ApprovedAt:     meta.ApprovedAt.Format(time.RFC3339),
		ExpirationDate: expirationDate,
	}
}

func tournamentMetadataToProto(meta *models.TournamentMetadata) *pb.TournamentMetadata {
	return &pb.TournamentMetadata{
		TournamentId:   meta.TournamentID,
		TournamentName: meta.TournamentName,
		StartDate:      meta.StartDate.Format(time.RFC3339),
		EndDate:        meta.EndDate.Format(time.RFC3339),
		Location:       meta.Location,
		Format:         meta.Format,
		League:         meta.League,
		Fee:            meta.Fee,
		Currency:       meta.Currency,
		Coordinator:    meta.Coordinator,
	}
}

func debateMetadataToProto(meta *models.DebateMetadata) *pb.DebateMetadata {
	return &pb.DebateMetadata{
		DebateId:      meta.DebateID,
		TournamentId:  meta.TournamentID,
		RoundNumber:   int32(meta.RoundNumber),
		IsElimination: meta.IsElimination,
		Room:          meta.Room,
		StartTime:     meta.StartTime.Format(time.RFC3339),
		EndTime:       meta.EndTime.Format(time.RFC3339),
		Team1:         meta.Team1,
		Team2:         meta.Team2,
		JudgePanel:    meta.JudgePanel,
		HeadJudge:     meta.HeadJudge,
		Motion:        meta.Motion,
	}
}

func reportMetadataToProto(meta *models.ReportMetadata) *pb.ReportMetadata {
	// Convert map[string]any to map[string]string
	keyMetrics := make(map[string]string)
	for k, v := range meta.KeyMetrics {
		// Convert any value to string
		keyMetrics[k] = fmt.Sprintf("%v", v)
	}

	return &pb.ReportMetadata{
		ReportId:    meta.ReportID,
		ReportType:  meta.ReportType,
		Period:      meta.Period,
		GeneratedAt: meta.GeneratedAt.Format(time.RFC3339),
		GeneratedBy: meta.GeneratedBy,
		Size:        meta.Size,
		DownloadUrl: meta.DownloadURL,
		Summary:     meta.Summary,
		KeyMetrics:  keyMetrics,
		ExpiresAt:   meta.ExpiresAt.Format(time.RFC3339),
		FileSize:    meta.FileSize,
	}
}
