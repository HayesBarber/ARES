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

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	intervalStr := os.Getenv("INTERVAL_SECONDS")
	intervalSeconds, err := strconv.Atoi(intervalStr)
	if err != nil || intervalSeconds < 1 {
		fmt.Println("Invalid interval, setting to default (30 seconds)")
		intervalSeconds = 30
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}

	bodyStr := os.Getenv("HEALTH_BODY")
	if bodyStr == "" {
		bodyStr = "{}"
	}

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("Tick at %v\n", t)

		resp, err := http.Post(baseURL+"/health", "application/json", bytes.NewBuffer([]byte(bodyStr)))
		if err != nil {
			fmt.Printf("Error making POST request: %v\n", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			continue
		}

		fmt.Printf("POST /health responded with status: %s\n", resp.Status)

		var healthResp HealthResponse
		err = json.Unmarshal(respBody, &healthResp)
		if err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			continue
		}
		fmt.Printf("Health state: %s, missing devices: %v, reason: %v\n",
			healthResp.State, healthResp.MissingDevices, healthResp.Reason)
	}
}
