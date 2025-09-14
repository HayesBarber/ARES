package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	healthyCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_healthy_total",
			Help: "Total number of times the service was healthy",
		},
	)
	moderateCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_moderate_total",
			Help: "Total number of times the service was in moderate state",
		},
	)
	unhealthyCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_unhealthy_total",
			Help: "Total number of times the service was unhealthy",
		},
	)
	missingDevicesCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_missing_devices_total",
			Help: "Total number of times each device was missing",
		},
		[]string{"device"},
	)
	reasonCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_unhealthy_reason_total",
			Help: "Total number of times the service was unhealthy for a given reason",
		},
		[]string{"reason"},
	)
)

func init() {
	prometheus.MustRegister(healthyCount)
	prometheus.MustRegister(moderateCount)
	prometheus.MustRegister(unhealthyCount)
	prometheus.MustRegister(missingDevicesCount)
	prometheus.MustRegister(reasonCount)
}

func RecordHealthMetrics(resp HealthResponse, service string) {
	switch resp.State {
	case Healthy:
		healthyCount.Inc()
	case Moderate:
		moderateCount.Inc()
		for _, device := range resp.MissingDevices {
			missingDevicesCount.WithLabelValues(device).Inc()
		}
	case Unhealthy:
		unhealthyCount.Inc()
		reasonCount.WithLabelValues(*resp.Reason).Inc()
	}
}
