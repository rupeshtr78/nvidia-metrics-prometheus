package prometheusmetrics

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rupeshtr78/nvidia-metrics/pkg/utils"
	"gopkg.in/yaml.v2"
)

var MetricsMap = make(map[string]*prometheus.GaugeVec)
var GuageMap = make(map[string]prometheus.Gauge)

type Metrics struct {
	MetricList []GpuMetric `yaml:"metrics"`
}

type MetricsV2 struct {
	MetricList []GpuMetricV2 `yaml:"metrics"`
}

type GpuMetric struct {
	Name   string   `yaml:"name"`
	Help   string   `yaml:"help"`
	Type   string   `yaml:"type"`
	Labels []string `yaml:"labels"`
}

type GpuMetricV2 struct {
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

// Metrics
type Metric string

const (
	GPU_CPU_UTILIZATION Metric = "gpu_cpu_utilization"
	GPU_MEM_UTILIZATION Metric = "gpu_mem_utilization"
	GPU_POWER_USAGE     Metric = "gpu_power_usage"
	GPU_RUNNING_PROCESS Metric = "gpu_running_process"
	GPU_TEMPERATURE     Metric = "gpu_temperature"
)

type Label string

// Labels
const (
	GPU_ID   Label = "gpu_id"
	GPU_NAME Label = "gpu_name"
)

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

func (m *MetricMap) GetMetricFromMap(metricName string) (*prometheus.GaugeVec, error) {
	if metric, ok := (*m)[metricName]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %v not found in map", metricName)
}

func LoadFromYAML(yamlFile string) (*Metrics, error) {
	file, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	app := &Metrics{}
	err = decoder.Decode(app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// GetLabelsForMetric returns the labels for the given metric name
func GetLabelsForMetric(metricName string, filePath string) ([]string, error) {

	// check if the metric exists in prometheus
	if _, ok := MetricsMap[metricName]; !ok {
		return nil, fmt.Errorf("metric %v not registered", metricName)
	}

	// read from config/metrics.yaml
	var m MetricsV2
	err := utils.LoadFromYAMLV2(filePath, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to read metrics yaml file %v", filePath)
	}

	// check if the metric exists in the yaml file
	if len(m.MetricList) == 0 {
		return nil, fmt.Errorf("metrics not found in the yaml file %v", filePath)
	}

	// get labels for the metric
	for _, metric := range m.MetricList {
		if metric.Name == metricName {
			labels, err := GetGPuLabels(metric.Labels)
			if err != nil {
				return nil, fmt.Errorf("failed to get labels %v", err)
			}
			return labels, nil
		}

	}

	return nil, fmt.Errorf("")

}

// GetGPuLabels returns the list of labels for the metric
func GetGPuLabels(labels GpuLabels) ([]string, error) {
	var labelList []string
	if labels == nil {
		return labelList, fmt.Errorf("no labels found")
	}

	for _, v := range labels {
		labelList = append(labelList, v)
	}
	return labelList, nil
}

// GetPromtheusLabels returns the prometheus labels
func GetPromtheusLabels(labels GpuLabels) (prometheus.Labels, error) {
	if labels == nil {
		return nil, fmt.Errorf("no labels found")
	}

	promLabels := prometheus.Labels(labels)
	return promLabels, nil
}
