package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
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
	GpuPState           int32
	GpuClock            uint32
	GpuEccErrors        uint64
	GpuFanSpeed         uint32
	GpuPeakFlops        float64
}

func isRegistered(metric config.Metric) bool {
	if _, ok := prometheusmetrics.RegisteredMetrics[metric.GetMetric()]; !ok {
		logger.Warn("metric not registered", zap.String("metric", metric.GetMetric()))
		return false
	}
	return true
}

// CollectDeviceMetrics collects all the metrics for the GPU device.
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

// CollectMemoryInfoMetrics collects the memory usage metrics for the GPU device.
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

// CollectPowerInfoMetrics collects the power usage metrics for the GPU device.
func CollectPowerInfoMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		SetDeviceMetric(handle, metric, metrics.GPUPowerUsage)
	}

	return err
}

// CollectRunningProcessMetrics collects the number of running processes on the GPU device.
func CollectRunningProcessMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		SetDeviceMetric(handle, metric, float64(metrics.GPURunningProcesses))
	}

	return err
}

// CollectTemperatureMetrics collects the temperature metrics for the GPU device.
func CollectTemperatureMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		metrics.GPUTemperature = float64(temperature)
		SetDeviceMetric(handle, metric, metrics.GPUTemperature)
	}

	return err
}

// CollectDeviceIdAsMetric collects the device id as a metric.
func CollectDeviceIdAsMetric(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	deviceId, err := handle.GetIndex()
	if err == nvml.SUCCESS {
		metrics.DeviceIndex = deviceId
		SetDeviceMetric(handle, metric, float64(metrics.DeviceIndex))
	}

	return err
}


func collectPStateMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	pState, err := handle.GetPerformanceState()
	if err == nvml.SUCCESS {
		metrics.GpuPState = int32(pState)
		SetDeviceMetric(handle, metric, float64(metrics.GpuPState))
	}

	return err

}

// collectMemoryClockMetrics collects the memory clock metrics for the GPU device.
func collectMemoryClockMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	memoryClock, err := handle.GetClock(nvml.CLOCK_MEM, nvml.CLOCK_ID_CURRENT)
	if err == nvml.SUCCESS {
		metrics.GpuClock = memoryClock
		SetDeviceMetric(handle, metric, float64(memoryClock))
	}

	return err
}

// collectMemoryClockMetrics collects the memory clock metrics for the GPU device.
func collectGpuClockMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	memoryClock, err := handle.GetClock(nvml.CLOCK_SM, nvml.CLOCK_ID_CURRENT)
	if err == nvml.SUCCESS {
		metrics.GpuClock = memoryClock
		SetDeviceMetric(handle, metric, float64(memoryClock))
	}

	return err
}

func collectGpuVideoClockMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	memoryClock, err := handle.GetClock(nvml.CLOCK_VIDEO, nvml.CLOCK_ID_CURRENT)
	if err == nvml.SUCCESS {
		metrics.GpuClock = memoryClock
		SetDeviceMetric(handle, metric, float64(memoryClock))
	}

	return err
}

func collectGpuGraphicsClockMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	memoryClock, err := handle.GetClock(nvml.CLOCK_GRAPHICS, nvml.CLOCK_ID_CURRENT)
	if err == nvml.SUCCESS {
		metrics.GpuClock = memoryClock
		SetDeviceMetric(handle, metric, float64(memoryClock))
	}

	return err
}

func collectEccCorrectedErrorsMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {

	eccErrors, err := handle.GetTotalEccErrors(nvml.MEMORY_ERROR_TYPE_CORRECTED, nvml.VOLATILE_ECC)
	if err == nvml.SUCCESS {
		metrics.GpuEccErrors = eccErrors
		SetDeviceMetric(handle, metric, float64(eccErrors))
	}
	return err
}

func collectEccUncorrectedErrorsMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	eccErrors, err := handle.GetTotalEccErrors(nvml.MEMORY_ERROR_TYPE_UNCORRECTED, nvml.VOLATILE_ECC)
	if err == nvml.SUCCESS {
		metrics.GpuEccErrors = eccErrors
		SetDeviceMetric(handle, metric, float64(eccErrors))
	}
	return err
}

// @TODO Fix this not collecting the right value
func collectFanSpeedMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	fans, err := handle.GetNumFans()
	if err != nvml.SUCCESS || fans == 0 {
		return err
	}
	fanSpeed, err := handle.GetFanSpeed_v2(fans)
	if err == nvml.SUCCESS {
		metrics.GpuFanSpeed = fanSpeed
		SetDeviceMetric(handle, metric, float64(fanSpeed))
	}
	return err
}

func collectPeakFlopsMetrics(handle nvml.Device, metrics *GPUDeviceMetrics, metric config.Metric) nvml.Return {
	// Retrieve max clock speed
	maxClock, err := handle.GetMaxClockInfo(nvml.CLOCK_GRAPHICS)
	if err != nvml.SUCCESS || maxClock == 0 {
		return err
	}

	currentClock, err := handle.GetClockInfo(nvml.CLOCK_GRAPHICS)
	if err != nvml.SUCCESS || currentClock == 0 {
		return err
	}

	deviceId, err := handle.GetIndex()
	if err != nvml.SUCCESS {
		return err
	}

	// get device flops
	peakFlops, err := config.GetGpuFlops(deviceId)
	if err != nvml.SUCCESS || peakFlops == 0 {
		return err
	}

	// Calculate effective FLOPS
	effectiveFLOPS := float64(currentClock) / float64(maxClock) * peakFlops
	// 1e15 is 1 PetaFLOP
	// 1e12 is 1 TeraFLOP

	// calculate in TFLOPS convert to PFLOPS in grafana
	pflops := effectiveFLOPS / 1e12

	metrics.GpuPeakFlops = pflops
	SetDeviceMetric(handle, metric, pflops)

	return err
}

// @TODO add metrics here
func collectComputeRunningProcesses(handle nvml.Device, metrics *GPUDeviceMetrics) nvml.Return {

	//handle.GetTotalEccErrors(nvml.MEMORY_ERROR_TYPE_CORRECTED, nvml.ECC_COUNTER_TYPE_COUNT)
	processInfos, _ := handle.GetComputeRunningProcesses()
	for _, processInfo := range processInfos {
		fmt.Printf("Process Info: %+v\n", processInfo)
		fmt.Printf("Process Info: %+v\n", processInfo.Pid)
		fmt.Printf("Process Info: %+v\n", processInfo.UsedGpuMemory)
		fmt.Printf("Process Info: %+v\n", processInfo.GpuInstanceId)
		fmt.Printf("Process Info: %+v\n", processInfo.ComputeInstanceId)
		_ = metrics
	}
	return nvml.SUCCESS

}

// Additional Metrics can be added here
//handle.GetActiveVgpus()
//handle.GetEncoderUtilization()
//handle.GetDecoderUtilization()
//handle.GetEccMode()
//handle.GetTotalEccErrors()
