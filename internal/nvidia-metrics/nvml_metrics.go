package nvidiametrics

import (


	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
)

// GetGPUUtilization returns the GPU utilization of the device.
// The GPU utilization is the percent of time over the past sample period during which one or more kernels was executing on the GPU.

func initNVML() {
	err := nvml.Init()
	if err != nil {

	}
}
