package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type EnvoyConfig struct {
	ENVOY_LISTENER_PORT           int
	ENVOY_ROUTE_TIMEOUT           string
	ENVOY_CORS_MAX_AGE            string
	ENVOY_BACKEND_CONNECT_TIMEOUT string
	BACKEND_SERVICE_HOST          string
	BACKEND_SERVICE_PORT          int
	ENVOY_ADMIN_PORT              int
}

func main() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config := EnvoyConfig{
		ENVOY_LISTENER_PORT:           viper.GetInt("ENVOY_LISTENER_PORT"),
		ENVOY_ROUTE_TIMEOUT:           viper.GetString("ENVOY_ROUTE_TIMEOUT"),
		ENVOY_CORS_MAX_AGE:            viper.GetString("ENVOY_CORS_MAX_AGE"),
		ENVOY_BACKEND_CONNECT_TIMEOUT: viper.GetString("ENVOY_BACKEND_CONNECT_TIMEOUT"),
		BACKEND_SERVICE_HOST:          viper.GetString("BACKEND_SERVICE_HOST"),
		BACKEND_SERVICE_PORT:          viper.GetInt("BACKEND_SERVICE_PORT"),
		ENVOY_ADMIN_PORT:              viper.GetInt("ENVOY_ADMIN_PORT"),
	}

	templatePath, _ := filepath.Abs("./envoy/envoy.yaml.template")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	outputPath, _ := filepath.Abs("./envoy.yaml")
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, config)
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}

	log.Println("Envoy configuration file generated successfully")
}
