package prometheusmetrics

import (
	"fmt"

	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// createGauge creates a new gauge with labels and sets the value
func CreateGauge(name string, labels GpuLabels, value float64) error {
	// Get the gauge vector from the metrics map
	// check if the metric exists in prometheus

	if MetricsMap == nil {
		logger.Error("Metrics map is nil")
		return fmt.Errorf("metrics map is nil")
	}

	gaugeVec, err := RegisteredMetrics.GetMetric(name)
	if err != nil {
		logger.Error("Failed to get metric from map", zap.Error(err))
		return err
	}

	// get prometheus labels
	gpuLabels, err := GetPromtheusLabels(labels)
	if err != nil {
		logger.Error("Failed to get prometheues labels", zap.Error(err))
		return err
	}

	// If registered, create a new gauge with labels
	gauge := gaugeVec.With(gpuLabels)
	// Set the value
	gauge.Set(value)

	logger.Debug("Created the gauge", zap.String("name", string(name)), zap.Any("labels", labels), zap.Float64("value", value))

	return nil
}

// SetGaugeMetric sets a gauge metric with the given name, labels, and value.
func SetGaugeMetric(name string, labels GpuLabels, value float64) {
	err := CreateGauge(name, labels, value)
	if err != nil {
		logger.Error("Failed to create gauge metric", zap.String("metric_name", string(name)), zap.Error(err))
	}
}
