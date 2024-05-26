#!/bin/sh


# https://catalog.ngc.nvidia.com/orgs/nvidia/teams/k8s/containers/container-toolkit
# The NVIDIA Container Toolkit enables users to build and run GPU-accelerated containers. 
# The toolkit includes a container runtime library and utilities to automatically configure containers to leverage NVIDIA GPUs.

# Add the package repositories
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
  && curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
    sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
    sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

# Update the packages list from the repository:
sudo apt-get update
# Install the NVIDIA Container Toolkit packages:
sudo apt-get install -y nvidia-container-toolkit

# Configuring Docker
# Configure the container runtime by using the nvidia-ctk command:

sudo nvidia-ctk runtime configure --runtime=docker
# The nvidia-ctk command modifies the /etc/docker/daemon.json file on the host. 
# The file is updated so that Docker can use the NVIDIA Container Runtime.

# Restart the Docker daemon:

sudo systemctl restart docker

# Comment this line if you want to run the container with the default runtime.
docker run --gpus all --rm nvcr.io/nvidia/k8s/container-toolkit:v1.15.0-ubi8 nvidia-smi
docker run --runtime=nvidia --rm nvcr.io/nvidia/k8s/container-toolkit:v1.15.0-ubi8 nvidia-smi

# Errors
# docker: Error response from daemon: could not select device driver "" with capabilities: [[gpu]].

docker configuration: Check if /etc/docker/daemon.json is configured correctly to use nvidia as the default runtime, it should look something like this:
{
  "default-runtime": "nvidia",
  "runtimes": {
    "nvidia": {
      "path": "/usr/bin/nvidia-container-runtime",
      "runtimeArgs": []
    }
  }
}

# current daemon.json
{
    "runtimes": {
        "nvidia": {
            "args": [],
            "path": "nvidia-container-runtime"
        }
    }
}

# msg="error running nvidia-toolkit: unable to initialize: unable to create pidfile: open /run/nvidia/toolkit.pid: no such file or directory"
# Edit your runtime configuration under /etc/nvidia-container-runtime/config.toml and uncomment the debug=... line.
# Run your container again to reproduce the issue and generate the logs.

# Check if the /run/nvidia directory exists and create it if it doesn't:
sudo mkdir -p /run/nvidia
sudo touch /run/nvidia/toolkit.pid
sudo chmod 0755 /run/nvidia
sudo chmod 777 /run/nvidia/toolkit.pid
docker run --runtime=nvidia --rm nvcr.io/nvidia/k8s/container-toolkit:v1.15.0-ubi8 nvidia-smi


# docker-compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version
docker-compose up -d > /dev/null 2>&1
docker-compose up -d > logs/nvidia-metrics.log 2>&1
docker-compose ps
