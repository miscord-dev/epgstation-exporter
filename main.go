package main

import (
	"flag"
	"github.com/miscord-dev/epgstation-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func main() {
	var (
		maxRetry = flag.Int("max-retry", 3, "max retry for the exporter")
		baseURL  = flag.String("base-url", "http://localhost:8888/api", "base URL for the exporter")
	)
	flag.Parse()
	e, err := metrics.New(
		metrics.WithBaseURL(*baseURL),
		metrics.WithMaxRetry(*maxRetry),
	)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
