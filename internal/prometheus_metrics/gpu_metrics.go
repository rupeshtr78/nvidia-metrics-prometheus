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
	gaugeVec, ok := MetricsMap[name]
	if !ok {
		logger.Error("Failed to find metric", zap.String("metric", name))
		return
	}

	// Set the prometheus metrics for the GPU maps to label values and sets the value
	gaugeVec.With(labels).Set(value)
}
