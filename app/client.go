package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
