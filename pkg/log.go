// Description: This file contains the log functions for the package.
// Using uber-go/zap for logging.
// In this package we initate the logger.
// go get -u go.uber.org/zap
// go get -u go.uber.org/zap/zapcore
package pkg

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger


func Logger() (logger *zap.Logger, err error) {
 config := zap.NewProductionConfig()
 config.EncoderConfig.TimeKey = "timestamp"
 config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

 logger, err = config.Build()

 if err != nil {
  err = fmt.Errorf("Failed to initialize logger: %v", err)
  os.Exit(1)
  return nil, err
 }

 return logger, nil

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
