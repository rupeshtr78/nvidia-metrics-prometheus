package prometheusmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"testing"
)

func labelsHelper(t *testing.T) (labels GpuLabels, actual LabelsMap) {
	t.Helper()
	labels = GpuLabels{"label1": "gpu_id", "label2": "gpu_name"}
	labelsMap := map[string]GpuLabels{"metric1": labels}

	return labels, labelsMap

}

func TestCreateLabelsMap(t *testing.T) {
	// Assign
	want := make(LabelsMap)
	// Act
	got := CreateLabelsMap()
	// Assert
	if len(want) != len(got) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestLabelsMap_AddLabels(t *testing.T) {
	// Assign
	l := CreateLabelsMap()
	want, _ := labelsHelper(t)
	// Act
	l.AddLabels("metric1", want)
	got := l["metric1"]
	// Assert
	if len(want) != len(got) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestLabelsMap_GetLabelsFromMap(t *testing.T) {
	// Assign
	l := CreateLabelsMap()
	want, _ := labelsHelper(t)
	l.AddLabels("metric1", want)
	// Act
	actual, _ := l.GetLabelsFromMap("metric1")
	// Assert
	if len(want) != len(actual) {
		t.Errorf("Expected %v, got %v", want, actual)
	}
}

func TestLabelsMap_GetLabelsFromMap_InvalidMetric(t *testing.T) {
	// Assign
	l := CreateLabelsMap()
	want, _ := labelsHelper(t)
	l.AddLabels("metric1", want)
	// Act
	actual, _ := l.GetLabelsFromMap("metric2")
	// Assert
	if actual != nil {
		t.Errorf("Expected %v, got %v", nil, actual)
	}
}

func TestCreateMetricsMap(t *testing.T) {
	// Assign
	expected := make(MetricMap)
	// Act
	got := CreateMetricsMap()
	// Assert
	if len(expected) != len(got) {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestMetricMap_AddMetric(t *testing.T) {
	// Assign
	m := CreateMetricsMap()
	expected := &prometheus.GaugeVec{}

	// Act
	m.AddMetric("metric1", expected)
	got := m["metric1"]
	// Assert
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestMetricMap_AddMetric_InvalidGauge(t *testing.T) {
	// Assign
	m := CreateMetricsMap()
	want := &prometheus.GaugeVec{}
	// Act
	m.AddMetric("metric1", nil)
	got := m["metric1"]
	// Assert
	if want == got {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestMetricMap_GetMetric(t *testing.T) {
	// Assign
	m := CreateMetricsMap()
	expected := &prometheus.GaugeVec{}
	m.AddMetric("metric1", expected)
	// Act
	got, _ := m.GetMetric("metric1")
	// Assert
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestMetricMap_GetMetric_InvalidMetric(t *testing.T) {
	// Assign
	m := CreateMetricsMap()
	expected := &prometheus.GaugeVec{}
	m.AddMetric("metric1", expected)
	// Act
	got, _ := m.GetMetric("metric2")
	// Assert
	if expected == got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}
