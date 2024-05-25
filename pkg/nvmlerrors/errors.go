package nvmlerrors

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

// Extend nvml.Return to include more error codes
type MetricReturn struct {
	Return      nvml.Return
	MetricError int32
}

const (
	ERR_NONE              = iota
	METRIC_NOT_REGISTERED = 1000
	LABEL_NOT_REGISTERED  = 1001
)

func (e *MetricReturn) Error() string {
	return e.String()
}

// String method
func (e *MetricReturn) String() string {
	// if metric error is less than or equal to the return error, return the return nvml error
	if e.MetricError <= int32(e.Return) {
		return e.Return.Error()
	}
	switch e.MetricError {
	case ERR_NONE:
		return nvml.SUCCESS.String()
	case METRIC_NOT_REGISTERED:
		return "ERR_METRIC_NOT_REGISTERED"
	default:
		return e.Return.Error()
	}
}

// ErrorInit is a function that initializes the MetricReturn struct
func ErrorInit() *MetricReturn {
	return &MetricReturn{}

}

// @TODO implement
// Errormain is a function that demonstrates the use of the MetricReturn struct
func Errormain() {
	me := ErrorInit()
	// me.MetricError = METRIC_NOT_REGISTERED
	fmt.Println(me.Error())

	me.Return = nvml.ERROR_OPERATING_SYSTEM
	fmt.Println(me.Error())

}
