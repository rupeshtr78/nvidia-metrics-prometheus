package main

import (
	"fmt"
	"time"

	"github.com/rupeshtr78/nvidia-metrics/api"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
)

var err = logger.GetLogger()

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
