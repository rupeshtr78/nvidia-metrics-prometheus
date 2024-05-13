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

// CreateGuage creates a new gauge and registers it with Prometheus
func CreateGuage(name string, help string, labels prometheus.Labels) error {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
	)

	// Register the metric
	if err := prometheus.Register(gauge); err != nil {
		logger.Error("failed to register metric", zap.String("metric", name), zap.Error(err))
		return err
	}

	// Add the metric to the metrics map
	if GuageMap == nil {
		GuageMap = make(map[string]prometheus.Gauge)
	}
	GuageMap[name] = gauge

	return nil
}
