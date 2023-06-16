package main

import (
	"fmt"

	"github.com/brianvoe/gofakeit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

const (
	pushgatewayIp      = "kcd-pushgateway"
	pushgatewayPort    = "9091"
	customersToProcess = 5
	delayToDelete      = 30
)

// Each staff (staffName) will process items for "customersToProcess" customers
func main() {
	pushgatewayUrl := fmt.Sprintf("http://%s:%s", pushgatewayIp, pushgatewayPort)
	logger, _ := zap.NewProduction()
	logger.Info("Starting the job with", zap.String("Pushgateway", pushgatewayUrl))
	staffName := gofakeit.FirstName()

	// in Pushgateway job key will have kcd-job
	// in Pushgateway instance key might be empty

	// in Prometheus exported_job key will have kcd-job
	// in Prometheus job key value will be based on the scrape config name
	// in Prometheus instace key will be pushgateway server url
	pushgateway := push.New(pushgatewayUrl, "kcd-job")

	opts := prometheus.GaugeOpts{
		Name:        "job_delay",
		Help:        "Time took for the Job to complete",
		ConstLabels: prometheus.Labels{"staff_name": staffName},
	}

	gauge := prometheus.NewGaugeVec(opts, []string{"customer_name"})

	for i := 0; i < customersToProcess; i++ {
		process(gauge)
	}
	// Only one metric so not using an registry
	pushgateway.Collector(gauge)
	if err := pushgateway.Push(); err != nil {
		logger.Error("Pushing the metrics to pushgateway failed", zap.Error(err))
		return
	}
	logger.Info("Pushed metrics to push gateway", zap.String("staff_name", staffName))

	// Wait and delete the metrics from the Pushgateway
	// time.Sleep(time.Second * delayToDelete)
	// if err := pushgateway.Delete(); err != nil {
	// 	logger.Error("Failed to delete the push gateway metrics", zap.Error(err))
	// 	return
	// }

	// logger.Info("Successfully delete the push gateway metrics", zap.String("staff_name", staffName))
}

func process(gauge *prometheus.GaugeVec) {
	timeConsumedInSec := gofakeit.Number(1, 10)
	jobName := gofakeit.FirstName()
	gauge.WithLabelValues(jobName).Set(float64(timeConsumedInSec))
}
