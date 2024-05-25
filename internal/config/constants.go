// This file contains the constants used in the metrcs.yaml file
// The collect metric functions has to be added in nividia-metrics.go file
// The label functions has to be added in nvidia_labels.go file
package config

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

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
	GPU_GRAPHICS_CLOCK         Metric = "gpu_graphics_clock"
	GPU_SM_CLOCK               Metric = "gpu_sm_clock"
	GPU_VIDEO_CLOCK            Metric = "gpu_video_clock"
	GPU_ECC_CORRECTED_ERRORS   Metric = "gpu_ecc_corrected_errors"
	GPU_ECC_UNCORRECTED_ERRORS Metric = "gpu_ecc_uncorrected_errors"
	GPU_FAN_SPEED              Metric = "gpu_fan_speed"
	GPU_PEAK_FLOPS_METRIC      Metric = "gpu_peak_flops_metric"
)

type Label string

// Labels
const (
	GPU_ID                 Label = "gpu_id"
	GPU_NAME               Label = "gpu_name"
	GPU_TEM_THRESHOLD      Label = "gpu_temperature_threshold"
	GPU_MEM_CLOCK_MAX      Label = "gpu_memory_clock_max"
	GPU_SM_CLOCK_MAX       Label = "gpu_sm_clock_max"
	GPU_GRAPHICS_CLOCK_MAX Label = "gpu_graphics_clock"
	GPU_CORES              Label = "gpu_cores"
	GPU_DRIVER_VERSION     Label = "gpu_driver_version"
	GPU_CUDA_VERSION       Label = "gpu_cuda_version"
	GPU_PEAK_FLOPS         Label = "gpu_peak_flops"
)

func (m Metric) GetMetric() string {
	return string(m)
}

func (l Label) GetLabel() string {
	return string(l)
}

type GpuFlops struct {
	GPUId  int
	TFLOPS float64
}

const (
	GPU_PEAK_FLOPS_P40      float64 = 12e12   // TF32 FLOPS P40
	GPU_PEAK_FLOPS_RTX_3060 float64 = 1.28e13 // 1.28e13 RTX 3060
)

// GpuFlopsData returns the GPU flops data add more GPUs here and TFOPS above
func GpuFlopsData() []GpuFlops {
	gpus := []GpuFlops{
		{GPUId: 0, TFLOPS: GPU_PEAK_FLOPS_RTX_3060},
		{GPUId: 1, TFLOPS: GPU_PEAK_FLOPS_P40},
	}
	return gpus
}

// GetGpuFlops returns the GPU flops for the given GPU ID
func GetGpuFlops(gpuId int) (float64, nvml.Return) {
	if gpuId < 0 {
		return 0, nvml.ERROR_INVALID_ARGUMENT
	}

	if len(GpuFlopsData()) == 0 {
		return 0, nvml.ERROR_NO_DATA
	}

	for _, gpu := range GpuFlopsData() {
		if gpu.GPUId == gpuId {
			if gpu.TFLOPS == 0 {
				return 0, nvml.ERROR_NOT_FOUND
			}
			return gpu.TFLOPS, nvml.SUCCESS
		}
	}
	return 0, nvml.ERROR_NOT_FOUND
}
