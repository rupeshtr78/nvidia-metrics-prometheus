package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
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
}

// CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetrics() {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device count", zap.Error(err))
		return
	}

	for i := 0; i < int(deviceCount); i++ {
		metrics, err := collectDeviceMetrics(i)
		if err != nil {
			logger.Error(
				"Error collecting metrics for GPU",
				zap.Int("gpu_index", i),
				zap.Error(err),
			)
			continue // Skip this GPU and proceed with the next one
		}
		// Use the collected metrics if needed
		// Replace this with actual usage.
		_ = metrics
	}

	// Here we have successfully collected metrics for all GPUs without errors.
	logger.Info("Successfully collected metrics for all GPUs")
}

// CollectGpuDeviceMetrics collects metrics for a single device and returns them in a GPUDeviceMetrics struct.
func collectDeviceMetrics(deviceIndex int) (*GPUDeviceMetrics, error) {
	handle, err := nvml.DeviceGetHandleByIndex(deviceIndex)
	if err != nvml.SUCCESS {
		logger.Error("Error getting device handle", zap.Int("device_index", deviceIndex), zap.Error(err))
		return nil, err
	}

	deviceName, err := handle.GetName()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device name", zap.Error(err))
		return nil, err
	}

	metrics := &GPUDeviceMetrics{
		DeviceIndex: deviceIndex,
	}

	labels := map[string]string{
		"gpu_id":   fmt.Sprintf("%d", deviceIndex),
		"gpu_name": deviceName,
	}

	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		metrics.GPUTemperature = float64(temperature)
		gauge.SetGaugeMetric("gpu_temperature", labels, metrics.GPUTemperature)
	}

	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metrics.GPUCPUUtilization = float64(utilization.Gpu)
		metrics.GPUMemUtilization = float64(utilization.Memory)
		gauge.SetGaugeMetric("gpu_cpu_utilization", labels, metrics.GPUCPUUtilization)
		gauge.SetGaugeMetric("gpu_mem_utilization", labels, metrics.GPUMemUtilization)
	}

	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		gauge.SetGaugeMetric("gpu_power_usage", labels, metrics.GPUPowerUsage)
	}

	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		gauge.SetGaugeMetric("gpu_running_process", labels, float64(metrics.GPURunningProcesses))
	}

	// Add more metrics here.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}
