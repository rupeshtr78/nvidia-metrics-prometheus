package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/rupeshtr78/nvidia-metrics/pkg/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func GetLabelList(labels map[string]string) ([]string, error) {
	var labelList []string
	if labels == nil {
		return labelList, fmt.Errorf("labels are nil")
	}

	for _, v := range labels {
		labelList = append(labelList, v)
	}
	return labelList, nil
}

func LoadFromYAMLV2(fileName string, model interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		logger.Error("error opening file", zap.String("fileName", fileName), zap.Error(err))
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("error closing file", zap.String("fileName", fileName), zap.Error(err))
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(model)
	if err != nil {
		logger.Error("error decoding yaml", zap.String("file", fileName), zap.Error(err))
		return err
	}

	return nil
}

// GetProjectDir will return the directory where the project is
func GetProjectDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return wd, err
	}
	wd = strings.Replace(wd, "/internal/utils", "", -1)
	return wd, nil
}
