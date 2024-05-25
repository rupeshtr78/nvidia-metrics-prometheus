package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rupeshtr78/nvidia-metrics/api"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("config", "config/metrics.yaml", "Path to the configuration file")
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error,fatal)")
	port 	 = flag.String("port", "9500", "Port to run the metrics server")
)

func main() {
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Config file is required")
	}

	filePath := filepath.Join(*configFile)

	// Register the metrics with Prometheus and start the metrics server
	// provide --config "config/metrics.yaml"
	err := prometheusmetrics.CreatePrometheusMetrics(filePath)
	if err != nil {
		logger.Fatal("Failed to create Prometheus metrics", zap.Error(err))
		os.Exit(1)
	}

	//// run the metrics server
	api.RunPrometheusMetricsServer()


}

// @TODO - Remove this function for testing only
func RunMetricsLocal() {

	// Initialize NVML before starting the metric collection loop
	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	for {
		nvidiametrics.CollectGpuMetrics()
		time.Sleep(30 * time.Second)
		for key, label := range prometheusmetrics.RegisteredLabels {
			logger.Debug("Registered label", zap.String("key", key), zap.Any("label", label))
		}
		// "key":"gpu_power_usage","label":{"label1":"gpu_id","label2":"gpu_name"}}

	}

	// get all labels for debugging

}
