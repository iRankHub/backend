package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	listenerPort, _ := strconv.Atoi(os.Getenv("ENVOY_LISTENER_PORT"))
	backendPort, _ := strconv.Atoi(os.Getenv("BACKEND_SERVICE_PORT"))
	adminPort, _ := strconv.Atoi(os.Getenv("ENVOY_ADMIN_PORT"))

	config := EnvoyConfig{
		ENVOY_LISTENER_PORT:           listenerPort,
		ENVOY_ROUTE_TIMEOUT:           os.Getenv("ENVOY_ROUTE_TIMEOUT"),
		ENVOY_CORS_MAX_AGE:            os.Getenv("ENVOY_CORS_MAX_AGE"),
		ENVOY_BACKEND_CONNECT_TIMEOUT: os.Getenv("ENVOY_BACKEND_CONNECT_TIMEOUT"),
		BACKEND_SERVICE_HOST:          os.Getenv("BACKEND_SERVICE_HOST"),
		BACKEND_SERVICE_PORT:          backendPort,
		ENVOY_ADMIN_PORT:              adminPort,
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