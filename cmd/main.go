package main

import (
	"fmt"
	"time"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics/api"
	"github.com/rupeshtr78/nvidia-metrics/pkg"
)

var logger, _ = pkg.Logger()

func main() {
	fmt.Println("Hello, nvidia-metrics")
	logger.Info("Starting nvidia-metrics")

	api.RunMetrics()

}

func RunMetricsLocal() {

	for {
		nvidiametrics.CollectGpuMetrics()
		time.Sleep(30 * time.Second)
	}
}
