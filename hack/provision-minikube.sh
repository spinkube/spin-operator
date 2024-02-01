#!/usr/bin/env bash
set -euo pipefail

# Provision a new minikube instance configured with the containerd shim for Spin.

echo "Starting minikube"
minikube start --container-runtime containerd

echo "Installing the containerd shim"
curl -fsSLO https://github.com/deislabs/containerd-wasm-shims/releases/download/v0.10.0/containerd-wasm-shims-v2-spin-linux-aarch64.tar.gz
tar -zxvf containerd-wasm-shims-v2-spin-linux-aarch64.tar.gz

echo "Copying the shim to minikube"
minikube cp containerd-shim-spin-v2 /usr/local/bin/
minikube ssh sudo chmod +x /usr/local/bin/containerd-shim-spin-v2

# just cleaning up
rm containerd-wasm-shims-v2-spin-linux-aarch64.tar.gz containerd-shim-spin-v2

echo "Configuring containerd"
if ! minikube ssh -- grep -q io.containerd.spin /etc/containerd/config.toml; then
  echo nope
  minikube ssh 'cat << EOF | sudo tee -a /etc/containerd/config.toml
[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.spin]
  runtime_type = "io.containerd.spin.v2"
EOF'
fi

echo "Restarting containerd"
minikube ssh sudo systemctl restart containerd

echo "Creating runtime class"
kubectl apply -f - <<EOF
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: wasmtime-spin-v2
handler: spin
EOF
