package prometheusmetrics

import (
	"fmt"

	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// createGauge creates a new gauge with labels and sets the value
func CreateGauge(name string, labels map[string]string, value float64) error {
	// Get the gauge vector from the metrics map
	// check if the metric exists in prometheus
	if MetricsMap == nil {
		logger.Error("Metrics map is nil")
		return fmt.Errorf("metrics map is nil")
	}

	gaugeVec, ok := MetricsMap[name]
	if !ok {
		return fmt.Errorf("failed to find metric: %s", name)
	}

	// If registered, create a new gauge with labels
	gauge := gaugeVec.With(labels)
	// Set the value
	gauge.Set(float64(value))
	// Add the gauge to the gauge map
	GuageMap[name] = gauge
	logger.Debug("Created the gauge", zap.String("name", name), zap.Any("labels", labels), zap.Float64("value", value))

	return nil
}

// SetGaugeMetric sets a gauge metric with the given name, labels, and value.
func SetGaugeMetric(name string, labels map[string]string, value float64) {
	err := CreateGauge(name, labels, value)
	if err != nil {
		logger.Error("Failed to create gauge metric", zap.String("metric_name", name), zap.Error(err))
	}
}
