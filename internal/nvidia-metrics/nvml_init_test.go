package nvidiametrics_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"

	nvidiametrics "github.com/rupeshtr78/nvidia-metrics/internal/nvidia-metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
)

type stubLogger struct {
	fatalCalled bool
}

func (s *stubLogger) Fatal(msg string, fields ...zapcore.Field) {
	s.fatalCalled = true
}

var _ = Describe("NvmlInit", func() {
	Context("When Init is called", func() {
		It("should initialize NVML", func() {
			stub := &stubLogger{}
			logger.GetLogger("debug", false, "")
			nvidiametrics.InitNVML()
			// Expect no fatal errors
			Expect(stub.fatalCalled).To(BeFalse())

		})
	})

})
