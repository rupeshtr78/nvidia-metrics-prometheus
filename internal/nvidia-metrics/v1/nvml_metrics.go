package nvidiametrics

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

type MetricCollector struct {
    metricConfigKey config.Metric

    collectionFunc  func(handle nvml.Device, metrics *nvidiametrics.GPUDeviceMetrics, configKey string) error
}

func (mc *MetricCollector) CollectIfNeeded(handle nvml.Device, metrics *nvidiametrics.GPUDeviceMetrics) {
   m := mc.metricConfigKey
	if isRegistered(m) {
        err := mc.collectionFunc(handle, metrics, mc.metricConfigKey.GetMetric())
        if err != nvml.SUCCESS {
            logger.Error("Error collecting metric", zap.String("metric_key", mc.metricConfigKey.GetMetric()), zap.Error(err))
        }
    }
}

func isRegistered(metric config.Metric) bool {

	if _, ok := prometheusmetrics.RegisteredMetrics[metric.GetMetric()]; !ok {
		logger.Warn("metric not registered", zap.String("metric", metric.GetMetric()))
		return false
	}
	return true
}

// func collectDeviceMetrics(deviceIndex int) (*GPUDeviceMetrics, error) {
//     // Existing setup code...

//     // Example usage for temperature metrics, assuming CollectTemperatureMetrics matches the expected function signature
//     tempCollector := &MetricCollector{
//         metricConfigKey: config.GPU_TEMPERATURE,
//         collectionFunc:  CollectTemperatureMetrics,
//     }
//     tempCollector.CollectIfNeeded(handle, metrics)

//     // Call unconditional metric collections directly
//     err = CollectUtilizationMetrics(handle, metrics)
//     if err != nvml.SUCCESS {
//         logger.Error("Error collecting utilization metrics", zap.Error(err))
//     }

//     // Continue with other metrics...

//     // For ECC corrected errors
//     eccCorrectedCollector := &MetricCollector{
//         metricConfigKey: config.GPU_ECC_CORRECTED_ERRORS,
//         collectionFunc:  collectEccCorrectedErrorsMetrics, // Adjust based on actual function
//     }
//     eccCorrectedCollector.CollectIfNeeded(handle, metrics)

//     // Repeat for other conditional metrics using similar pattern

//     return metrics, nil
// }
