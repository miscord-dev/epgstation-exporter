package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace is the namespace of the metrics
	namespace = "epgstation"
)

// NewCounterVec creates a new CounterVec with "epgstation" as the namespace
func NewCounterVec(opts prometheus.CounterOpts, labels []string) *prometheus.CounterVec {
	opts.Namespace = namespace
	return prometheus.NewCounterVec(opts, labels)
}

// NewGaugeVec creates a new GaugeVec with "epgstation" as the namespace
func NewGaugeVec(opts prometheus.GaugeOpts, labels []string) *prometheus.GaugeVec {
	opts.Namespace = namespace
	return prometheus.NewGaugeVec(opts, labels)
}

// NewHistogramVec creates a new HistogramVec with "epgstation" as the namespace
func NewHistogramVec(opts prometheus.HistogramOpts, labels []string) *prometheus.HistogramVec {
	opts.Namespace = namespace
	return prometheus.NewHistogramVec(opts, labels)
}
