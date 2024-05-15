// Deprecated: This file is deprecated and will be removed in future.
package nvidiametrics

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	gauge "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// Deprecated CollectGpuMetrics collects metrics for all the GPUs.
func CollectGpuMetricsV0() {
	nvml.Init()
	defer nvml.Shutdown()
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		fmt.Println("Error getting device count:", err)
		return
	}
	// Loop through all the devices
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
		logger.Debug("GPU metrics", zap.String("name", deviceName), zap.Uint("temperature", uint(temperature)))
		// Set the prometheus metrics for the GPU maps to label values and sets the value
		// metrics.GpuId.WithLabelValues(fmt.Sprintf("%d", i)).Set(float64(i))
		// metrics.GpuName.WithLabelValues(fmt.Sprintf("%v", deviceName)).Set(float64(i))
		// metrics.GpuTemperature.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(temperature))
		// metrics.GpuUtilization.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(utilization.Gpu))
		// metrics.GpuMemoryUtilization.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(utilization.Memory))
		// metrics.GpuPowerUsageMetric.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(gpuPowerUsage) / 1000)
		// metrics.GpuRunningProcess.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(len(runningProcess)))
		prometheusmetrics.CreateGauge("gpu_id", map[string]string{"gpu_id": fmt.Sprintf("%d", i)}, float64(i))
		prometheusmetrics.CreateGauge("gpu_name", map[string]string{"gpu_name": fmt.Sprintf("%v", deviceName)}, float64(i))
		prometheusmetrics.CreateGauge("gpu_temperature", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(temperature))
		prometheusmetrics.CreateGauge("gpu_cpu_utilization", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(utilization.Gpu))
		prometheusmetrics.CreateGauge("gpu_mem_utilization", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(utilization.Memory))
		prometheusmetrics.CreateGauge("gpu_power_usage", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(gpuPowerUsage)/1000)
		prometheusmetrics.CreateGauge("gpu_running_process", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(len(runningProcess)))
	}
}

// Collectprometheus_metrics collects metrics for all the GPUs.
func CollectGpuMetricsV1() {
	deviceCount, err := nvml.DeviceGetCount()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device count", zap.Error(err))
		return
	}

	for i := 0; i < int(deviceCount); i++ {
		collectDeviceMetricsV1(i)
	}
}

// collectDeviceMetrics collects metrics for a single device.
func collectDeviceMetricsV1(deviceIndex int) {
	handle, err := nvml.DeviceGetHandleByIndex(deviceIndex)
	if err != nvml.SUCCESS {
		logger.Error("Error getting device handle", zap.Int("device_index", deviceIndex), zap.Error(err))
		return
	}

	deviceName, err := handle.GetName()
	if err != nvml.SUCCESS {
		logger.Error("Error getting device name", zap.Error(err))
		return
	}

	labels := map[string]string{
		"gpu_id":   fmt.Sprintf("%d", deviceIndex),
		"gpu_name": deviceName,
	}

	temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
	if err == nvml.SUCCESS {
		gauge.SetGaugeMetric("gpu_temperature", labels, float64(temperature))
	}

	utilization, err := handle.GetUtilizationRates()
	if err == nvml.SUCCESS {
		gauge.SetGaugeMetric("gpu_cpu_utilization", labels, float64(utilization.Gpu))
		gauge.SetGaugeMetric("gpu_mem_utilization", labels, float64(utilization.Memory))
	}

	gpuPowerUsage, err := handle.GetPowerUsage()
	if err == nvml.SUCCESS {
		gauge.SetGaugeMetric("gpu_power_usage", labels, float64(gpuPowerUsage)/1000) // Assuming power is in mW and we want W.
	}

	runningProcess, err := handle.GetComputeRunningProcesses()
	if err == nvml.SUCCESS {
		gauge.SetGaugeMetric("gpu_running_process", labels, float64(len(runningProcess)))
	}

	// Add more metrics here as needed.
	logger.Debug("Collected GPU metrics", zap.Int("device_index", deviceIndex))
}
