package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
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
	GpuPState           int32
}

func CollectUtilizationMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metrics.GPUCPUUtilization = float64(utilization.Gpu)
		SetDeviceMetric(handle, config.GPU_GPU_UTILIZATION, metrics.GPUCPUUtilization)

		metrics.GPUMemUtilization = float64(utilization.Memory)
		SetDeviceMetric(handle, config.GPU_MEM_UTILIZATION, metrics.GPUMemUtilization)
	}

	return err
}

func CollectMemoryInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {
	memoryInfo, err := handle.GetMemoryInfo()
	if err == nvml.SUCCESS {
		// Memory usage is in bytes, converting to GB.
		metrics.GPUMemoryUsed = uint64(memoryInfo.Used) / 1024 / 1024 //  memory is in bytes and we want MB
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

// P0/P1 - Maximum 3D performance
// P2/P3 - Balanced 3D performance-power
// P8 - Basic HD video playback
// P10 - DVD playback
// P12 - Minimum idle power consumption
// PState is the current performance state of the GPU device.
func collectPStateMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	pState, err := handle.GetPerformanceState()
	if err == nvml.SUCCESS {
		metrics.GpuPState = int32(pState)
		SetDeviceMetric(handle, metric, float64(metrics.GpuPState))
	}

	return err

}
