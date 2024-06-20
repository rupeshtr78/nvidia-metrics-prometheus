package nvidiametrics

import (
	"context"
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var labelManager = NewLabelFunction()

// CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetrics(ctx context.Context) {
	deviceCount := CollectGPUDeviceCount(ctx)
	if deviceCount == 0 {
		logger.Error("No GPU devices found")
		return
	}

	// Add label functions
	labelManager.AddFunctions()

	for i := 0; i < deviceCount; i++ {
		metrics, err := collectDeviceMetrics(ctx, i)
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
func collectDeviceMetrics(ctx context.Context, deviceIndex int) (*GPUDeviceMetrics, error) {
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

	metrics := NewGPUDeviceMetrics()
	metrics.DeviceIndex = deviceIndex

	// Collect Device Metrics
	if isRegistered(config.GPU_TEMPERATURE) {
		err = metrics.CollectTemperatureMetrics(ctx, handle, config.GPU_TEMPERATURE)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting temperature metrics", zap.Error(err))
		}
	}

	err = metrics.CollectUtilizationMetrics(ctx, handle)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting utilization metrics", zap.Error(err))
	}

	err = metrics.CollectMemoryInfoMetrics(ctx, handle)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting memory info metrics", zap.Error(err))
	}

	if isRegistered(config.GPU_POWER_USAGE) {
		err = metrics.CollectPowerInfoMetrics(ctx, handle, config.GPU_POWER_USAGE)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting power info metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_RUNNING_PROCESS) {
		err = metrics.CollectRunningProcessMetrics(ctx, handle, config.GPU_RUNNING_PROCESS)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting running process metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_ID_METRIC) {
		err = metrics.CollectDeviceIdAsMetric(ctx, handle, config.GPU_ID_METRIC)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting device id as metric", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_P_STATE) {
		err = metrics.collectPStateMetrics(ctx, handle, config.GPU_P_STATE)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting p state metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_ECC_CORRECTED_ERRORS) {
		err = metrics.collectEccCorrectedErrorsMetrics(ctx, handle, config.GPU_ECC_CORRECTED_ERRORS)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting ECC corrected errors metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_ECC_UNCORRECTED_ERRORS) {
		err = metrics.collectEccUncorrectedErrorsMetrics(ctx, handle, config.GPU_ECC_UNCORRECTED_ERRORS)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting ECC uncorrected errors metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_SM_CLOCK) {
		err = metrics.collectGpuClockMetrics(ctx, handle, config.GPU_SM_CLOCK)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting GPU clock metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_GRAPHICS_CLOCK) {
		err = metrics.collectGpuGraphicsClockMetrics(ctx, handle, config.GPU_GRAPHICS_CLOCK)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting GPU graphics clock metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_VIDEO_CLOCK) {
		err = metrics.collectGpuVideoClockMetrics(ctx, handle, config.GPU_VIDEO_CLOCK)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting GPU video clock metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_MEMORY_CLOCK) {
		err = metrics.collectMemoryClockMetrics(ctx, handle, config.GPU_MEMORY_CLOCK)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting memory clock metrics", zap.Error(err))
		}
	}

	if isRegistered(config.GPU_PEAK_FLOPS_METRIC) {
		err = metrics.collectPeakFlopsMetrics(ctx, handle, config.GPU_PEAK_FLOPS_METRIC)
		if err != nvml.SUCCESS {
			logger.Error("Error collecting peak flops metrics", zap.Error(err))
		}
	}

	// @TODO Add more metrics here.

	// @TODO Fix fan speed metrics
	// err = collectFanSpeedMetrics(handle, metrics, config.GPU_FAN_SPEED)
	// if err != nvml.SUCCESS {
	// 	logger.Error("Error collecting fan speed metrics", zap.Error(err))
	// }

	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}

// CollectGPUDeviceCount collects the number of GPU devices.
func CollectGPUDeviceCount(ctx context.Context) int {
	var deviceCount int
	var err nvml.Return
	select {
	case <-ctx.Done():
		return 0
	default:
		deviceCount, err = nvml.DeviceGetCount()
		if err != nvml.SUCCESS {
			logger.Error("Error getting device count", zap.Error(fmt.Errorf("nvml error code: %d", err)))
			return 0
		}
	}

	return deviceCount
}
