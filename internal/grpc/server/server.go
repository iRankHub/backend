package server

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/models"
)

func StartGRPCServer(queries *models.Queries, privateKey ed25519.PrivateKey) error {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create the AuthServer
	authServer, err := NewAuthServer(queries, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create AuthServer: %v", err)
	}

	// Register the AuthService with the gRPC server
	proto.RegisterAuthServiceServer(grpcServer, authServer)

	// Start the gRPC server on a specific port
	grpcPort := "0.0.0.0:8080"  // Changed to listen on all interfaces
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server started on %s", grpcPort)
	return grpcServer.Serve(listener)
}