package nvidiametrics

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var labelManager = NewLabelFunction()

// CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetrics() {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device count", zap.Error(err))
		return
	}

	labelManager.AddFunctions()

	for i := 0; i < deviceCount; i++ {
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
		// To Replace this with actual usage.
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

	logger.Debug(
		"Collected device info",
		zap.Int("device_index", deviceIndex),
		zap.String("device_name", deviceName),
		// zap.Int("pci_bus_id", int(pciInfo.)),
	)

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

		gauge.SetGaugeMetric(config.GPU_TEMPERATURE.GetMetric(), labels, metrics.GPUTemperature)
	}

	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metricCpu := config.GPU_CPU_UTILIZATION.GetMetric()
		metricCpuLabels := labelManager.GetMetricLabelValues(handle, metricCpu)
		metrics.GPUCPUUtilization = float64(utilization.Gpu)

		metricMem := config.GPU_MEM_UTILIZATION.GetMetric()
		metricMemLabels := labelManager.GetMetricLabelValues(handle, metricMem)
		metrics.GPUMemUtilization = float64(utilization.Memory)
		gauge.SetGaugeMetric(metricCpu, metricCpuLabels, metrics.GPUCPUUtilization)
		gauge.SetGaugeMetric(config.GPU_MEM_UTILIZATION.GetMetric(), metricMemLabels, metrics.GPUMemUtilization)
	}

	//memoryInfo, err := handle.GetMemoryInfo()
	//if err == nvml.SUCCESS {
	//	// Memory usage is in bytes, converting to GB.
	//	metrics.GPUMemoryUsed = float64(memoryInfo.Used) / 1024 / 1024 / 1024
	//	metrics.GPUMemoryTotal = float64(memoryInfo.Total) / 1024 / 1024 / 1024
	//	metrics.GPUMemoryFree = float64(memoryInfo.Free) / 1024 / 1024 / 1024
	//	gauge.SetGaugeMetric(config.GPU_MEMORY_USED.GetMetric(), labels, metrics.GPUMemoryUsed)
	//	gauge.SetGaugeMetric(config.GPU_MEMORY_TOTAL.GetMetric(), labels, metrics.GPUMemoryTotal)
	//	gauge.SetGaugeMetric(config.GPU_MEMORY_FREE.GetMetric(), labels, metrics.GPUMemoryFree)
	//}

	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		gauge.SetGaugeMetric(config.GPU_POWER_USAGE.GetMetric(), labels, metrics.GPUPowerUsage)
	}

	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		gauge.SetGaugeMetric(config.GPU_RUNNING_PROCESS.GetMetric(), labels, float64(metrics.GPURunningProcesses))
	}

	// Add more metrics here.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}
