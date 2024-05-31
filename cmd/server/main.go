package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/iRankHub/backend/internal/config"
	"github.com/iRankHub/backend/internal/database/postgres"
	"github.com/iRankHub/backend/internal/grpc/proto"
	"github.com/iRankHub/backend/internal/grpc/server"
	"github.com/iRankHub/backend/internal/models"
)

func main() {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("failed to load database config: %v", err)
	}

	conn, err := postgres.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	err = postgres.CreateDatabase(conn, dbConfig.Name)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	// Reconnect to the newly created database
	conn.Close(context.Background())
	dbConfig.Name = viper.GetString("DB_NAME")
	conn, err = postgres.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatalf("failed to reconnect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Run migrations
	migrationPath := "file://internal/database/postgres/migrations"
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	err = postgres.RunMigrations(connString, migrationPath)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Initialize database connection for the generated models
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("failed to initialize database connection: %v", err)
	}
	defer db.Close()

	// Create a new instance of the Queries struct
	queries := models.New(db)

	// Use the generated query functions from the Queries struct
	// Example usage:
	users, err := queries.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("failed to get users: %v", err)
	}
	log.Printf("Retrieved %d users from the database", len(users))

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the AuthService with the gRPC server
	proto.RegisterAuthServiceServer(grpcServer, server.NewAuthServer(db))

	// Start the gRPC server on a specific port
	grpcPort := ":50051"
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}
	log.Printf("gRPC server started on%s", grpcPort)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}