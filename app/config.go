package main

import (
	"os"
	"strconv"
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
