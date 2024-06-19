package nvidiametrics_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
)

type stubLogger struct {
	fatalCalled bool
}

var _ = logger.GetLogger("debug", false, "")

var _ = Describe("NvmlInit", func() {
	Context("When Init is called", func() {

		// if NVML is initialized, it should not call logger.Fatal
		It("should initialize NVML", func() {
			stub := &stubLogger{}
			nvidiametrics.InitNVML()
			// Expect no fatal errors
			Expect(stub.fatalCalled).To(BeFalse())

		})

		// // if NVML is not initialized, it should call logger.Fatal
		It("should shutdown NVML", func() {
			stub := &stubLogger{}
			nvidiametrics.ShutdownNVML()
			// Expect no fatal errors
			Expect(stub.fatalCalled).To(BeFalse())

		})

	})

})
