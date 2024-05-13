package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/pkg"
	"go.uber.org/zap"
)

var logger, _ = pkg.Logger()

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

		logger.Debug("GPU metrics", zap.String("name", name), zap.Uint("temperature", uint(temperature)))

		// Set the prometheus metrics for the GPU
		gpuId.WithLabelValues(fmt.Sprintf("gpu%d", i)).Set(float64(i))
		gpuName.WithLabelValues(fmt.Sprintf("gpu%d", i)).Set(deviceName)
		gpuTemperature.WithLabelValues(fmt.Sprintf("gpu%d", i)).Set(float64(temperature))
		gpuUtilization.WithLabelValues(fmt.Sprintf("gpu%d", i)).Set(float64(utilization.Gpu))
		gpuMemoryUtilization.WithLabelValues(fmt.Sprintf("gpu%d", i)).Set(float64(utilization.Memory))
	}
}
