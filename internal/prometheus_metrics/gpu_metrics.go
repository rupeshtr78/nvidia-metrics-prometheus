package prometheusmetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var MetricsMap = make(map[string]*prometheus.GaugeVec)
var GuageMap = make(map[string]prometheus.Gauge)

type Metrics struct {
	MetricList []GpuMetric `yaml:"metrics"`
}

type GpuMetric struct {
	Name   string   `yaml:"name"`
	Help   string   `yaml:"help"`
	Type   string   `yaml:"type"`
	Labels []string `yaml:"labels"`
}

// NewGaugeVec creates a new gauge vector and registers it with Prometheus.
func RegisterMetric(gpuMetric GpuMetric) (*prometheus.GaugeVec, error) {
	if gpuMetric.Type != "gauge" {
		err := fmt.Errorf("unsupported metric type: %s", gpuMetric.Type)
		logger.Error("unsupported metric type", zap.String("type", gpuMetric.Type))
		return nil, err
	}

	// Create a new gauge vector
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: gpuMetric.Name,
			Help: gpuMetric.Help,
		},
		gpuMetric.Labels,
	)

	// Check if the metric is already registered
	registered := IsMetricRegisteredV2(gaugeVec)

	if !registered {
		// Register the metric
		if err := prometheus.Register(gaugeVec); err != nil {
			logger.Error("failed to register metric", zap.String("metric", gpuMetric.Name), zap.Error(err))
			return nil, err
		}
	}

	logger.Info("Registered metric", zap.String("metric", gpuMetric.Name))
	return gaugeVec, nil
}

// CreatePrometheusMetrics reads from config/metrics.yaml and create prometheus metrics
func CreatePrometheusMetrics(filePath string) error {
	// 	// read from config/metrics.yaml
	m, err := LoadFromYAML(filePath)
	if err != nil {
		logger.Error("Failed to read metrics yaml file", zap.String("file", filePath), zap.Error(err))
		return err
	}

	if len(m.MetricList) == 0 {
		logger.Error("No metrics found in the yaml file", zap.String("file", filePath))
		return fmt.Errorf("no metrics found in the yaml file")
	}

	// create prometheus metrics from yaml
	for _, metric := range m.MetricList {
		gaugeVec, err := RegisterMetric(metric)
		if err != nil {
			return err
		}

		// Add the metric to the metrics map
		MetricsMap[metric.Name] = gaugeVec

	}

	return nil
}

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

// IsMetricRegistered attempts to register a collector and checks if it is already registered.
func IsMetricRegistered(collector prometheus.Collector) (bool, error) {
	err := prometheus.Register(collector)
	if err == nil {
		return false, nil // Successfully registered, not previously registered
	}

	if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
		return true, nil // Already registered
	}

	// Other kinds of errors indicate different problems with registration.
	return false, err
}

// IsMetricRegistered checks if the collector is already registered with Prometheus.
func IsMetricRegisteredV2(collector prometheus.Collector) bool {
	err := prometheus.Register(collector)
	if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
		return true // Already registered
	}
	if err == nil {
		prometheus.Unregister(collector) // Unregister if it was just registered for the check
	}
	return false
}

// DeleteMetric reads from config/metrics.yaml and deleted prometheus metrics
func DeleteMetrics(filePath string) error {
	// 	// read from config/metrics.yaml
	m, err := LoadFromYAML(filePath)
	if err != nil {
		logger.Error("Failed to read metrics yaml file", zap.String("file", filePath), zap.Error(err))
		return err
	}

	if len(m.MetricList) == 0 {
		logger.Error("No metrics found in the yaml file", zap.String("file", filePath))
		return fmt.Errorf("no metrics found in the yaml file")
	}

	// delete prometheus metrics from yaml
	for _, metric := range m.MetricList {
		err := DeleteMetrics(metric.Name)
		if err != nil {
			return err
		}

		// delete the metric from the metrics map
		delete(MetricsMap, metric.Name)

	}

	return nil
}
