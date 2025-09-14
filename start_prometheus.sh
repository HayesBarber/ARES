#!/bin/bash

docker run -d \
  --name prometheus \
  --restart unless-stopped \
  --network host \
  -v prometheus-volume:/prometheus \
  -v ./prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:latest \
  --config.file=/etc/prometheus/prometheus.yml
