package prometheusmetrics

import (
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func LoadFromYAML(yamlFile string) (*Metrics, error) {
	file, err := os.Open(yamlFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	app := &Metrics{}
	err = decoder.Decode(app)
	if err != nil {
		return nil, err
	}

	return app, nil
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
