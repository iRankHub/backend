package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Ssl		 string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	dbConfig := &DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetInt("DB_PORT"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		Name:     viper.GetString("DB_NAME"),
		Ssl: 	  viper.GetString("DB_SSL"), // Optional: Set to "disable" to disable SSL/TLS connection.
	}

	return dbConfig, nil
}
