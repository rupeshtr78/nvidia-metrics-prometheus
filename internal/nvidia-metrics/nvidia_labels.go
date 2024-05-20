package nvidiametrics

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// DeviceInfoRetriever: A function type for retrieving device information.
type DeviceInfoRetriever func(device nvml.Device) (any, nvml.Return)

type LabelFunctions map[string]DeviceInfoRetriever

func NewLabelFunction() *LabelFunctions {
	l := make(LabelFunctions)
	return &l
}

func NewDeviceInfoRetriever(f func(device nvml.Device) (any, nvml.Return)) DeviceInfoRetriever {
	return DeviceInfoRetriever(f)
}

func (lf *LabelFunctions) AddLabel(labelName string, f DeviceInfoRetriever) {
	(*lf)[labelName] = f
}

func (lf *LabelFunctions) GetLabelFunc(labelName string) (DeviceInfoRetriever, error) {
	if f, ok := (*lf)[labelName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("label function not found")
}

func (lf *LabelFunctions) SetLabelFunc(labelName string, f DeviceInfoRetriever) {
	(*lf)[labelName] = f
}

func (lf *LabelFunctions) FetchDeviceLabelValue(device nvml.Device, labelName string) any {

	labelFunc, err := lf.GetLabelFunc(labelName)
	if err != nil {
		logger.Error("Error fetching label function", zap.String("label_name", labelName))
		return nil
	}

	value, ret := labelFunc(device)
	if ret != nvml.SUCCESS {
		logger.Error("Error fetching label value", zap.String("label_name", labelName))
		return nil
	}
	return value

}

// example usage
func GetLabelValue(device nvml.Device, name string) {

	labelFunc := NewLabelFunction()
	labelFunc.AddLabel(string(config.GPU_ID), NewDeviceInfoRetriever(func(device nvml.Device) (any, nvml.Return) {
		return device.GetIndex()
	}))

	// get the label value
	value := labelFunc.FetchDeviceLabelValue(device, name)
	_ = value

}

// GetLabelKeys returns the label keys for the given metric name.
// Example: "key":"gpu_power_usage","label":{"label1":"gpu_id","label2":"gpu_name"}}
func GetLabelKeys(metricName string) map[string]string {
	labelKeys := make(map[string]string)

	if _, ok := prometheusmetrics.RegisteredLabels[metricName]; !ok {
		logger.Warn("Metric not found in registered labels", zap.String("metric_name", metricName))
		return labelKeys
	}

	keys := prometheusmetrics.RegisteredLabels[metricName]
	// iterate over the keys and add label name to the map
	for i, key := range keys {
		if len(key) == 0 {
			continue
		}
		// adding dummy values for now will be updates while setting guage
		labelKeys[key] = i
	}

	return labelKeys

}
