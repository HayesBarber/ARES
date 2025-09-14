docker run -d \
  --name my-prometheus \
  --restart unless-stopped \
  -p 9090:9090 \
  -v prometheus-volume:/prometheus \
  -v "$(pwd)"/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus:latest \
  --config.file=/etc/prometheus/prometheus.yml
