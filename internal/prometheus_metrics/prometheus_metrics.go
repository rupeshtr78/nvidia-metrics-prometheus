package prometheusmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

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
