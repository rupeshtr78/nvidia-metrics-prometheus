package api

import (
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"go.uber.org/zap"
)

func RunMetrics() {
	// Initialize NVML before starting the metric collection loop
	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	http.Handle("/metrics", promhttp.Handler())

	// Start a separate goroutine for collecting metrics at regular intervals
	go func() {
		for {
			nvidiametrics.CollectGpuMetrics()
			time.Sleep(5 * time.Second)
		}
	}()

	// ListenAndServe on port 9500
	err := http.ListenAndServe(":9500", nil)
	if err != nil {
		log.Fatal("HTTP server failed", zap.Error(err))
	}
}

func RunPrometheusMetricsServer() {
	// Initialize NVML before starting the metric collection loop
	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	// Start collecting GPU metrics every 5 seconds
	startMetricsCollection(5 * time.Second)

	// Start the HTTP server to expose metrics
	StartPrometheusServer()
}

func startMetricsCollection(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				nvidiametrics.CollectGpuMetrics()
			}
		}
	}()
}

func StartPrometheusServer() {

	server := &http.Server{
		Addr:         ":9500",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      promhttp.Handler(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}
}
