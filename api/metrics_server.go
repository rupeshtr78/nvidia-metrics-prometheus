package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	nvidiaMetrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "gpu_metrics_processed_ops_total",
		Help: "The total number of gpu metrics processed events",
	})
)

func RunPrometheusMetricsServer(ctx context.Context, address string, interval time.Duration) {
	// Initialize NVML before starting the metric collection loop
	nvidiaMetrics.InitNVML()
	defer nvidiaMetrics.ShutdownNVML()

	// Start collecting GPU metrics every 5 seconds
	startMetricsCollection(ctx, interval)

	// Start the HTTP server to expose metrics
	err := StartPrometheusServer(address)
	if err != nil {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}
}

func startMetricsCollection(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			nvidiaMetrics.CollectGpuMetrics(ctx)
			opsProcessed.Inc()
		}
	}()
}

func StartPrometheusServer(address string) error {
	server := &http.Server{
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      promhttp.Handler(),
	}

	logger.Info("Starting Prometheus server", zap.String("address", address), zap.String("path", "/metrics"))

	err := server.ListenAndServe()

	return err

}

// @TODO Remove after testing
// RunMetrics starts HTTP server on port 9500 to expose GPU metrics [Deprecated]
func RunMetrics() {
	// Initialize NVML before starting the metric collection loop
	nvidiaMetrics.InitNVML()
	defer nvidiaMetrics.ShutdownNVML()

	http.Handle("/metrics", promhttp.Handler())

	ctx := context.TODO()

	// Start a separate goroutine for collecting metrics at regular intervals
	go func() {
		for {
			nvidiaMetrics.CollectGpuMetrics(ctx)
			time.Sleep(5 * time.Second)
		}
	}()

	// ListenAndServe on port 9500
	err := http.ListenAndServe(":9500", nil)
	if err != nil {
		log.Fatal("HTTP server failed", zap.Error(err))
	}
}
