package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func postHealthCheck(client *http.Client, url string, body string) (HealthResponse, error) {
	var healthResp HealthResponse

	resp, err := client.Post(url+"/health", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return healthResp, err
	}

	if resp.StatusCode >= 400 {
		return healthResp, fmt.Errorf("HTTP error: %s", resp.Status)
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

func checkWithBackoff(client *http.Client, url, body string, maxRetries int) (HealthResponse, error) {
	var healthResp HealthResponse
	var err error

	backoff := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		healthResp, err = postHealthCheck(client, url, body)
		if err != nil {
			return healthResp, err
		}

		if healthResp.State == Healthy {
			return healthResp, nil
		}

		if attempt < maxRetries {
			sleep := backoff + time.Duration(rand.Intn(500))*time.Millisecond
			fmt.Printf("Moderate or Unhealthy state, retrying in %v (attempt %d)\n", sleep, attempt+1)
			time.Sleep(sleep)
			backoff *= 2
		}
	}

	return healthResp, nil
}
