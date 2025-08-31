package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

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

		healthResp, err := checkWithBackoff(client, config.BaseURL, config.HealthBody, config.MaxRetries)
		if err != nil {
			fmt.Printf("Error making POST request: %v\n", err)
			continue
		}

		reason := "<nil>"
		if healthResp.Reason != nil {
			reason = *healthResp.Reason
		}

		fmt.Printf("Health state: %s, missing devices: %v, reason: %s\n",
			healthResp.State, healthResp.MissingDevices, reason)
	}
}
