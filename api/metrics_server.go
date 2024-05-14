package api

import (
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

	// ListenAndServe blocks the execution; if it returns an error, NVML will be shut down due to the defer statement above
	err := http.ListenAndServe(":9500", nil)
	if err != nil {
		log.Fatal("HTTP server failed", zap.Error(err))
	}
}
