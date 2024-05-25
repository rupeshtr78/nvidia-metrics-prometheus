#!/bin/bash

docker run -e CONFIG_FILE=/path/to/config.yaml \
           -e LOG_LEVEL=debug \
           -e PORT=9500 \
           -e HOST=0.0.0.0 \
           -e INTERVAL=5 \
           nvidia-gpu-exporter:1.0