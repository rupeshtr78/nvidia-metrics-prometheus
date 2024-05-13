package main

import (
	"fmt"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
)

func main() {
	fmt.Println("Hello, nvidia-metrics")
	nvidiametrics.CollectGpuMetrics()
}
