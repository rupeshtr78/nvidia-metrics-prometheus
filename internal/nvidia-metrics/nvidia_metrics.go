package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
)

func CollectUtilizationMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metrics.GPUCPUUtilization = float64(utilization.Gpu)
		SetDeviceMetric(handle, config.GPU_CPU_UTILIZATION, metrics.GPUCPUUtilization)

		metrics.GPUMemUtilization = float64(utilization.Memory)
		SetDeviceMetric(handle, config.GPU_MEM_UTILIZATION, metrics.GPUMemUtilization)
	}

	return err
}

func CollectMemoryInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
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

func CollectPowerInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		SetDeviceMetric(handle, metric, metrics.GPUPowerUsage)
	}

	return err
}

func CollectRunningProcessMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		SetDeviceMetric(handle, metric, float64(metrics.GPURunningProcesses))
	}

	return err
}

func CollectTemperatureMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		metrics.GPUTemperature = float64(temperature)
		SetDeviceMetric(handle, metric, metrics.GPUTemperature)
	}

	return err
}

func CollectDeviceIdAsMetric(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	deviceId, err := handle.GetIndex()
	if err == nvml.SUCCESS {
		metrics.DeviceIndex = deviceId
		SetDeviceMetric(handle, metric, float64(metrics.DeviceIndex))
	}

	return err
}

// @TODO added as label remove if not needed
func CollectTempThresholdShutdownAsMetric(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	metricValue, err := handle.GetTemperatureThreshold(nvml.TEMPERATURE_THRESHOLD_SHUTDOWN)
	if err == nvml.SUCCESS {
		SetDeviceMetric(handle, metric, float64(metricValue))
	}
	return err
}

func CollectClock(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	metricValue, err := handle.GetClock(nvml.CLOCK_MEM, nvml.CLOCK_ID_CURRENT)
	if err == nvml.SUCCESS {
		SetDeviceMetric(handle, metric, float64(metricValue))
	}

	return err
}
