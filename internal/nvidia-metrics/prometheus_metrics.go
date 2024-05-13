package nvidiametrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	gpuId = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_id",
			Help: "ID of the GPU.",
		},
		[]string{"gpu_id"},
	)

	gpuName = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_name",
			Help: "Name of the GPU.",
		},
		[]string{"gpu_name"},
	)

	gpuUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_cpu_utilization",
			Help: "GPU utilization in percent.",
		},
		[]string{"gpu_id", "gpu_name"},
	)
	gpuMemoryUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_mem_utilization",
			Help: "GPU memory utilization in percent.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	gpuTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_temperature",
			Help: "Temperature of the GPU in degrees Celsius.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	gpuPowerUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_power_usage",
			Help: "Power usage of the GPU in watts.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	gpuRunningProcess = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_running_process",
			Help: "Number of running processes on the GPU.",
		},
		[]string{"gpu_id", "gpu_name"},
	)

	gpuUtilizationRates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "gpu_utilization_rates",
			Help: "Utilization rates of the GPU.",
		},
		[]string{"gpu_id", "gpu_name"},
	)
)

func init() {
	prometheus.MustRegister(gpuId)
	prometheus.MustRegister(gpuName)
	prometheus.MustRegister(gpuUtilization)
	prometheus.MustRegister(gpuTemperature)
	prometheus.MustRegister(gpuMemoryUtilization)
	prometheus.MustRegister(gpuPowerUsageMetric)
	prometheus.MustRegister(gpuRunningProcess)
	prometheus.MustRegister(gpuUtilizationRates)
}
