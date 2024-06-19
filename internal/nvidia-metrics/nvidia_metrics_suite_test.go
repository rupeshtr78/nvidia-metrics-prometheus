package nvidiametrics_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNvidiaMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NvidiaMetrics Suite")
}
