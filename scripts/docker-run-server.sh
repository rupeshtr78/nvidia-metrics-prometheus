#!/bin/bash
docker run -p 9500:9500 --gpus all nvidia-gpu-exporter:1.0

docker run -p 9500:9500 --gpus all nvidia-gpu-exporter:1.1

docker run -p 9500:9500 --gpus all -e HOST=0.0.0.0 nvidia-gpu-exporter:1.0


docker run -p 9500:9500 --gpus all -e CONFIG_FILE=/path/to/config.yaml \
           -e LOG_LEVEL=debug \
           -e PORT=9500 \
           -e HOST=0.0.0.0 \
           -e INTERVAL=5 \
           nvidia-gpu-exporter:1.0

# Verify that the port is open and accessible on the host machine
sudo ufw allow 9500/tcp
sudo ufw reload

# Verify that the port is open
sudo netstat -tuln | grep 9500