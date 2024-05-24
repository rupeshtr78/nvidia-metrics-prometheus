package prometheusmetrics

import (
	"fmt"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mock logger

var getGPuLabelsFunc func(GpuLabels) ([]string, error)

func TestRegisterMetric(t *testing.T) {

	defer func() {
		getGPuLabelsFunc = GetGPuLabels
	}()

	tests := []struct {
		name        string
		gpuMetric   GpuMetric
		mockLabels  func(GpuLabels) ([]string, error)
		expectError bool
	}{
		{
			name: "ValidGauge",
			gpuMetric: GpuMetric{
				Name: config.Metric("gpu_utilization"),
				Help: "GPU utilization",
				Type: "gauge",
				Labels: GpuLabels{
					"label1": "gpu_id",
					"label2": "gpu_name",
				},
			},
			mockLabels: func(labels GpuLabels) ([]string, error) {
				return []string{"gpu_id", "gpu_name"}, nil
			},
			expectError: false,
		},
		{
			name: "UnsupportedType",
			gpuMetric: GpuMetric{
				Name: config.Metric("gpu_utilization"),
				Help: "GPU utilization",
				Type: "counter", // Unsupported type
				Labels: GpuLabels{
					"label1": "gpu_id",
					"label2": "gpu_name",
				},
			},
			mockLabels:  nil,
			expectError: true,
		},
		{
			name: "FailedToGetLabels",
			gpuMetric: GpuMetric{
				Name: config.Metric("gpu_utilization"),
				Help: "GPU utilization",
				Type: "gauge",
				Labels: GpuLabels{
					"label1": "", // Invalid label
				},
			},
			mockLabels: func(labels GpuLabels) ([]string, error) {
				return nil, fmt.Errorf("failed to get labels")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockLabels != nil {
				getGPuLabelsFunc = tt.mockLabels
			}

			// Act
			gaugeVec, err := RegisterMetric(tt.gpuMetric)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "Expected error for test case: %s", tt.name)
				assert.Nil(t, gaugeVec, "Expected gaugeVec to be nil for test case: %s", tt.name)
			} else {
				assert.NoError(t, err, "Expected no error for test case: %s", tt.name)
				assert.NotNil(t, gaugeVec, "Expected gaugeVec to be non-nil for test case: %s", tt.name)
			}
		})
	}
}

func TestCreatePrometheusMetrics_WithValidFile(t *testing.T) {

	tmpFilePath2 := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")
	content := []byte(`
metrics:
  - name: gpu_name_test
    type: gauge
    help: "Name of the GPU."
    labels:
      label1: gpu_name`)

	_ = os.WriteFile(tmpFilePath2, []byte(content), 0666)

	err := CreatePrometheusMetrics(tmpFilePath2)

	m := RegisteredMetrics
	got, _ := m.GetMetric("gpu_name_test")
	assert.NoError(t, err)
	assert.NotNil(t, got)
}

func TestCreatePrometheusMetrics_WithInvalidFile(t *testing.T) {
	err := CreatePrometheusMetrics("nonexistent.yaml")

	assert.Error(t, err, "error opening file")
}

func TestCreatePrometheusMetrics_WithInvalidYAML(t *testing.T) {
	tmpFilePath2 := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")
	content := []byte(`
metrics:
  name: gpu_name_test
   - type: gauge
    help: "Name of the GPU."
    labels:
      label1: gpu_name`)

	_ = os.WriteFile(tmpFilePath2, []byte(content), 0666)

	err := CreatePrometheusMetrics(tmpFilePath2)

	assert.Error(t, err)
}

func TestCreatePrometheusMetrics_WithNoMetrics(t *testing.T) {
	tmpfile := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")

	content := []byte("metricList: []")
	_ = os.WriteFile(tmpfile, []byte(content), 0666)

	err := CreatePrometheusMetrics(tmpfile)

	assert.Error(t, err)
}

func TestCreatePrometheusMetrics_WithUnsupportedMetricType(t *testing.T) {
	tmpFilePath2 := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")
	content := []byte(`
metrics:
  - name: gpu_name_test
    type: histogram
    help: "Name of the GPU."
    labels:
      label1: gpu_name`)

	_ = os.WriteFile(tmpFilePath2, []byte(content), 0666)

	err := CreatePrometheusMetrics(tmpFilePath2)

	assert.Error(t, err)
}
