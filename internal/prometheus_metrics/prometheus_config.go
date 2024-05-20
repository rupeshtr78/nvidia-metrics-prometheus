package prometheusmetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var MetricsMap = make(map[string]*prometheus.GaugeVec)
var GuageMap = make(map[string]prometheus.Gauge)

type Metrics struct {
	MetricList []GpuMetric `yaml:"metrics"`
}

type GpuMetric struct {
	Name   string    `yaml:"name"`
	Help   string    `yaml:"help"`
	Type   string    `yaml:"type"`
	Labels GpuLabels `yaml:"labels"` // {label1: gpu_id, label2: gpu_name}
}

type GpuLabels map[string]string

// Labels for the metrics
type LabelsMap map[string]GpuLabels

// Metrics for the GPU
type MetricMap map[string]*prometheus.GaugeVec

// CreateLabelsMap creates a new LabelsMap
func CreateLabelsMap() LabelsMap {
	l := make(LabelsMap)
	return l
}

func (l *LabelsMap) AddLabels(metricName string, labels GpuLabels) {
	(*l)[metricName] = labels
}

func (l *LabelsMap) GetLabelsFromMap(metricName string) (GpuLabels, error) {
	if labels, ok := (*l)[metricName]; ok {
		return labels, nil
	}
	return nil, fmt.Errorf("labels %v not found in map", metricName)
}

func CreateMetricsMap() MetricMap {
	m := make(MetricMap)
	return m
}

func (m *MetricMap) AddMetric(metricName string, metric *prometheus.GaugeVec) {
	(*m)[metricName] = metric
}

func (m *MetricMap) GetMetric(metricName string) (*prometheus.GaugeVec, error) {
	if metric, ok := (*m)[metricName]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v not found in map", metricName)
}
