package config

// Metrics
type Metric string

const (
	GPU_ID_METRIC       Metric = "gpu_id"
	GPU_CPU_UTILIZATION Metric = "gpu_cpu_utilization"
	GPU_MEM_UTILIZATION Metric = "gpu_mem_utilization"
	GPU_POWER_USAGE     Metric = "gpu_power_usage"
	GPU_RUNNING_PROCESS Metric = "gpu_running_process"
	GPU_TEMPERATURE     Metric = "gpu_temperature"
	GPU_MEMORY_USED     Metric = "gpu_memory_used"
	GPU_MEMORY_TOTAL    Metric = "gpu_memory_total"
	GPU_MEMORY_FREE     Metric = "gpu_memory_free"
)

type Label string

// Labels
const (
	GPU_ID            Label = "gpu_id"
	GPU_NAME          Label = "gpu_name"
	GPU_TEM_THRESHOLD Label = "gpu_temperature_threshold"
)

func (m Metric) GetMetric() string {
	return string(m)
}

func (l Label) GetLabel() string {
	return string(l)
}
