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

	err = CollectTemperatureMetrics(handle, metrics, config.GPU_TEMPERATURE)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting temperature metrics", zap.Error(err))
	}

	err = CollectUtilizationMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting utilization metrics", zap.Error(err))
	}

	err = CollectMemoryInfoMetrics(handle, metrics)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting memory info metrics", zap.Error(err))
	}

	err = CollectPowerInfoMetrics(handle, metrics, config.GPU_POWER_USAGE)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting power info metrics", zap.Error(err))
	}

	err = CollectRunningProcessMetrics(handle, metrics, config.GPU_RUNNING_PROCESS)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting running process metrics", zap.Error(err))
	}

	err = CollectDeviceIdAsMetric(handle, metrics, config.GPU_ID_METRIC)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting device id as metric", zap.Error(err))
	}

	err = collectPStateMetrics(handle, metrics, config.GPU_P_STATE)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting p state metrics", zap.Error(err))
	}

	err = collectEccCorrectedErrorsMetrics(handle, metrics, config.GPU_ECC_CORRECTED_ERRORS)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting ECC corrected errors metrics", zap.Error(err))
	}

	err = collectEccUncorrectedErrorsMetrics(handle, metrics, config.GPU_ECC_UNCORRECTED_ERRORS)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting ECC uncorrected errors metrics", zap.Error(err))
	}

	err = collectGpuClockMetrics(handle, metrics, config.GPU_SM_CLOCK)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting GPU clock metrics", zap.Error(err))
	}

	err = collectGpuGraphicsClockMetrics(handle, metrics, config.GPU_GRAPHICS_CLOCK)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting GPU graphics clock metrics", zap.Error(err))
	}

	err = collectGpuVideoClockMetrics(handle, metrics, config.GPU_VIDEO_CLOCK)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting GPU video clock metrics", zap.Error(err))
	}

	err = collectMemoryClockMetrics(handle, metrics, config.GPU_MEMORY_CLOCK)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting memory clock metrics", zap.Error(err))
	}

	err = collectFanSpeedMetrics(handle, metrics, config.GPU_FAN_SPEED)
	if err != nvml.SUCCESS {
		logger.Error("Error collecting fan speed metrics", zap.Error(err))
	}

	// @TODO Add more metrics here.

	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}

func updateUtilizationMetrics(handle nvml.Device, utilization nvml.Utilization) {
	updateMetric(config.GPU_CPU_UTILIZATION, handle, float64(utilization.Gpu))
	updateMetric(config.GPU_MEM_UTILIZATION, handle, float64(utilization.Memory))
}

func updateMetric(m config.Metric, handle nvml.Device, value float64) {
	metric := m.GetMetric()
	labels := labelManager.GetMetricLabelValues(handle, metric)
	gauge.SetGaugeMetric(metric, labels, value)
}
