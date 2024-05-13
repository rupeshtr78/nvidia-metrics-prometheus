package main

import (
	"fmt"
	"os"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var err = logger.GetLogger()

func main() {
	fmt.Println("Hello, nvidia-metrics")
	logger.Info("Starting nvidia-metrics")

	// Register the metrics with Prometheus
	err := prometheusmetrics.CreatePrometheusMetrics("config/metrics.yaml")
	if err != nil {
		logger.Fatal("Failed to create Prometheus metrics", zap.Error(err))
		os.Exit(1)
	}

	// run the metrics server
	// api.RunMetrics()
	RunMetricsLocal()

}

func RunMetricsLocal() {

	nvidiametrics.CollectGpuMetrics()
}
