#!/bin/sh
set -e

docker load < k3s-airgap-images-amd64.tar

cp k3s /usr/local/bin

chmod +x install.sh /usr/local/bin/k3s

INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_BIN_DIR=/usr/local/bin ./install.sh