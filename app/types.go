package main

type HealthState string

const (
	Healthy   HealthState = "healthy"
	Moderate  HealthState = "moderate"
	Unhealthy HealthState = "unhealthy"
)

type HealthResponse struct {
	State          HealthState `json:"state"`
	MissingDevices []string    `json:"missing_devices,omitempty"`
	Reason         *string     `json:"reason,omitempty"`
}
