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
      []string{"gpu"},
      )

  gpuName = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "gpu_name",
      Help: "Name of the GPU.",
      },
      []string{"gpu"},
      )

  gpuUtilization = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "gpu_utilization",
      Help: "GPU utilization in percent.",
      },
      []string{"gpu"},
      )

  gpuTemperature = prometheus.NewGaugeVec(
    prometheus.GaugeOpts{
      Name: "gpu_temperature",
      Help: "Temperature of the GPU in degrees Celsius.",
      },
      []string{"gpu"},
      )
)

func init() {
  prometheus.MustRegister(gpuId)
  prometheus.MustRegister(gpuName)
  prometheus.MustRegister(gpuUtilization)
  prometheus.MustRegister(gpuTemperature)
}
