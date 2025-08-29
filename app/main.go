package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func parseEnvConfig() EnvConfig {
	intervalStr := os.Getenv("INTERVAL_SECONDS")
	intervalSeconds, err := strconv.Atoi(intervalStr)
	if err != nil || intervalSeconds < 1 {
		intervalSeconds = 30
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}

	healthBody := os.Getenv("HEALTH_BODY")
	if healthBody == "" {
		healthBody = "{}"
	}

	timeoutStr := os.Getenv("HTTP_TIMEOUT_SECONDS")
	httpTimeoutSeconds, err := strconv.Atoi(timeoutStr)
	if err != nil || httpTimeoutSeconds < 1 {
		httpTimeoutSeconds = 0
	}

	return EnvConfig{
		IntervalSeconds:    intervalSeconds,
		BaseURL:            baseURL,
		HealthBody:         healthBody,
		HTTPTimeoutSeconds: httpTimeoutSeconds,
	}
}

func postHealthCheck(client *http.Client, url string, body string) (HealthResponse, error) {
	var healthResp HealthResponse

	resp, err := client.Post(url+"/health", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return healthResp, err
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return healthResp, err
	}

	err = json.Unmarshal(respBody, &healthResp)
	if err != nil {
		return healthResp, err
	}

	return healthResp, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	config := parseEnvConfig()

	fmt.Printf("Loaded EnvConfig: %+v\n", config)

	var client *http.Client
	if config.HTTPTimeoutSeconds < 1 {
		client = &http.Client{}
	} else {
		client = &http.Client{
			Timeout: time.Duration(config.HTTPTimeoutSeconds) * time.Second,
		}
	}

	ticker := time.NewTicker(time.Duration(config.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("Tick at %v\n", t)

		healthResp, err := postHealthCheck(client, config.BaseURL, config.HealthBody)
		if err != nil {
			fmt.Printf("Error making POST request: %v\n", err)
			continue
		}

		fmt.Printf("Health state: %s, missing devices: %v, reason: %v\n",
			healthResp.State, healthResp.MissingDevices, healthResp.Reason)
	}
}
