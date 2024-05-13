package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// Assumed correct logger initialization elsewhere
var log = logger.GetLogger()

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

// CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetrics() {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device count", zap.Error(err))
		return
	}

	for i := 0; i < int(deviceCount); i++ {
		collectDeviceMetrics(i)
	}
}

// collectDeviceMetrics collects metrics for a single device.
func collectDeviceMetrics(deviceIndex int) {
	handle, err := nvml.DeviceGetHandleByIndex(deviceIndex)
	if err != nvml.SUCCESS {
		logger.Error("Error getting device handle", zap.Int("device_index", deviceIndex), zap.Error(err))
		return
	}

	deviceName, err := handle.GetName()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device name", zap.Error(err))
		return
	}

	logger.Debug("Collecting GPU metrics", zap.Int("device_index", deviceIndex), zap.String("device_name", deviceName))

	labels := map[string]string{
		"gpu_id":   fmt.Sprintf("%d", deviceIndex),
		"gpu_name": deviceName,
	}

	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		setGaugeMetric("gpu_temperature", labels, float64(temperature))
	}

	logger.Debug("Collected GPU temperature", zap.Int("device_index", deviceIndex), zap.Float64("temperature", float64(temperature)))

	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		setGaugeMetric("gpu_cpu_utilization", labels, float64(utilization.Gpu))
		setGaugeMetric("gpu_mem_utilization", labels, float64(utilization.Memory))
	}

	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		setGaugeMetric("gpu_power_usage", labels, float64(gpuPowerUsage)/1000) // Assuming power is in mW and we want W.
	}

	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		setGaugeMetric("gpu_running_process", labels, float64(len(runningProcess)))
	}

	// Add more metrics here as needed.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
}

// setGaugeMetric sets a gauge metric with the given name, labels, and value.
func setGaugeMetric(name string, labels map[string]string, value float64) {
	err := prometheusmetrics.CreateGauge(name, labels, value)
	if err != nil {
		logger.Error("Failed to create gauge metric", zap.String("metric_name", name), zap.Error(err))
	}
}
