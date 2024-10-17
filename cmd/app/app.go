package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"

	"github.com/iRankHub/backend/internal/config"
	"github.com/iRankHub/backend/internal/grpc/server"
	"github.com/iRankHub/backend/internal/utils"
)

func StartBackend() {
    // Load the database configuration from environment variables
	dbConfig := &config.DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Name:     viper.GetString("DB_NAME"),
		Ssl:      viper.GetString("DB_SSL"),
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
	db.SetMaxOpenConns(100)  // Potential number of concurrent users
	db.SetMaxIdleConns(50)   // Half of max open connections
	db.SetConnMaxLifetime(time.Minute * 5)  // Recycle connections after 5 minutes

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

    // Run database migrations
    migrationPath := "file://internal/database/postgres/migrations"
    m, err := migrate.New(migrationPath, connString)
    if err != nil {
        log.Fatalf("Failed to create migrate instance: %v", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to run database migrations: %v", err)
    }
    log.Println("Successfully ran database migrations")

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