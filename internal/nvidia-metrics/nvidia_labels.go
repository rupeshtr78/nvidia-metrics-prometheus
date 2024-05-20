package nvidiametrics

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// DeviceInfo: A function type for retrieving device information.
type DeviceInfo func(device nvml.Device) (any, nvml.Return)

type LabelFunctions map[string]DeviceInfo

func NewLabelFunction() LabelFunctions {
	l := make(LabelFunctions)
	return l
}

func NewDeviceInfo(f func(device nvml.Device) (any, nvml.Return)) DeviceInfo {
	return f
}

func (lf LabelFunctions) AddLabel(labelName string, f DeviceInfo) {
	if lf == nil {
		logger.Error("Label functions map is nil")
	}
	lf[labelName] = f
}

func (lf LabelFunctions) GetLabelFunc(labelName string) (DeviceInfo, error) {
	if f, ok := (lf)[labelName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("label function not found")
}

func (lf LabelFunctions) SetLabelFunc(labelName string, f DeviceInfo) {
	(lf)[labelName] = f
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

	logger.Debug("Label keys", zap.Any("label_keys", labelKeys))

	return labelKeys

}

// AddLabelFunction adds the label function to the map
func (lm LabelFunctions) AddFunctions() {

	labelFunc := NewLabelFunction()
	labelFunc.AddLabel(config.GPU_ID.GetLabel(), NewDeviceInfo(func(device nvml.Device) (any, nvml.Return) {
		return device.GetIndex()
	}))

	labelFunc.AddLabel(config.GPU_NAME.GetLabel(), NewDeviceInfo(func(device nvml.Device) (any, nvml.Return) {
		return device.GetName()
	}))

	logger.Debug("Label functions", zap.Any("label_functions", lm))

	// add the label function to the map

}

func (lf LabelFunctions) FetchDeviceLabelValue(device nvml.Device, labelName string) any {

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

func (lm LabelFunctions) GetLabelValue(device nvml.Device, labelName string) string {
	// get the label value
	value := lm.FetchDeviceLabelValue(device, labelName)
	return fmt.Sprintf("%v", value)
}

func (lm LabelFunctions) GetMetricLabelValues(device nvml.Device, metricName string) map[string]string {
	labelValues := GetLabelKeys(metricName)

	// iterate over the label functions and get the label values
	for labelName, _ := range labelValues {
		labelValues[labelName] = lm.GetLabelValue(device, labelName)
	}

	return labelValues
}
