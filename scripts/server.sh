#!/bin/bash

# Setting the environment variables
export CONFIG_FILE="config/metrics.yaml"
export LOG_LEVEL="debug"
export PORT="9500"
export HOST="0.0.0.0"
export INTERVAL="5"
export LOG_FILE_PATH="logs/gpu-metrics.log"
export LOG_TO_FILE="false"

# Running the server
go run cmd/main.go
