package interceptors

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/iRankHub/backend/internal/grpc/proto/common"
)

type MetadataInterceptor struct{}

func NewMetadataInterceptor() *MetadataInterceptor {
	return &MetadataInterceptor{}
}

// UnaryServerInterceptor intercepts unary RPC calls
func (i *MetadataInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientMeta := extractClientMetadata(ctx)
		newCtx := context.WithValue(ctx, "client_metadata", clientMeta)
		return handler(newCtx, req)
	}
}

// StreamServerInterceptor intercepts stream RPC calls
func (i *MetadataInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientMeta := extractClientMetadata(ss.Context())
		newCtx := context.WithValue(ss.Context(), "client_metadata", clientMeta)
		wrappedStream := newWrappedServerStream(ss, newCtx)
		return handler(srv, wrappedStream)
	}
}

func extractClientMetadata(ctx context.Context) *common.ClientMetadata {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &common.ClientMetadata{}
	}

	// Get client info from metadata
	clientInfoValues := md.Get("x-client-info")
	if len(clientInfoValues) == 0 {
		// Try to get individual headers
		return &common.ClientMetadata{
			IpAddress:  getFirstValue(md, "x-real-ip", "x-forwarded-for"),
			UserAgent:  getFirstValue(md, "user-agent"),
			DeviceInfo: getFirstValue(md, "sec-ch-ua"),
		}
	}

	// Parse client info JSON
	var clientInfo common.ClientMetadata
	if err := json.Unmarshal([]byte(clientInfoValues[0]), &clientInfo); err != nil {
		return &common.ClientMetadata{}
	}

	return &clientInfo
}

func getFirstValue(md metadata.MD, keys ...string) string {
	for _, key := range keys {
		if values := md.Get(key); len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// wrappedServerStream wraps grpc.ServerStream to modify context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func newWrappedServerStream(ss grpc.ServerStream, ctx context.Context) *wrappedServerStream {
	return &wrappedServerStream{
		ServerStream: ss,
		ctx:          ctx,
	}
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
