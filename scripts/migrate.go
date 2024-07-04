package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/iRankHub/backend/internal/config"
)

func main() {
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatalf("failed to load database config: %v", err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

		m, err := migrate.New("file://internal/database/postgres/migrations", connString)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}