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
func RegisterMetric(gpuMetric GpuMetric) error {
	if gpuMetric.Type != "gauge" {
		err := fmt.Errorf("unsupported metric type: %s", gpuMetric.Type)
		logger.Error("unsupported metric type", zap.String("type", gpuMetric.Type))
		return err
	}

	// Create a new gauge vector
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: gpuMetric.Name,
			Help: gpuMetric.Help,
		},
		gpuMetric.Labels,
	)

	// Register the metric
	if err := prometheus.Register(gaugeVec); err != nil {
		logger.Error("failed to register metric", zap.String("metric", gpuMetric.Name), zap.Error(err))
		return err
	}

	// Add the metric to the metrics map
	if MetricsMap == nil {
		MetricsMap = make(map[string]*prometheus.GaugeVec)
	}

	// Add the metric to the metrics map
	MetricsMap[gpuMetric.Name] = gaugeVec

	return nil

}

// CreatePrometheusMetrics reads from config/metrics.yaml and create prometheus metrics
func CreatePrometheusMetrics(filePath string) {
	// 	// read from config/metrics.yaml
	m, err := LoadFromYAML(filePath)
	if err != nil {
		logger.Error("Failed to load metrics from yaml file", zap.String("file", filePath), zap.Error(err))
		return
	}

	if len(m.MetricList) == 0 {
		logger.Error("No metrics found in the yaml file", zap.String("file", filePath))
		return
	}

	// create prometheus metrics from yaml
	for _, metric := range m.MetricList {
		err := RegisterMetric(metric)
		if err != nil {
			logger.Error("Failed to create prometheus metric", zap.Error(err))
		}
	}
}

// createGauge creates a new gauge with labels and sets the value
func CreateGauge(name string, labels map[string]string, value float64) {
	// Get the gauge vector from the metrics map
	// check if the metric exists in prometheus

	gaugeVec, ok := MetricsMap[name]
	if !ok {
		logger.Error("Failed to find metric", zap.String("metric", name))
		return
	}

	// Check if the metric is already registered
	registered, err := IsMetricRegistered(gaugeVec)
	if err != nil && !registered {
		logger.Error("Not registered metric", zap.String("metric", name), zap.Error(err))
		return
	}

	// If registered, get the existing gauge
	if m, ok := MetricsMap[name]; ok {
		// Create a new gauge with labels
		gauge := m.With(labels)
		// Set the value
		gauge.Set(value)
		// Add the gauge to the gauge map
		GuageMap[name] = gauge
	} else {
		logger.Error("Failed to create gauge", zap.String("metric", name))
		return
	}
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
