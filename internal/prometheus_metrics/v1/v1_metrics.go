package prometheusmetrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

var MetricsMap = make(map[string]*prometheus.GaugeVec)

var (
	GpuId = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_id",
			Help: "ID of the GPU.",
		},
		[]string{"gpu_id"},
	)

	GpuName = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_name",
			Help: "Name of the GPU.",
		},
		[]string{"gpu_name"},
	)

	GpuUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_cpu_utilization",
			Help: "GPU utilization in percent.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	GpuMemoryUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_mem_utilization",
			Help: "GPU memory utilization in percent.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	GpuTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_temperature",
			Help: "Temperature of the GPU in degrees Celsius.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	GpuPowerUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_power_usage",
			Help: "Power usage of the GPU in watts.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	GpuRunningProcess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_running_process",
			Help: "Number of running processes on the GPU.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	GpuUtilizationRates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_utilization_rates",
			Help: "Utilization rates of the GPU.",
		},
		[]string{"gpu_id", "gpu_name"},
	)
)

func init() {
	prometheus.MustRegister(GpuId)
	prometheus.MustRegister(GpuName)
	prometheus.MustRegister(GpuUtilization)
	prometheus.MustRegister(GpuTemperature)
	prometheus.MustRegister(GpuMemoryUtilization)
	prometheus.MustRegister(GpuPowerUsageMetric)
	prometheus.MustRegister(GpuRunningProcess)
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

// IsMetricRegistered checks if the collector is already registered with Prometheus.
func IsMetricRegisteredV2(collector prometheus.Collector) bool {
	err := prometheus.Register(collector)
	if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
		return true // Already registered
	}
	if err == nil {
		prometheus.Unregister(collector) // Unregister if it was just registered for the check
	}
	return false
}

// DeleteMetric reads from config/metrics.yaml and deleted prometheus metrics
func DeleteMetrics(filePath string) error {
	// 	// read from config/metrics.yaml
	m, err := prometheusmetrics.LoadFromYAML(filePath)
	if err != nil {
		logger.Error("Failed to read metrics yaml file", zap.String("file", filePath), zap.Error(err))
		return err
	}

	if len(m.MetricList) == 0 {
		logger.Error("No metrics found in the yaml file", zap.String("file", filePath))
		return fmt.Errorf("no metrics found in the yaml file")
	}

	// delete prometheus metrics from yaml
	for _, metric := range m.MetricList {
		err := DeleteMetrics(metric.Name)
		if err != nil {
			return err
		}

		// delete the metric from the metrics map
		delete(MetricsMap, metric.Name)
		logger.Info("Deleted metric", zap.String("metric", metric.Name))

	}

	return nil
}
