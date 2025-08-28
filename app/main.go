package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
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

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("Tick at %v\n", t)

		resp, err := http.Post(baseURL+"/health", "application/json", bytes.NewBuffer([]byte("{}")))
		if err != nil {
			fmt.Printf("Error making POST request: %v\n", err)
			continue
		}
		resp.Body.Close()
		fmt.Printf("POST /health responded with status: %s\n", resp.Status)
	}
}
