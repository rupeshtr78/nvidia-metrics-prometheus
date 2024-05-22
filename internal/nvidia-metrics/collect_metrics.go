package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
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
		// @TODO add this to slice of metrics for cli client
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
		"Collecting metrics for device",
		zap.Int("device_index", deviceIndex),
		zap.String("device_name", deviceName),
	)

	metrics := &GPUDeviceMetrics{
		DeviceIndex: deviceIndex,
	}

	// Collect Device Metrics

	err = collectTemperatureMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting temperature metrics", zap.Error(err))
	}

	err = collectUtilizationMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting utilization metrics", zap.Error(err))
	}

	err = collectMemoryInfoMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting memory info metrics", zap.Error(err))
	}

	err = collectPowerInfoMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting power info metrics", zap.Error(err))
	}

	err = collectRunningProcessMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting running process metrics", zap.Error(err))
	}

	err = collectDeviceIdAsMetric(handle, metrics, config.GPU_ID_METRIC)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting device id as metric", zap.Error(err))
	}

	// Add more metrics here.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}

func collectUtilizationMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metrics.GPUCPUUtilization = float64(utilization.Gpu)
		SetDeviceMetric(handle, config.GPU_CPU_UTILIZATION, metrics.GPUCPUUtilization)

		metrics.GPUMemUtilization = float64(utilization.Memory)
		SetDeviceMetric(handle, config.GPU_MEM_UTILIZATION, metrics.GPUMemUtilization)
	}

	return err
}

func collectMemoryInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	memoryInfo, err := handle.GetMemoryInfo()
	if err == nvml.SUCCESS {
		// Memory usage is in bytes, converting to GB.
		metrics.GPUMemoryUsed = uint64(memoryInfo.Used) / 1024 / 1024 / 1024
		metrics.GPUMemoryTotal = uint64(memoryInfo.Total) / 1024 / 1024 / 1024
		metrics.GPUMemoryFree = uint64(memoryInfo.Free) / 1024 / 1024 / 1024

		SetDeviceMetric(handle, config.GPU_MEMORY_USED, float64(metrics.GPUMemoryUsed))
		SetDeviceMetric(handle, config.GPU_MEMORY_TOTAL, float64(metrics.GPUMemoryTotal))
		SetDeviceMetric(handle, config.GPU_MEMORY_FREE, float64(metrics.GPUMemoryFree))
	}

	return err
}

func collectPowerInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		SetDeviceMetric(handle, config.GPU_POWER_USAGE, metrics.GPUPowerUsage)
	}

	return err
}

func collectRunningProcessMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		SetDeviceMetric(handle, config.GPU_RUNNING_PROCESS, float64(metrics.GPURunningProcesses))
	}

	return err
}

func collectTemperatureMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		metrics.GPUTemperature = float64(temperature)
		SetDeviceMetric(handle, config.GPU_TEMPERATURE, metrics.GPUTemperature)
	}

	return err
}

func collectDeviceIdAsMetric(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	deviceId, err := handle.GetIndex()
	if err == nvml.SUCCESS {
		metrics.DeviceIndex = deviceId
		SetDeviceMetric(handle, metric, float64(metrics.DeviceIndex))
	}

	return err
}
