package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Ssl      string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	dbConfig := &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Ssl:      os.Getenv("DB_SSL"),
	}

	return dbConfig, nil
}