package nvidiametrics

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"

	nvml "github.com/NVIDIA/go-nvml/pkg/nvml"
)

var _ = Describe("ShutdownNVML", func() {
	var (
		err error
	)

	JustBeforeEach(func() {
		err = nvml.Shutdown()
		if err != nil && err != nvml.SUCCESS {
			logger.Fatal("Failed to shutdown NVML", zap.Error(err))
		} else {
			logger.Info("Shutdown NVML Successfully")
		}
	})

	Context("when calling Shutdown", func() {
		It("should not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
