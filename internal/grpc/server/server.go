package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/grpc/proto/authentication"
	"github.com/iRankHub/backend/internal/grpc/proto/debate_management"
	"github.com/iRankHub/backend/internal/grpc/proto/notification"
	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/grpc/proto/user_management"
	authserver "github.com/iRankHub/backend/internal/grpc/server/authentication"
	debateserver "github.com/iRankHub/backend/internal/grpc/server/debate_management"
	notificationserver "github.com/iRankHub/backend/internal/grpc/server/notification"
	tournamentserver "github.com/iRankHub/backend/internal/grpc/server/tournament_management"
	userserver "github.com/iRankHub/backend/internal/grpc/server/user_management"
)

func StartGRPCServer(db *sql.DB) error {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create and register all your servers
	authServer, err := authserver.NewAuthServer(db)
	if err != nil {
		return fmt.Errorf("failed to create AuthServer: %v", err)
	}
	authentication.RegisterAuthServiceServer(grpcServer, authServer)

	userManagementServer, err := userserver.NewUserManagementServer(db)
	if err != nil {
		return fmt.Errorf("failed to create UserManagementServer: %v", err)
	}
	user_management.RegisterUserManagementServiceServer(grpcServer, userManagementServer)

	tournamentServer, err := tournamentserver.NewTournamentServer(db)
	if err != nil {
		return fmt.Errorf("failed to create TournamentServer: %v", err)
	}
	tournament_management.RegisterTournamentServiceServer(grpcServer, tournamentServer)

	debateServer, err := debateserver.NewDebateServer(db)
	if err != nil {
		return fmt.Errorf("failed to create DebateServer: %v", err)
	}
	debate_management.RegisterDebateServiceServer(grpcServer, debateServer)

	notificationServer, err := notificationserver.NewNotificationServer(db)
	if err != nil {
		return fmt.Errorf("failed to create NotificationServer: %v", err)
	}
	notification.RegisterNotificationServiceServer(grpcServer, notificationServer)

	// Read the gRPC server port from the environment
	viper.SetDefault("GRPC_PORT", "8080")
	grpcPort := viper.GetString("GRPC_PORT")

	// Start the gRPC server on the specified port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server started on 0.0.0.0:%s", grpcPort)
	return grpcServer.Serve(listener)
}
