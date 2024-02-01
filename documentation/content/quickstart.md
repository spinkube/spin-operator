- [Quickstart](#quickstart)

# Quickstart

TODO
Helm plus the `kubectl apply` for CRDs is recommended. Documentation coming soon.

## Prerequisites

Ensure necessary [prerequisites](./prerequisites.md) are installed.

### Setting Up Your Kubernetes Cluster

1. Create a Kubernetes k3d cluster that has containerd-wasm-shim pre-requistes installed:

```
k3d cluster create wasm-cluster --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.10.0 -p "8081:80@loadbalancer" --agents 2
```

2. Apply the Runtime Class:

```
kubectl apply -f spin-runtime-class.yaml
```

## Running the Sample Application

1. `make install` to install the SpinApp CRD on to the cluster
2. `make run` to build and run the controller locally
3. In a different terminal window: `kubectl apply -f config/samples/shim-executor.yaml`
4. `kubectl apply -f config/samples/simple.yaml`
5. `kubectl port-forward svc/simple-spinapp 8083:80`
6. In a different terminal window: `curl localhost:8083/hello`

You should see:

```bash
Hello world from Spin!
```

If you want to test the admission webhooks you'll need to follow the instructions [here](operator_development.md) to deploy the operator to the cluster. We disable webhooks when using `make run` because that would require us to locally setup TLS certificates.
