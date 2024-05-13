package api

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
)

func RunMetrics() {

	nvidiametrics.InitNVML()
	defer nvidiametrics.ShutdownNVML()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		for {
			nvidiametrics.CollectGpuMetrics()
			time.Sleep(10 * time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":9500", nil))
}
