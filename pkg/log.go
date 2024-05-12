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

func init() {
 config := zap.NewProductionConfig()
 config.EncoderConfig.TimeKey = "timestamp"
 config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
 var err error
 logger, err = config.Build()
 if err != nil {
  fmt.Printf("Error initializing logger: %v\n", err)
  os.Exit(1)
 }
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
