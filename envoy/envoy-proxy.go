package envoy

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

func StartEnvoyProxy() error {
	// Load configuration from .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %v", err)
	}

	// Run the Envoy config generator
	log.Println("Generating Envoy configuration...")
	genCmd := exec.Command("go", "run", "./cmd/envoy/envoy.go")
	if output, err := genCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to generate Envoy config: %v\nOutput: %s", err, output)
	}

	// Get Envoy Proxy configuration from Viper
	envoyContainerName := viper.GetString("ENVOY_CONTAINER_NAME")
	envoyListenerPort := viper.GetString("ENVOY_LISTENER_PORT")
	envoyAdminPort := viper.GetString("ENVOY_ADMIN_PORT")

	// Check if the Envoy Proxy container is already running
	checkCmd := exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", envoyContainerName))
	output, _ := checkCmd.Output()
	if len(output) > 0 {
		log.Printf("Envoy Proxy container '%s' is already running", envoyContainerName)
		return nil
	}

	// Build the Envoy Proxy image using the Dockerfile
	log.Println("Building Envoy Proxy image...")
	buildCmd := exec.Command("docker", "build", "-t", envoyContainerName, ".")
	buildCmd.Stdout = log.Writer()
	buildCmd.Stderr = log.Writer()
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Envoy Proxy image: %v", err)
	}

	// Create and start the Envoy Proxy container
	log.Printf("Starting Envoy Proxy container: %s", envoyContainerName)
	configPath, _ := filepath.Abs("./envoy.yaml")
	runCmd := exec.Command("docker", "run", "-d", "--name", envoyContainerName,
		"-p", fmt.Sprintf("%s:%s", envoyListenerPort, envoyListenerPort),
		"-p", fmt.Sprintf("%s:%s", envoyAdminPort, envoyAdminPort),
		"-v", fmt.Sprintf("%s:/etc/envoy/envoy.yaml", configPath),
		envoyContainerName)

	output, err = runCmd.CombinedOutput()
	if err != nil {
		log.Printf("Docker run command: %s", runCmd.String())
		log.Printf("Docker run output: %s", string(output))
		return fmt.Errorf("failed to start Envoy Proxy container: %v", err)
	}

	// Check if the container is actually running
	checkCmd = exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", envoyContainerName))
	checkOutput, checkErr := checkCmd.Output()
	if checkErr != nil || len(checkOutput) == 0 {
		log.Printf("Container created but not running. Checking logs...")
		logCmd := exec.Command("docker", "logs", envoyContainerName)
		logOutput, _ := logCmd.CombinedOutput()
		log.Printf("Container logs: %s", string(logOutput))
		return fmt.Errorf("container created but not running: %v", checkErr)
	}

	log.Printf("Envoy Proxy container started: %s", envoyContainerName)
	return nil
}