package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rupeshtr78/nvidia-metrics/api"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Starting prometheus nvidia-metrics")

	// Deletes the metric first
	// prometheusmetrics.DeleteMetrics("config/metrics.yaml")

	// Register the metrics with Prometheus
	err := prometheusmetrics.CreatePrometheusMetrics("config/metrics.yaml")
	if err != nil {
		logger.Fatal("Failed to create Prometheus metrics", zap.Error(err))
		os.Exit(1)
	}

	// run the metrics server
	api.RunMetrics()
	// RunMetricsLocal()

}

func RunMetricsLocal() {

	// Initialize NVML before starting the metric collection loop
	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	for {
		nvidiametrics.CollectGpuMetrics()
		time.Sleep(30 * time.Second)
	}
}
