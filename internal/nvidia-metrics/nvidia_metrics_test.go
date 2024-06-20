package nvidiametrics

import (
	"context"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	"github.com/stretchr/testify/mock"
)

// MockNvmlDevice is a mock implementation of the nvml.Device interface.
type MockNvmlDevice struct {
	mock.Mock
	nvml.Device
}

var ctx = context.TODO()

func (m *MockNvmlDevice) GetUtilizationRates() (nvml.Utilization, nvml.Return) {
	args := m.Called()
	return args.Get(0).(nvml.Utilization), args.Get(1).(nvml.Return)
}

var _ = Describe("GPUDeviceMetrics", func() {
	var (
		gpuDeviceMetrics *GPUDeviceMetrics
		mockHandle       *MockNvmlDevice
	)

	BeforeEach(func() {
		gpuDeviceMetrics = &GPUDeviceMetrics{}
		mockHandle = new(MockNvmlDevice)
		InitNVML()
		defer ShutdownNVML()
	})

	Context("CollectUtilizationMetrics", func() {
		It("should collect GPU and Memory utilization metrics correctly when GetUtilizationRates returns success", func() {
			utilization := nvml.Utilization{
				Gpu:    50,
				Memory: 60,
			}

			mockHandle.On("GetUtilizationRates").Return(utilization, nvml.SUCCESS).Once()
			mockHandle.On("SetDeviceMetric", mock.Anything, config.GPU_GPU_UTILIZATION, 50.0).Return().Once()
			mockHandle.On("SetDeviceMetric", mock.Anything, config.GPU_MEM_UTILIZATION, 60.0).Return().Once()

			err := gpuDeviceMetrics.CollectUtilizationMetrics(ctx, mockHandle)
			Expect(err).To(Equal(nvml.SUCCESS))

			Expect(gpuDeviceMetrics.GPUCPUUtilization).To(Equal(50.0))
			Expect(gpuDeviceMetrics.GPUMemUtilization).To(Equal(60.0))

			mockHandle.AssertExpectations(GinkgoT())
		})

		It("should not update metrics if GetUtilizationRates does not return success", func() {
			mockHandle.On("GetUtilizationRates").Return(nvml.Utilization{}, nvml.ERROR_UNKNOWN).Once()

			err := gpuDeviceMetrics.CollectUtilizationMetrics(ctx, mockHandle)
			Expect(err).To(Equal(nvml.ERROR_UNKNOWN))

			Expect(gpuDeviceMetrics.GPUCPUUtilization).To(Equal(float64(0)))
			Expect(gpuDeviceMetrics.GPUMemUtilization).To(Equal(float64(0)))

			mockHandle.AssertNotCalled(GinkgoT(), "SetDeviceMetric")
		})
	})
})
