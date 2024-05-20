package prometheusmetrics

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
