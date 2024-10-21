package app

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/iRankHub/backend/envoy"
	"github.com/iRankHub/backend/internal/config"
	"github.com/iRankHub/backend/internal/grpc/server"
	"github.com/iRankHub/backend/internal/utils"
	"github.com/iRankHub/backend/scripts"
)

var envSetup sync.Once

func StartBackend() {
	// Check if we're in development mode
	isDevMode, err := checkDevMode()
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No .env file found. Assuming production mode.")
			isDevMode = false
		} else {
			log.Fatalf("Error checking development mode: %v", err)
		}
	}

	if isDevMode {
		// Use sync.Once to ensure this only runs once
		envSetup.Do(func() {
			err := script.SetEnvVars()
			if err != nil {
				log.Fatalf("Failed to set environment variables: %v", err)
			} else {
				fmt.Println("Environment variables set successfully")
			}

			// Format Go code
			fmt.Println("Formatting Go code...")
			err = script.FormatGoCode()
			if err != nil {
				log.Printf("Failed to format Go code: %v", err)
			} else {
				fmt.Println("Go code formatted successfully")
			}
		})

		// Start Envoy proxy in development mode
		err := envoy.StartEnvoyProxy()
		if err != nil {
			log.Printf("Warning: Failed to start Envoy proxy: %v", err)
			// Note: We're logging this as a warning instead of fatally exiting
		} else {
			fmt.Println("Envoy proxy started or was already running")
		}
	} else {
		log.Println("Running in production mode")
	}

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
	if isDevMode {
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

func checkDevMode() (bool, error) {
	file, err := os.Open(".env")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "GO_ENV=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && strings.TrimSpace(parts[1]) == "development" {
				return true, nil
			}
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading .env file: %w", err)
	}

	return false, nil
}
