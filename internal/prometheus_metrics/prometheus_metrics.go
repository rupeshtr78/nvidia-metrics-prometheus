package prometheusmetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"github.com/rupeshtr78/nvidia-metrics/pkg/utils"
	"go.uber.org/zap"
)

var RegisteredMetrics = CreateMetricsMap()
var RegisteredLabels = CreateLabelsMap()

// RegisterMetric NewGaugeVec creates a new gauge vector and registers it with Prometheus.
func RegisterMetric(gpuMetric GpuMetric) (*prometheus.GaugeVec, error) {
	if gpuMetric.Type != "gauge" {
		err := fmt.Errorf("unsupported metric type: %s", gpuMetric.Type)
		logger.Error("unsupported metric type", zap.String("type", gpuMetric.Type))
		return nil, err
	}

	labels, err := GetGPuLabels(gpuMetric.Labels)
	if err != nil {
		logger.Error("failed to get labels", zap.Error(err))
		return nil, err
	}

	// Create a new gauge vector
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: gpuMetric.Name.GetMetric(),
			Help: gpuMetric.Help,
		},
		labels,
	)

	// Unregister first; if not registered, no operations will be performed
	if !prometheus.Unregister(gaugeVec) {
		logger.Warn("metric was already registered", zap.String("metric", gpuMetric.Name.GetMetric()))
	}

	err = prometheus.Register(gaugeVec)
	if err != nil {
		logger.Error("failed to register metric", zap.Error(err))
		return nil, err
	}

	logger.Info("Verified registration of", zap.String("metric", gpuMetric.Name.GetMetric()))
	return gaugeVec, nil
}

// CreatePrometheusMetrics reads from config/metrics.yaml and create prometheus metrics
func CreatePrometheusMetrics(filePath string) error {
	var m Metrics
	// 	// read from config/metrics.yaml

	err := utils.LoadFromYAMLV2(filePath, &m)
	if err != nil {
		logger.Fatal("Failed to read metrics yaml file", zap.String("file", filePath), zap.Error(err))
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

		RegisteredMetrics.AddMetric(metric.Name.GetMetric(), gaugeVec)
		// metricMap[metric.Name] = gaugeVec
		// MetricsMap[metric.Name] = gaugeVec
		// Add Labels to the labels map

		RegisteredLabels.AddLabels(metric.Name.GetMetric(), metric.Labels)

	}

	return nil
}
