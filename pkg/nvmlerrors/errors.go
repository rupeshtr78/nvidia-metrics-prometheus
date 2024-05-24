package nvmlerrors

import (
	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

// Extend nvml.Return to include more error codes
type MetricReturn struct {
	Return      nvml.Return
	MetricError int32
}

const (
	ERR_NONE              = iota
	METRIC_NOT_REGISTERED = 100
)

const (
	ERR_NONE_STR              = "No error"
	ERR_METRIC_NOT_REGISTERED = "Metric not registered"
)

func (e *MetricReturn) Error() string {
	return e.String()
}

// String method
func (e *MetricReturn) String() string {

	if e.MetricError <= int32(e.Return) {
		return e.Return.Error()
	}
	switch e.MetricError {
	case ERR_NONE:
		return nvml.SUCCESS.String()
	case METRIC_NOT_REGISTERED:
		return ERR_METRIC_NOT_REGISTERED
	default:
		return e.Return.String()
	}
}

func HandleErrors(err *MetricReturn) {
	if err.MetricError != ERR_NONE {
		err.Return = nvml.ERROR_UNKNOWN
	} else if err.Return != nvml.SUCCESS {
		err.MetricError = int32(err.Return)
	}
}
