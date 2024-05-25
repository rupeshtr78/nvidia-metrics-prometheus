package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetLogger(t *testing.T) {
	err := GetLogger("info")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("Expected logger to be initialized, got nil")
	}
}

func TestLogLevels(t *testing.T) {

	tests := []struct {
		level   string
		message string
	}{
		{"debug", "This is a debug message"},
		{"info", "This is an info message"},
		{"warn", "This is a warn message"},
		{"error", "This is an error message"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			switch tt.level {
			case "debug":
				Debug(tt.message)
			case "info":
				Info(tt.message)
			case "warn":
				Warn(tt.message)
			case "error":
				Error(tt.message)
			}
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected zapcore.Level
	}{
		{"debug", zap.DebugLevel},
		{"info", zap.InfoLevel},
		{"warn", zap.WarnLevel},
		{"error", zap.ErrorLevel},
		{"fatal", zap.FatalLevel},
		{"unknown", zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			atomicLevel := setLogLevel(tt.level)
			if atomicLevel.Level() != tt.expected {
				t.Fatalf("Expected %v, got %v", tt.expected, atomicLevel.Level())
			}
		})
	}
}
