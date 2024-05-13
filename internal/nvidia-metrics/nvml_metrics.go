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

		gpuId.WithLabelValues(device.UUID).Set(float64(device.UUID))
	}

}
