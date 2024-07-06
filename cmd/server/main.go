package main

import (
	"database/sql"
	"fmt"
	"time"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"

	"github.com/iRankHub/backend/envoy"
	"github.com/iRankHub/backend/internal/config"
	"github.com/iRankHub/backend/internal/grpc/server"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"

)

func main() {
	// Load the database configuration using Viper
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	dbConfig := &config.DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Name:     viper.GetString("DB_NAME"),
	}

	// Connect to the database
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database")
	defer db.Close()

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

	// Initialize database connection for the generated models
	queries := models.New(db)

	// Generate PASETO key pair
	privateKey, publicKey, err := utils.GeneratePasetoKeyPair()
	if err != nil {
		log.Fatalf("Failed to generate PASETO key pair: %v", err)
	}

	// Set the public key for token validation
	utils.SetPublicKey(publicKey)

	// Start the gRPC server
	go func() {
		if err := server.StartGRPCServer(queries, privateKey); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Give the gRPC server a moment to start
	time.Sleep(2 * time.Second)

	// Start the Envoy Proxy server
	if err := envoy.StartEnvoyProxy(); err != nil {
		log.Fatalf("Failed to start Envoy Proxy: %v", err)
	}
}