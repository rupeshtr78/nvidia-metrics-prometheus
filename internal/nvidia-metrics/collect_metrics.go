package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetrics() {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device count", zap.Error(err))
		return
	}

	for i := 0; i < int(deviceCount); i++ {
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
		// Replace this with actual usage.
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
		"Collected device info",
		zap.Int("device_index", deviceIndex),
		zap.String("device_name", deviceName),
		// zap.Int("pci_bus_id", int(pciInfo.)),
	)

	metrics := &GPUDeviceMetrics{
		DeviceIndex: deviceIndex,
	}

	labels := map[string]string{
		"gpu_id":   fmt.Sprintf("%d", deviceIndex),
		"gpu_name": deviceName,
	}

	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		metrics.GPUTemperature = float64(temperature)
		gauge.SetGaugeMetric("gpu_temperature", labels, metrics.GPUTemperature)
	}

	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		metrics.GPUCPUUtilization = float64(utilization.Gpu)
		metrics.GPUMemUtilization = float64(utilization.Memory)
		gauge.SetGaugeMetric("gpu_cpu_utilization", labels, metrics.GPUCPUUtilization)
		gauge.SetGaugeMetric("gpu_mem_utilization", labels, metrics.GPUMemUtilization)
	}

	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		metrics.GPUPowerUsage = float64(gpuPowerUsage) / 1000 // Assuming power is in mW and we want W.
		gauge.SetGaugeMetric("gpu_power_usage", labels, metrics.GPUPowerUsage)
	}

	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		metrics.GPURunningProcesses = len(runningProcess)
		gauge.SetGaugeMetric("gpu_running_process", labels, float64(metrics.GPURunningProcesses))
	}

	// Add more metrics here.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
	return metrics, nil
}

// func GetNvidiaLabels(handle nvml.Device) {

// handle.GetPciInfo()
// handle.GetTemperatureThreshold(nvml.TEMPERATURE_THRESHOLD_SLOWDOWN)
// 	handle.GetTemperatureThreshold(nvml.TEMPERATURE_THRESHOLD_SHUTDOWN)
// 	handle.GetClockInfo(nvml.CLOCK_GRAPHICS)
// 	handle.GetClockInfo(nvml.CLOCK_SM)
// 	handle.GetClockInfo(nvml.CLOCK_MEM)
// 	handle.GetClockInfo(nvml.CLOCK_VIDEO)
// 	handle.GetClockInfo(nvml.CLOCK_COUNT)
// 	handle.GetSupportedVgpus()
// 	handle.GetVgpuMetadata()
// 	handle.GetVgpuUtilization()
// 	handle.GetVgpuProcessUtilization()
// 	handle.GetAttributes()
// 	handle.GetEncoderCapacity()
// 	handle.GetEncoderStats()
// 	handle.GetActiveVgpus()
// 	handle.GetComputeMode()
// 	handle.GetEccMode()
// 	handle.GetTotalEccErrors()
// 	handle.GetDetailedEccErrors()
// 	handle.GetMemoryErrorCounter()
// 	handle.GetVbiosVersion()
// 	handle.GetBridgeChipInfo()
// 	handle.GetInforomVersion(nvml.INFOROM_OEM)
// 	handle.GetInforomVersion(nvml.INFOROM_ECC)
// 	handle.GetInforomVersion(nvml.INFOROM_POWER)
// 	handle.GetInforomVersion(nvml.INFOROM_COUNT)
// 	handle.GetInforomImageVersion()
// 	handle.GetMaxClockInfo(nvml.CLOCK_GRAPHICS)
// 	handle.GetMaxClockInfo(nvml.CLOCK_SM)
// 	handle.GetMaxClockInfo(nvml.CLOCK_MEM)
// 	handle.GetMaxClockInfo(nvml.CLOCK_VIDEO)
// 	handle.GetMaxClockInfo(nvml.CLOCK_COUNT)
// 	handle.GetApplicationsClock(nvml.CLOCK_GRAPHICS)
// 	handle.GetApplicationsClock(nvml.CLOCK_SM)
// 	handle.GetApplicationsClock(nvml.CLOCK_MEM)
// 	handle.GetApplicationsClock(nvml.CLOCK_VIDEO)
// 	handle.GetApplicationsClock(nvml.CLOCK_COUNT)
// 	handle.GetDefaultApplicationsClock(nvml.CLOCK_GRAPHICS)
// 	handle.GetDefaultApplicationsClock(nvml.CLOCK_SM)
// 	handle.GetDefaultApplicationsClock(nvml.CLOCK_MEM)
// 	handle.GetDefaultApplicationsClock(nvml.CLOCK_VIDEO)
// 	handle.GetDefaultApplicationsClock(nvml.CLOCK_COUNT)
// 	handle.GetCudaComputeCapability()
// 	handle.GetNvLinkCapability()
// 	handle.GetNvLinkRemotePciInfo()
// }
