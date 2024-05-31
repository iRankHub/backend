package postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"

	"github.com/iRankHub/backend/internal/config"
)

func ConnectDatabase(cfg *config.DatabaseConfig) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return conn, nil
}

func CreateDatabase(conn *pgx.Conn, dbName string) error {
	// Check if the database exists
	var exists bool
	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %v", err)
	}

	if !exists {
		// Create the database if it doesn't exist
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		fmt.Printf("Database '%s' created successfully\n", dbName)
	} else {
		fmt.Printf("Database '%s' already exists\n", dbName)
	}

	return nil
}

func RunMigrations(connString string, migrationPath string) error {
	m, err := migrate.New(migrationPath, connString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}
	return nil
}
