package utils

import (
	"context"
	"github.com/iRankHub/backend/internal/grpc/proto/common"
)

// GetClientMetadata extracts client metadata from context
func GetClientMetadata(ctx context.Context) *common.ClientMetadata {
	if meta, ok := ctx.Value("client_metadata").(*common.ClientMetadata); ok {
		return meta
	}
	return &common.ClientMetadata{}
}

// WithAttemptCount creates new metadata with updated attempt count
func WithAttemptCount(ctx context.Context, count int32) context.Context {
	meta := GetClientMetadata(ctx)
	meta.AttemptCount = count
	return context.WithValue(ctx, "client_metadata", meta)
}

// GetIPAddress returns client IP address
func GetIPAddress(ctx context.Context) string {
	return GetClientMetadata(ctx).IpAddress
}

// GetDeviceInfo returns client device info
func GetDeviceInfo(ctx context.Context) string {
	return GetClientMetadata(ctx).DeviceInfo
}

// GetUserAgent returns client user agent
func GetUserAgent(ctx context.Context) string {
	return GetClientMetadata(ctx).UserAgent
}
