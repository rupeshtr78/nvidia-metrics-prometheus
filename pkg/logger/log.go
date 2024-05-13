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
	err := GetLogger()

	if err != nil {
		fmt.Println(err)
	}

	setLogLevel("debug")
}

func GetLogger() (err error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = config.Build()

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

func setLogLevel(level string) {
	switch level {
	case "debug":
		logger.Info("Setting log level to debug")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.DebugLevel)
		logger, _ = config.Build()
	case "info":
		logger.Info("Setting log level to info")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.InfoLevel)
		logger, _ = config.Build()
	case "warn":
		logger.Info("Setting log level to warn")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.WarnLevel)
		logger, _ = config.Build()
	case "error":
		logger.Info("Setting log level to error")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.ErrorLevel)
		logger, _ = config.Build()
	case "fatal":
		logger.Info("Setting log level to fatal")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.FatalLevel)
		logger, _ = config.Build()
	default:
		logger.Info("Setting log level to info")
		config := zap.NewProductionConfig()
		config.Level.SetLevel(zapcore.InfoLevel)
		logger, _ = config.Build()
	}
}
