- [Running Locally](#running-locally)
  - [Prerequisites](#prerequisites)
  - [Fetch Spin Operator (Source Code)](#fetch-spin-operator-source-code)
  - [Setting Up Kubernetes Cluster](#setting-up-kubernetes-cluster)
  - [Installation With Make](#installation-with-make)
  - [Running the Sample Application](#running-the-sample-application)

# Running Locally

## Prerequisites

Please ensure that your system has all of the [prerequisites](./prerequisites.md) installed before continuing.

## Fetch Spin Operator (Source Code)

Clone the Spin Operator repository:

```bash
git clone https://github.com/spinkube/spin-operator.git
```

Change into the Spin Operator directory:

```bash
cd spin-operator
```

## Setting Up Kubernetes Cluster

Run the following command to create a Kubernetes k3d cluster that has [the containerd-wasm-shims](https://github.com/deislabs/containerd-wasm-shims) pre-requisites installed:

```bash
k3d cluster create wasm-cluster --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.10.0 -p "8081:80@loadbalancer" --agents 2
```

Run the following command to create the Runtime Class:

```bash
kubectl apply -f - <<EOF
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: wasmtime-spin-v2
handler: spin
EOF
```

## Installation With Make

Run the following command to install the Custom Resource Definition (CRD) into the cluster:

```bash
make install
```

## Running the Sample Application

Run the following command to run the Spin Operator locally:

```bash
make run
```

Run the following command, in a different terminal window:

```bash
$ kubectl apply -f ./config/samples/simple.yaml --validate=false
```

Run the following command to obtain the name of the pod you have running:

```bash
kubectl get pods
```

The above command will return information similar to the following:

```bash
NAME                              READY   STATUS    RESTARTS   AGE
simple-spinapp-5b8d8d69b4-snscs   1/1     Running   0          3m40s

```

Using the `NAME` from above, run the following `kubectl` command to listen on port 8083 locally and forward to port 80 in the pod:

```bash
kubectl port-forward simple-spinapp-5b8d8d69b4-snscs 8083:80
```

The above command will return the following forwarding mappings:

```bash
Forwarding from 127.0.0.1:8083 -> 80
Forwarding from [::1]:8083 -> 80
```

Run the following command, in a different terminal window:

```bash
curl localhost:8083/hello
```

The above command will return the following message:

```bash
Hello world from Spin!
```
