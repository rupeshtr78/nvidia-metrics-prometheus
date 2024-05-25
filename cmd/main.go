package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rupeshtr78/nvidia-metrics/api"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("config", "config/metrics.yaml", "Path to the configuration file")
	logLevel   = flag.String("loglevel", "info", "Log level (debug, info, warn, error,fatal)")
	port       = flag.String("port", "9500", "Port to run the metrics server")
	host       = flag.String("host", "localhost", "Host to run the metrics server")
	interval   = flag.String("interval", "5", "Time interval in seconds to scrape metrics")
)

func main() {
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Config file is required")
	}

	// Initialize the logger
	err := logger.GetLogger(*logLevel)
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}

	filePath := filepath.Join(*configFile)

	// Register the metrics with Prometheus
	err = prometheusmetrics.CreatePrometheusMetrics(filePath)
	if err != nil {
		logger.Fatal("Failed to create Prometheus metrics", zap.Error(err))
		os.Exit(1)
	}

	// get the address from the flag host and port
	address := *host + ":" + *port

	//  get the metrics scrape interval
	t, err := strconv.Atoi(*interval)
	if err != nil {
		logger.Fatal("Failed to convert time to integer", zap.Error(err))
		os.Exit(1)
	}

	var scrapreInterval time.Duration
	if t > 0 {
		scrapreInterval = time.Duration(t) * time.Second
	}

	// start the metrics server
	api.RunPrometheusMetricsServer(address, scrapreInterval)

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
