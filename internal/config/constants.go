// This file contains the constants used in the metrcs.yaml file
// The collect metric functions has to be added in nividia-metrics.go file
// The label functions has to be added in nvidia_labels.go file
package config

// Metrics
type Metric string

const (
	GPU_ID_METRIC              Metric = "gpu_id_metric"
	GPU_GPU_UTILIZATION        Metric = "gpu_gpu_utilization"
	GPU_MEM_UTILIZATION        Metric = "gpu_mem_utilization"
	GPU_POWER_USAGE            Metric = "gpu_power_usage"
	GPU_RUNNING_PROCESS        Metric = "gpu_running_process"
	GPU_TEMPERATURE            Metric = "gpu_temperature"
	GPU_MEMORY_USED            Metric = "gpu_memory_used"
	GPU_MEMORY_TOTAL           Metric = "gpu_memory_total"
	GPU_MEMORY_FREE            Metric = "gpu_memory_free"
	GPU_P_STATE                Metric = "gpu_p_state"
	GPU_MEMORY_CLOCK           Metric = "gpu_memory_clock"
	GPU_ECC_CORRECTED_ERRORS   Metric = "gpu_ecc_corrected_errors"
	GPU_ECC_UNCORRECTED_ERRORS Metric = "gpu_ecc_uncorrected_errors"
)

type Label string

// Labels
const (
	GPU_ID             Label = "gpu_id"
	GPU_NAME           Label = "gpu_name"
	GPU_TEM_THRESHOLD  Label = "gpu_temperature_threshold"
	GPU_MEM_CLOCK_MAX  Label = "gpu_memory_clock_max"
	GPU_CORES          Label = "gpu_cores"
	GPU_DRIVER_VERSION Label = "gpu_driver_version"
	GPU_CUDA_VERSION   Label = "gpu_cuda_version"
)

func (m Metric) GetMetric() string {
	return string(m)
}

func (l Label) GetLabel() string {
	return string(l)
}
