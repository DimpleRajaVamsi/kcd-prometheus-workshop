package metrics

import (
	"context"
	"net"
	"net/http"

	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var registry *prometheus.Registry

var apiInvokeCounter *prometheus.CounterVec
var connectionsGauge prometheus.Gauge
var delayHistogram prometheus.Histogram
var delaySummary prometheus.Summary

const (
	apiInvokeCounterName = "api_invoke_count"
	apiInvokeCounterDesc = "How many times each API is invoked"

	connectionsGaugeName = "present_connections"
	connectionsGaugeDesc = "How many active http connections"

	delayHistogramName = "delay_histogram"
	delayHistogramDesc = "Histogram distrubution of the delay in seconds"

	delaySummaryName = "delay_summary"
	delaySummaryDesc = "Summary of the delay in seconds"
)

func IntMetrics(ctx context.Context) {
	logger := utils.GetLogger(ctx)
	if registry == nil {
		registry = prometheus.NewPedanticRegistry()
		logger.Info("Created new Prometheus registry")
	}
	initApiInvokeCountCounter(ctx)
	initConnectionsMetric(ctx)
	initDelayHistogram(ctx)
	initDelaySummary(ctx)
}

func initApiInvokeCountCounter(ctx context.Context) {
	logger := utils.GetLogger(ctx)
	if apiInvokeCounter == nil {
		opts := prometheus.Opts{
			Name: apiInvokeCounterName,
			Help: apiInvokeCounterDesc,
		}
		apiInvokeCounter = prometheus.NewCounterVec(prometheus.CounterOpts(opts), []string{"api_name", "success"})
		apiInvokeCounter.WithLabelValues("beer", "true")
		apiInvokeCounter.WithLabelValues("beer", "false")
		apiInvokeCounter.WithLabelValues("car", "true")
		apiInvokeCounter.WithLabelValues("car", "false")
		apiInvokeCounter.WithLabelValues("delay", "true")
		apiInvokeCounter.WithLabelValues("delay", "false")
		if err := registry.Register(apiInvokeCounter); err != nil {
			logger.Error("Failed to registry the metric", zap.Error(err))
		}
		logger.Info("Created new metric", zap.String("name", apiInvokeCounterName))
	}
}

func initConnectionsMetric(ctx context.Context) {
	logger := utils.GetLogger(ctx)
	if connectionsGauge == nil {
		opts := prometheus.Opts{
			Name: connectionsGaugeName,
			Help: connectionsGaugeDesc,
		}
		connectionsGauge = prometheus.NewGauge(prometheus.GaugeOpts(opts))
		if err := registry.Register(connectionsGauge); err != nil {
			logger.Error("Failed to registry the metric", zap.Error(err))
		}
	}
}

func initDelayHistogram(ctx context.Context) {
	logger := utils.GetLogger(ctx)
	if delayHistogram == nil {
		opts := prometheus.HistogramOpts{Name: delayHistogramName,
			Help: delayHistogramDesc, Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}
		delayHistogram = prometheus.NewHistogram(opts)
		if err := registry.Register(delayHistogram); err != nil {
			logger.Error("Failed to registry the metric", zap.Error(err))
		}
	}
}

func initDelaySummary(ctx context.Context) {
	logger := utils.GetLogger(ctx)
	if delaySummary == nil {
		opts := prometheus.SummaryOpts{Name: delaySummaryName,
			Help:       delaySummaryDesc,
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}}
		delaySummary = prometheus.NewSummary(opts)
		if err := registry.Register(delaySummary); err != nil {
			logger.Error("Failed to registry the metric", zap.Error(err))
		}
	}
}

func GetPrometheusHttpHandler(ctx context.Context) http.Handler {
	logger := utils.GetLogger(ctx)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	logger.Info("Created HTTP handler for prometheus metrics")
	return handler
}

func GetApiInvokeCountCounter(ctx context.Context) *prometheus.CounterVec {
	return apiInvokeCounter
}

func GetDelayHistogram(ctx context.Context) prometheus.Histogram {
	return delayHistogram
}

func GetDelaySummary(ctx context.Context) prometheus.Histogram {
	return delaySummary
}

func MonitorConnectionsMetric(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		connectionsGauge.Inc()
	case http.StateClosed, http.StateHijacked:
		connectionsGauge.Dec()
	}
}
