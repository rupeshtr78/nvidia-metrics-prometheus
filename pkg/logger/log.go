// Description: This file contains the log functions for the package.
// Using uber-go/zap for logging.
// In this package we initate the logger.
// go get -u go.uber.org/zap
// go get -u go.uber.org/zap/zapcore
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger is the global logger instance.
var logger *zap.Logger

func init() {
	level := "info"
	err := GetLogger(level)

	if err != nil {
		fmt.Println(err)
	}

}

func GetLogger(level string) (err error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.InitialFields = map[string]interface{}{
		"app": "nvidia-metrics",
	}
	config.Level = setLogLevel(level)

	// Enable caller information reporting.
	// Modify the EncoderConfig for the caller key and format
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // Use ShortCallerEncoder to log the relative path
	config.EncoderConfig.CallerKey = "caller"                      // Specify the key used for the caller in structured logs

	logger, err = config.Build(zap.AddCallerSkip(1)) // Skip one level to account for this wrapper.

	if err != nil {
		err = fmt.Errorf("Failed to initialize logger: %v", err)
		return err
	}

	return nil

}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

func setLogLevel(level string) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()
	switch level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	return atomicLevel

}
