package server

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"

)

func StartGRPCServer(queries *models.Queries, privateKey ed25519.PrivateKey) error {
    // Create a new gRPC server with logging interceptor
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
            log.Printf("Received gRPC request: %s", info.FullMethod)
            return handler(ctx, req)
        }),
    )

    // Register the AuthService with the gRPC server
    authServer, err := NewAuthServer(queries, privateKey)
    if err != nil {
        return fmt.Errorf("failed to create AuthServer: %v", err)
    }
    proto.RegisterAuthServiceServer(grpcServer, authServer)

    // Start the gRPC server on a specific port
    grpcPort := "0.0.0.0:8080"
    listener, err := net.Listen("tcp", grpcPort)
    if err != nil {
        return fmt.Errorf("failed to listen on port %s: %v", grpcPort, err)
    }
    log.Printf("gRPC server starting on %s", grpcPort)

    // Serve gRPC requests
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Printf("gRPC server failed to serve: %v", err)
    }
    return err
}
