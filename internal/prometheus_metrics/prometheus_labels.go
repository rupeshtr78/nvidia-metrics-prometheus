package prometheusmetrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rupeshtr78/nvidia-metrics/pkg/utils"
)

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

// GetLabelsForMetric returns the labels for the given metric name
func GetLabelsForMetric(metricName string, filePath string) ([]string, error) {

	// check if the metric exists in prometheus
	if _, ok := RegisteredMetrics[metricName]; !ok {
		return nil, fmt.Errorf("metric %v not registered", metricName)
	}

	// read from config/metrics.yaml
	var m Metrics
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
