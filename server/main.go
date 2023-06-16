package main

import (
	"context"
	"net/http"

	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/metrics"
	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/suggestions"
	"github.com/DimpleRajaVamsi/kcd-prometheus-workshop/server/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	mainCtx := context.Background()
	mainlogger, _ := zap.NewProduction()
	startServer(mainCtx, mainlogger)
}

func startServer(ctx context.Context, logger *zap.Logger) {

	serverLogger := logger.Named("server")
	serverCtx := context.WithValue(ctx, utils.LoggerKey, serverLogger)

	mux := http.NewServeMux()

	// Init prometheus metrics
	metrics.IntMetrics(serverCtx)
	prometheusHandler := metrics.GetPrometheusHttpHandler(serverCtx)
	mux.Handle("/metrics", prometheusHandler)

	// Handle requests for beer suggestions
	beerSuggestion := suggestions.GetBeer(serverCtx)
	mux.HandleFunc("/beer", beerSuggestion.Suggest)

	// Handle requests for car suggestions
	carSuggestion := suggestions.GetCar(serverCtx)
	mux.HandleFunc("/car", carSuggestion.Suggest)

	// Handle requests for Delay suggestions
	delaySuggestions := suggestions.GetDelay(serverCtx)
	mux.HandleFunc("/delay", delaySuggestions.Suggest)

	// Handle requests for Default Prometheus metrics
	mux.HandleFunc("/default-metrics", promhttp.Handler().ServeHTTP)

	serverLogger.Info("Starting http server")
	server := http.Server{Addr: ":8089", Handler: mux, ConnState: metrics.MonitorConnectionsMetric}
	err := server.ListenAndServe()
	if err != nil {
		serverLogger.Error("Running Suggestions server failed", zap.Error(err))
	}
}
