package utils

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromYAMLV2_WithValidFile(t *testing.T) {
	tmpFilePath := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")

	content := []byte("key: value")
	err := os.WriteFile(tmpFilePath, []byte(content), 0666)
	if err != nil {
		log.Fatalf("Unable to write test data to temp file: %v", err)
	}

	var model map[string]string
	err = LoadFromYAMLV2(tmpFilePath, &model)

	assert.NoError(t, err)
	assert.Equal(t, "value", model["key"])
}

func TestLoadFromYAMLV2_WithInvalidFile(t *testing.T) {
	var model map[string]string
	err := LoadFromYAMLV2("nonexistent.yaml", &model)

	assert.Error(t, err)
}

func TestLoadFromYAMLV2_WithInvalidYAML(t *testing.T) {
	tmpFilePath := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")

	content := []byte("value")
	err := os.WriteFile(tmpFilePath, []byte(content), 0666)
	if err != nil {
		log.Fatalf("Unable to write test data to temp file: %v", err)
	}

	var model map[string]string
	err = LoadFromYAMLV2(tmpFilePath, &model)

	assert.Error(t, err)
}

func TestLoadFromYAMLV2_WithInvalidModel(t *testing.T) {
	tmpFilePath := fmt.Sprintf("%s/%s", t.TempDir(), "test_data.yaml")

	content := []byte("key: value")
	err := os.WriteFile(tmpFilePath, []byte(content), 0666)
	if err != nil {
		log.Fatalf("Unable to write test data to temp file: %v", err)
	}

	var model string
	err = LoadFromYAMLV2(tmpFilePath, &model)

	assert.Error(t, err)
}
