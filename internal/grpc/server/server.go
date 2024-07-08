package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	authserver "github.com/iRankHub/backend/internal/grpc/server/authentication"
)

func StartGRPCServer(db *sql.DB) error {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create the AuthServer
	authServer, err := authserver.NewAuthServer(db)
	if err != nil {
		return fmt.Errorf("failed to create AuthServer: %v", err)
	}

	// Register the AuthService with the gRPC server
	authentication.RegisterAuthServiceServer(grpcServer, authServer)

	// Start the gRPC server on a specific port
	grpcPort := "0.0.0.0:8080"
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server started on %s", grpcPort)
	return grpcServer.Serve(listener)
}