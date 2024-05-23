package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
)

// SetDeviceMetric sets the metric value for the given device
func SetDeviceMetric(handle nvml.Device, metricConfig config.Metric, metricValue float64) {
	metric := metricConfig.GetMetric()
	metricLabels := labelManager.GetMetricLabelValues(handle, metric)
	gauge.SetGaugeMetric(metric, metricLabels, metricValue)
}

// AddFunctions adds the label function to the map
func (lf LabelFunctions) AddFunctions() {

	lf.Add(config.GPU_ID.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		index, ret := device.GetIndex()
		return index, ret
	})

	lf.Add(config.GPU_NAME.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		name, ret := device.GetName()
		return name, ret
	})

	// GPU temperature threshold protections can shut down system when it hits the temp.limit,
	lf.Add(config.GPU_TEM_THRESHOLD.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		threshold, ret := device.GetTemperatureThreshold(nvml.TEMPERATURE_THRESHOLD_SHUTDOWN)
		return threshold, ret
	})

	//determines the rate at which the GPU can access and manipulate data stored in the VRAM
	lf.Add(config.GPU_MEM_CLOCK.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		clock, ret := device.GetClock(nvml.CLOCK_MEM, nvml.CLOCK_ID_CURRENT)
		return clock, ret

	})

	lf.Add(config.GPU_CORES.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		cores, ret := device.GetNumGpuCores()
		return cores, ret
	})

	// @TODO add additional label function to the map
	// lf.Add(config.GPU_POWER.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
	// 	nvml.SystemGetCudaDriverVersion()
	// 	pci, ret := nvml.DeviceGetPciInfo(device)
	//     nvml.DeviceGetSupportedMemoryClocks(device)
	// })
}
