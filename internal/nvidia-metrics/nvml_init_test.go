package nvidiametrics_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
)

var _ = Describe("NvmlInit", func() {
	Context("When Init is called", func() {
		It("should initialize NVML", func() {
			nvidiametrics.InitNVML()
			// Expect no fatal errors
			Expect(logger.Fatal).To(BeNil())

		})
	})

})
