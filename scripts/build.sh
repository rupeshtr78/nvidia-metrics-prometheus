#!/bin/bash

TARGETOS=linux
TARGETARCH=amd64
 


GOOS=${TARGETOS:-linux}  GOARCH=${TARGETARCH:-amd64} go build -a -o bin/nvidia-metrics ./cmd/main.go