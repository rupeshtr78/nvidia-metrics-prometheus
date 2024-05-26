// Description: This file contains the log functions for the package.
// Using uber-go/zap for logging.
// In this package we initate the logger.
// go get -u go.uber.org/zap
// go get -u go.uber.org/zap/zapcore
package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger is the global logger instance.
var logger *zap.Logger

// @TODO - remove after testing
func Loggerinit(level string, fileLog bool, filePath string) {
	err := GetLogger(level, fileLog, filePath)

	if err != nil {
		fmt.Println(err)
	}

}

func GetLogger(level string, fileLog bool, filePath string) (err error) {
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
	config.EncoderConfig.CallerKey = "caller"

	if fileLog {
		config.OutputPaths = []string{"stdout", filePath}
		config.ErrorOutputPaths = []string{"stderr", filePath}
	} else {
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
	}

	logger, err = config.Build(zap.AddCallerSkip(1)) // Skip one level to account for this wrapper.

	if err != nil {
		err = fmt.Errorf("failed to initialize logger: %v", err)
		return err
	}

	return nil

}

// // Configure the logger to write to the file
func configureZapCore(fileName string) zapcore.Core {
	// Create a file for logging
	file, err := createLogFile(fileName)
	if err != nil {
		fmt.Println("Failed to create log file: ", err)
	}

	// Configure the encoder
	encoderCfg := zap.NewProductionEncoderConfig()

	// Create a new JSON encoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// Create a new core
	writeSyncer := zapcore.AddSync(file)

	// Create a new core
	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)

	return core

}

func createLogFile(fileName string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(file), nil
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
