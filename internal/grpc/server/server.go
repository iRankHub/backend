package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/grpc/proto/user_management"
	authserver "github.com/iRankHub/backend/internal/grpc/server/authentication"
	userserver "github.com/iRankHub/backend/internal/grpc/server/user_management"
	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	tournamentserver "github.com/iRankHub/backend/internal/grpc/server/tournament_management"

)

func StartGRPCServer(db *sql.DB) error {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create the Auth server
	authServer, err := authserver.NewAuthServer(db)
	if err != nil {
		return fmt.Errorf("failed to create AuthServer: %v", err)
	}

	// Register the AuthService with the gRPC server
	authentication.RegisterAuthServiceServer(grpcServer, authServer)

	// Create the User server
	userManagementServer, err := userserver.NewUserManagementServer(db)
	if err != nil {
		return fmt.Errorf("failed to create UserManagementServer: %v", err)
	}
	// Register the UserManagementService with the gRPC server
	user_management.RegisterUserManagementServiceServer(grpcServer, userManagementServer)

		// Create the Tournament server
	tournamentServer, err := tournamentserver.NewTournamentServer(db)
	if err != nil {
		return fmt.Errorf("failed to create TournamentServer: %v", err)
	}
	// Register the TournamentService with the gRPC server
	tournament_management.RegisterTournamentServiceServer(grpcServer, tournamentServer)

	// Start the gRPC server on a specific port
	grpcPort := "0.0.0.0:8080"
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server started on %s", grpcPort)
	return grpcServer.Serve(listener)
}