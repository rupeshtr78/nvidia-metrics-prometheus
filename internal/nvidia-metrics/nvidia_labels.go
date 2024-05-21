package nvidiametrics

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/rupeshtr78/nvidia-metrics/internal/config"
	prometheusmetrics "github.com/rupeshtr78/nvidia-metrics/internal/prometheus_metrics"
	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
)

// DeviceInfo A function type for retrieving device information.
type DeviceInfo func(device nvml.Device) (any, nvml.Return)

type LabelFunctions map[string]DeviceInfo

func NewLabelFunction() LabelFunctions {
	return make(LabelFunctions)
}

func NewDeviceInfo(f func(device nvml.Device) (any, nvml.Return)) DeviceInfo {
	return f
}

func (d DeviceInfo) ToString() string {
	return fmt.Sprintf("%v", d)
}

func (d DeviceInfo) GetFunction() func(device nvml.Device) (any, nvml.Return) {
	return d
}

func (lf LabelFunctions) Add(labelName string, f DeviceInfo) {
	lf[labelName] = f
}

func (lf LabelFunctions) GetLabelFunc(labelName string) (func(device nvml.Device) (any, nvml.Return), error) {
	if lf == nil {
		return nil, fmt.Errorf("label function map empty") // TODO: return error)
	}

	if f, ok := (lf)[labelName]; ok {
		logger.Debug("Label function found", zap.String("label_name", labelName))
		return f.GetFunction(), nil
	}
	return nil, fmt.Errorf("label function not found for label %s", labelName)
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
		// adding dummy values for now will be updates while setting gauge
		labelKeys[key] = i
	}

	return labelKeys

}

// AddFunctions adds the label function to the map
func (lf LabelFunctions) AddFunctions() {

	lf.Add(config.GPU_ID.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		index, ret := device.GetIndex()
		return index, ret
	})

	lf.Add(config.GPU_NAME.GetLabel(), func(device nvml.Device) (any, nvml.Return) {
		return device.GetName()
	})

	for labelName, labelFunc := range lf {
		logger.Debug("listing label function", zap.String("label_name", labelName), zap.String("label_func", labelFunc.ToString()))
	}

	// add the label function to the map

}

// FetchDeviceLabelValue fetches the label value for the given device and label name
func (lf LabelFunctions) FetchDeviceLabelValue(device nvml.Device, labelName string) any {

	labelFunc, err := lf.GetLabelFunc(labelName)
	if err != nil {
		return err
	}

	value, ret := labelFunc(device)
	if ret != nvml.SUCCESS {
		logger.Error("Error fetching label value", zap.String("label_name", labelName))
		return nil
	}
	return value

}

// GetLabelValue returns the label value for the given device and label name
func (lf LabelFunctions) GetLabelValue(device nvml.Device, labelName string) string {
	// get the label value
	value := lf.FetchDeviceLabelValue(device, labelName)
	return fmt.Sprintf("%v", value)
}

// GetMetricLabelValues returns all the label values for the given device and metric name
func (lf LabelFunctions) GetMetricLabelValues(device nvml.Device, metricName string) map[string]string {
	labelValues := GetLabelKeys(metricName)

	// iterate over the label functions and get the label values
	for labelName := range labelValues {
		labelValues[labelName] = lf.GetLabelValue(device, labelName)
	}

	return labelValues
}
