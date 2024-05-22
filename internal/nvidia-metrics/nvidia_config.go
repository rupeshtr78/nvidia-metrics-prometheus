package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// GPUDeviceMetrics represents the collected metrics for a GPU device.
type GPUDeviceMetrics struct {
	DeviceIndex         int
	GPUTemperature      float64
	GPUCPUUtilization   float64
	GPUMemUtilization   float64
	GPUPowerUsage       float64
	GPURunningProcesses int
	GPUMemoryUsed       uint64
	GPUMemoryTotal      uint64
	GPUMemoryFree       uint64
}

// InitNVML initializes the NVML library.
func InitNVML() {
	if err := nvml.Init(); err != nvml.SUCCESS {
		logger.Fatal("Failed to initialize NVML", zap.Error(err))
	}
	logger.Info("Initialized NVML")
}

// ShutdownNVML shuts down the NVML library.
func ShutdownNVML() {
	if err := nvml.Shutdown(); err != nvml.SUCCESS {
		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
	}
	logger.Info("Shutdown NVML")
}

func SetDeviceMetric(handle nvml.Device, metricConfig config.Metric, metricValue float64) {
	metric := metricConfig.GetMetric()
	metricLabels := labelManager.GetMetricLabelValues(handle, metric)
	gauge.SetGaugeMetric(metric, metricLabels, metricValue)
}
