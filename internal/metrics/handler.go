package metrics

import (
	"context"
	"fmt"
	"github.com/miscord-dev/epgstation-exporter/api/epgstation"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

const (
	namespace       = "epgstation"
	defaultBaseURL  = "http://localhost:8888/api"
	defaultMaxRetry = 3
)

var (
	epgStationRuleReserves = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "rule_reserves"),
		"The number of reserves by rule",
		[]string{"rule_id", "is_time_specification"},
		nil,
	)
)

type exporter struct {
	c        *epgstation.Client
	baseURL  string
	maxRetry int
	logger   *slog.Logger
}

type Option func(*exporter)

// WithBaseURL sets the base URL for the exporter
func WithBaseURL(baseURL string) Option {
	return func(e *exporter) {
		e.baseURL = baseURL
	}
}

// WithMaxRetry sets the max retry for the exporter
func WithMaxRetry(maxRetry int) Option {
	return func(e *exporter) {
		e.maxRetry = maxRetry
	}
}

// WithLogger sets the logger for the exporter
func WithLogger(logger *slog.Logger) Option {
	return func(e *exporter) {
		e.logger = logger
	}
}

// New returns a new EPGStation exporter
func New(opts ...Option) (prometheus.Collector, error) {
	e := &exporter{
		baseURL:  defaultBaseURL,
		maxRetry: defaultMaxRetry,
		logger:   slog.Default(),
	}

	for _, opt := range opts {
		opt(e)
	}

	c, err := epgstation.NewClient(e.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EPGStation API client: %w", err)
	}
	e.c = c

	return e, nil
}

// Describe implements the prometheus.Collector interface.
func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- epgStationRuleReserves
}

// Collect implements the prometheus.Collector interface.
func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	e.logger.Info("Collecting metrics")

	if err := e.collectMetrics(ch); err != nil {
		e.logger.Error("Failed to collect metrics", slog.String("error", err.Error()))
		// TODO: count errors and expose it as metrics
		return
	}

	e.logger.Info("Collected metrics")
}

// collectMetrics collects metrics from EPGStation
func (e *exporter) collectMetrics(ch chan<- prometheus.Metric) error {
	e.logger.Debug("Collecting metrics from EPGStation")

	var rules []epgstation.Rule
	var err error
	for i := 0; i < e.maxRetry; i++ {
		e.logger.Debug("Getting rules from EPGStation", slog.Int("attempt", i))
		rules, err = getEPGStationRules(e.c)
		if err != nil {
			e.logger.Warn("Failed to get rules from EPGStation, retrying...", slog.String("error", err.Error()))
		}
	}
	if err != nil {
		e.logger.Error("Failed to get rules from EPGStation", slog.String("error", err.Error()))
		return err
	}

	e.logger.Debug("Got rules from EPGStation", rules)

	for _, rule := range rules {
		e.logger.Debug("Collecting metrics for rule", rule)
		ch <- prometheus.MustNewConstMetric(
			epgStationRuleReserves,
			prometheus.GaugeValue,
			float64(derefIntPtr(rule.ReservesCnt, 0)),
			fmt.Sprintf("%d", rule.Id),
			fmt.Sprintf("%t", rule.IsTimeSpecification),
		)
	}

	e.logger.Debug("Collected metrics from EPGStation")

	return nil
}

// getEPGStationRules returns the rules obtained via EPGStation API
// TODO(musaprg): Refactor this function, move to other package
func getEPGStationRules(c *epgstation.Client) ([]epgstation.Rule, error) {
	r, err := c.GetRules(context.Background(), &epgstation.GetRulesParams{
		Offset:  nil,
		Limit:   intPtr(0), // get all rules
		Type:    getReserveTypePtr("all"),
		Keyword: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rules: %w", err)
	}

	res, err := epgstation.ParseGetRulesResponse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("request failed with %d: %+v", res.StatusCode(), res.JSONDefault)
	}

	return res.JSON200.Rules, nil
}
