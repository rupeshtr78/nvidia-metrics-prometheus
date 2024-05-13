package nvidiametrics

import (
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/rupeshtr78/nvidia-metrics/pkg"
	"go.uber.org/zap"
)

var logger, _ = pkg.Logger()

// GetGPUUtilization returns the GPU utilization of the device.
// The GPU utilization is the percent of time over the past sample period during which one or more kernels was executing on the GPU.

func initNVML() {
	err := nvml.Init()
	if err != nil {
		logger.Fatal("Failed to initialize NVML", zap.Error(err))

	}
}

func shutdownNVML() {
	err := nvml.Shutdown()
	if err != nil {
		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
	}
}

func CollectGpuMetrics() {

	count, err := nvml.GetDeviceCount()
	if err != nil {
		logger.Fatal("Failed to get device count", zap.Error(err))
	}

	if count == 0 {
		logger.Fatal("No devices found")
	}

	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			logger.Fatal("Failed to get device", zap.Error(err))
		}

		deviceId, err := device.GetGPUInstanceId()
		if err != nil {
			logger.Error("Failed to get GPU instance ID", zap.Error(err))
		}

		deviceName := device.Model
		if err != nil {
			logger.Error("Failed to get device name", zap.Error(err))
		}

		gpuUtilization, err := device.GetUtilizationRates()
		if err != nil {
			logger.Error("Failed to get GPU utilization", zap.Error(err))
		}

		gpuTemperature, err := device.GetTemperature()
		if err != nil {
			logger.Error("Failed to get GPU temperature", zap.Error(err))
		}

		logger.Info("GPU metrics",
			zap.Uint("device_id", deviceId),
			zap.String("device_name", deviceName),
			zap.Uint("gpu_utilization", gpuUtilization.GPU),
			zap.Uint("memory_utilization", gpuUtilization.Memory),
			zap.Uint("temperature", gpuTemperature),
		)

	}

}
