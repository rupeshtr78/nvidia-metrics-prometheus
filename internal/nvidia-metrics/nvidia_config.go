package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
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
}

// InitNVML initializes the NVML library.
func InitNVML() {
	if err := nvml.Init(); err != nvml.SUCCESS {
		logger.Fatal("Failed to initialize NVML", zap.Error(err))
	}
	logger.Info("Initialized NVML")
}

// ShutdownNVML shuts down the NVML library.
func ShutdownNVML() {
	if err := nvml.Shutdown(); err != nvml.SUCCESS {
		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
	}
	logger.Info("Shutdown NVML")
}

// GetLabelKeys returns the label keys for the given metric name.
// Example: "key":"gpu_power_usage","label":{"label1":"gpu_id","label2":"gpu_name"}}
func GetLabelKeys(metricName string) map[string]string {
	labelKeys := make(map[string]string)

	if _, ok := prometheusmetrics.RegisteredLabels[metricName]; !ok {
		logger.Warn("Metric not found in registered labels", zap.String("metric_name", metricName))
		return labelKeys
	}

	keys := prometheusmetrics.RegisteredLabels[metricName]
	// iterate over the keys and add label name to the map
	for i, key := range keys {
		if len(key) == 0 {
			continue
		}
		// adding dummy values for now will be updates while setting guage
		labelKeys[key] = i
	}

	return labelKeys

}

type LabelFunction struct {
	LabelFunc map[string]func(nvml.Device) (any, nvml.Return)
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
// handle.GetSupportedVgpus()
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
