package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iRankHub/backend/internal/config"
	"github.com/iRankHub/backend/internal/grpc/server"
	"github.com/iRankHub/backend/internal/utils"
)

func StartBackend() {
    // Load the database configuration from environment variables
    dbConfig, err := config.LoadDatabaseConfig()
    if err != nil {
        log.Fatalf("Failed to load database configuration: %v", err)
    }

    // Connect to the database
    connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.Ssl)
    db, err := sql.Open("pgx", connString)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    log.Println("Successfully connected to the database")
    defer db.Close()

    // Set connection pool settings
    db.SetMaxOpenConns(100)
    db.SetMaxIdleConns(50)
    db.SetConnMaxLifetime(time.Minute * 5)

    // Verify connection
    if err = db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    // Run database migrations only in development
    if os.Getenv("GO_ENV") != "production" {
        migrationPath := "file://internal/database/postgres/migrations"
        m, err := migrate.New(migrationPath, connString)
        if err != nil {
            log.Fatalf("Failed to create migrate instance: %v", err)
        }

        if err := m.Up(); err != nil {
            if err == migrate.ErrNoChange {
                log.Println("No new migrations to apply")
            } else {
                log.Fatalf("Failed to apply migrations: %v", err)
            }
        } else {
            log.Println("Successfully applied database migrations")
        }
    } else {
        log.Println("Skipping migrations in production environment")
    }

    // Initialize the token configuration
    err = utils.InitializeTokenConfig()
    if err != nil {
        log.Fatalf("Failed to initialize token configuration: %v", err)
    }

    // Start the token cleanup goroutine
    utils.StartTokenCleanup()

    // Start the gRPC server
    if err := server.StartGRPCServer(db); err != nil {
        log.Fatalf("Failed to start gRPC server: %v", err)
    }
}