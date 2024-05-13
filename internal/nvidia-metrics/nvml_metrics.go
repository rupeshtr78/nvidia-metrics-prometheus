package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var err = logger.GetLogger()

// GetGPUUtilization returns the GPU utilization of the device.
// The GPU utilization is the percent of time over the past sample period during which one or more kernels was executing on the GPU.

func InitNVML() {
	err := nvml.Init()
	if err != nvml.SUCCESS {
		logger.Fatal("Failed to initialize NVML", zap.Error(err))

	}

	logger.Info("Initialized NVML")
}

func ShutdownNVML() {
	err := nvml.Shutdown()
	if err != nvml.SUCCESS {
		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
	}

	logger.Info("Shutdown NVML")
}

func CollectGpuMetrics() {
	nvml.Init()
	defer nvml.Shutdown()

	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		fmt.Println("Error getting device count:", err)
		return
	}

	for i := 0; i < int(deviceCount); i++ {
		handle, err := nvml.DeviceGetHandleByIndex(i)
		if err != nvml.SUCCESS {
			fmt.Println("Error getting device handle:", err)
			continue
		}

		deviceName, err := handle.GetName()
		if err != nvml.SUCCESS {
			fmt.Println("Error getting device name:", err)
			continue
		}

		temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
		if err != nvml.SUCCESS {
			fmt.Println("Error getting temperature:", err)
			continue
		}

		utilization, err := handle.GetUtilizationRates()
		if err != nvml.SUCCESS {
			fmt.Println("Error getting utilization rates:", err)
			continue
		}

		gpuPowerUsage, err := handle.GetPowerUsage()
		if err != nvml.SUCCESS {
			fmt.Println("Error getting power usage:", err)
			continue
		}

		runningProcess, err := handle.GetComputeRunningProcesses()
		if err != nvml.SUCCESS {
			fmt.Println("Error getting running processes:", err)
			continue
		}

		utilizationRates, err := handle.GetUtilizationRates()
		if err != nvml.SUCCESS {
			fmt.Println("Error getting utilization rates:", err)
			continue
		}

		logger.Debug("GPU metrics", zap.String("name", deviceName), zap.Uint("temperature", uint(temperature)))

		// Set the prometheus metrics for the GPU
		gpuId.WithLabelValues(fmt.Sprintf("gpu_id: %d", i)).Set(float64(i))
		gpuName.WithLabelValues(fmt.Sprintf("gpu_name: %v", deviceName)).Set(float64(i))
		gpuTemperature.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(temperature))
		gpuUtilization.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(utilization.Gpu))
		gpuMemoryUtilization.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(utilization.Memory))
		gpuPowerUsageMetric.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(gpuPowerUsage) / 1000)
		gpuRunningProcess.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(len(runningProcess)))
		gpuUtilizationRates.WithLabelValues(fmt.Sprintf("gpu_id: %d, gpu_name: %v", i, deviceName)).Set(float64(utilizationRates.Gpu))
	}
}
