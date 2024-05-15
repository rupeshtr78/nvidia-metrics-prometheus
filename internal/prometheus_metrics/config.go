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

// Labels for the metrics
var LabelsMap = make(map[string]map[string]string)

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
	Name   string            `yaml:"name"`
	Help   string            `yaml:"help"`
	Type   string            `yaml:"type"`
	Labels map[string]string `yaml:"labels"` // {label1: gpu_id, label2: gpu_name}
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
		return nil, fmt.Errorf("Failed to read metrics yaml file %v", filePath)
	}

	// check if the metric exists in the yaml file
	if len(m.MetricList) == 0 {
		return nil, fmt.Errorf("No metrics found in the yaml file %v", filePath)
	}

	// get labels for the metric
	for _, metric := range m.MetricList {
		if metric.Name == metricName {
			labels, err := utils.GetLabelList(metric.Labels)
			if err != nil {
				return nil, fmt.Errorf("failed to get labels %v", err)
			}
			return labels, nil
		}

	}

	return nil, fmt.Errorf("")

}
