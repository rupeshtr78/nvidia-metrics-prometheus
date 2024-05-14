package nvidiametrics

// import (
// 	"fmt"

// 	"github.com/NVIDIA/go-nvml/pkg/nvml"
// 	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
// 	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
// 	"go.uber.org/zap"
// )

// var err = logger.GetLogger()

// // GetGPUUtilization returns the GPU utilization of the device.
// // The GPU utilization is the percent of time over the past sample period during which one or more kernels was executing on the GPU.

// func InitNVML() {
// 	err := nvml.Init()
// 	if err != nvml.SUCCESS {
// 		logger.Fatal("Failed to initialize NVML", zap.Error(err))

// 	}

// 	logger.Info("Initialized NVML")
// }

// func ShutdownNVML() {
// 	err := nvml.Shutdown()
// 	if err != nvml.SUCCESS {
// 		logger.Fatal("Failed to shutdown NVML", zap.Error(err))
// 	}

// 	logger.Info("Shutdown NVML")
// }

// func CollectGpuMetrics() {
// 	nvml.Init()
// 	defer nvml.Shutdown()

// 	deviceCount, err := nvml.DeviceGetCount()
// 	if err != nvml.SUCCESS {
// 		fmt.Println("Error getting device count:", err)
// 		return
// 	}

// 	// Loop through all the devices
// 	for i := 0; i < int(deviceCount); i++ {
// 		handle, err := nvml.DeviceGetHandleByIndex(i)
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting device handle:", err)
// 			continue
// 		}

// 		deviceName, err := handle.GetName()
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting device name:", err)
// 			continue
// 		}

// 		temperature, err := handle.GetTemperature(nvml.TEMPERATURE_GPU)
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting temperature:", err)
// 			continue
// 		}

// 		utilization, err := handle.GetUtilizationRates()
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting utilization rates:", err)
// 			continue
// 		}

// 		gpuPowerUsage, err := handle.GetPowerUsage()
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting power usage:", err)
// 			continue
// 		}

// 		runningProcess, err := handle.GetComputeRunningProcesses()
// 		if err != nvml.SUCCESS {
// 			fmt.Println("Error getting running processes:", err)
// 			continue
// 		}

// 		logger.Debug("GPU metrics", zap.String("name", deviceName), zap.Uint("temperature", uint(temperature)))

// 		// Set the prometheus metrics for the GPU maps to label values and sets the value
// 		// metrics.GpuId.WithLabelValues(fmt.Sprintf("%d", i)).Set(float64(i))
// 		// metrics.GpuName.WithLabelValues(fmt.Sprintf("%v", deviceName)).Set(float64(i))
// 		// metrics.GpuTemperature.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(temperature))
// 		// metrics.GpuUtilization.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(utilization.Gpu))
// 		// metrics.GpuMemoryUtilization.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(utilization.Memory))
// 		// metrics.GpuPowerUsageMetric.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(gpuPowerUsage) / 1000)
// 		// metrics.GpuRunningProcess.WithLabelValues(fmt.Sprintf(("%d"), i), fmt.Sprintf("%v", deviceName)).Set(float64(len(runningProcess)))

// 		prometheusmetrics.CreateGauge("gpu_id", map[string]string{"gpu_id": fmt.Sprintf("%d", i)}, float64(i))
// 		prometheusmetrics.CreateGauge("gpu_name", map[string]string{"gpu_name": fmt.Sprintf("%v", deviceName)}, float64(i))
// 		prometheusmetrics.CreateGauge("gpu_temperature", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(temperature))
// 		prometheusmetrics.CreateGauge("gpu_cpu_utilization", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(utilization.Gpu))
// 		prometheusmetrics.CreateGauge("gpu_mem_utilization", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(utilization.Memory))
// 		prometheusmetrics.CreateGauge("gpu_power_usage", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(gpuPowerUsage)/1000)
// 		prometheusmetrics.CreateGauge("gpu_running_process", map[string]string{"gpu_id": fmt.Sprintf("%d", i), "gpu_name": fmt.Sprintf("%v", deviceName)}, float64(len(runningProcess)))
// 	}
// }
