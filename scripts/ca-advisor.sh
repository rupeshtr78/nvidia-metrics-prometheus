#!/bin/sh

VERSION=v0.49.1 # use the latest release version from https://github.com/google/cadvisor/releases
sudo docker run \
  --volume=/:/rootfs:ro \
  --volume=/var/run:/var/run:ro \
  --volume=/sys:/sys:ro \
  --volume=/var/lib/docker/:/var/lib/docker:ro \
  --volume=/dev/disk/:/dev/disk:ro \
  --publish=9600:8080 \
  --detach=true \
  --name=cadvisor \
  --privileged \
  --device=/dev/kmsg \
  gcr.io/cadvisor/cadvisor:$VERSION


  # # Prometheus
  # - job_name: 'cadvisor'
  #   scrape_interval: 10s 
  #   metrics_path: '/metrics'
  #   static_configs:
  #     - targets: ['localhost:9600']
  #       labels:
  #         group: 'cadvisor'