- [Quickstart](#quickstart)
  - [Prerequisites](#prerequisites)
  - [Set up your Kubernetes cluster](#set-up-your-kubernetes-cluster)
  - [Install Spin Operator](#install-spin-operator)
  - [Run the sample application](#run-the-sample-application)

# Quickstart

This Quickstart guide demonstrates how to set up a new Kubernetes cluster, install the Spin Operator and deploy your first Spin application.

## Prerequisites

Ensure necessary [prerequisites](./prerequisites.md) are installed.

For this Quickstart in particular, you will need:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [k3d](./prerequisites.md#k3d) - a lightweight Kubernetes distribution that runs on Docker
- [Docker](./prerequisites.md#docker) - for running k3d
- [Helm](./prerequisites.md#helm) - the package manager for Kubernetes

<!-- NOTE: remove this prerequisite when the runtime-class and CRDs can be applied from their release artifacts, i.e. when repo and release are public -->

Also, ensure you have cloned this repository and have navigated to the root of the project:

```console
git clone git@github.com:spinkube/spin-operator.git
cd spin-operator
```

### Set up Your Kubernetes Cluster

1. Create a Kubernetes cluster with a k3d image that includes the [containerd-shim-spin](https://github.com/spinkube/containerd-shim-spin) prerequisite already installed:

<!-- TODO: update below with ghcr.io/spinkube/containerd-shim-spin/examples/k3d:<tag> -->

```console
k3d cluster create wasm-cluster \
  --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.11.0 \
  --port "8081:80@loadbalancer" \
  --agents 2
```

> Note: Spin Operator requires a few Kubernetes resources that are installed globally to the cluster. We create these directly through `kubectl` as a best practice, since their lifetimes are usually managed separately from a given Spin Operator installation.

2. Apply the [Runtime Class](../../spin-runtime-class.yaml) used for scheduling Spin apps onto nodes running the shim:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml' -->

```console
kubectl apply -f spin-runtime-class.yaml
```

3. Apply the [Custom Resource Definitions](./glossary-of-terms.md#custom-resource-definition-crd) used by the Spin Operator:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->

```console
make install
```

## Install Spin Operator

Now that your Kubernetes cluster is prepared, you can install the Spin Operator via its Helm chart:

<!-- TODO: remove '--devel' flag once we have our first non-prerelease chart available, e.g. when v0.1.0 of this project is out -->

```console
helm install spin-operator \
  --namespace spin-operator \
  --create-namespace \
  --devel \
  --wait \
  oci://ghcr.io/spinkube/spin-operator
```

This will create all of the Kubernetes resources required by Spin Operator under the Kubernetes namespace `spin-operator`. It may take a moment for the installation to complete as dependencies are installed and pods are spinning up.

## Run the Sample Application

You are now ready to deploy Spin applications onto the cluster!

<!-- TODO: if/when we have the option and if we wanted to, we could mention that the kwasm operator isn't needed when using k3d, as the containerd-shim-spin is already present. Installation could be skipped via --set kwasm-operator.enabled=false -->

1. Create your first application in the same `spin-operator` namespace that the operator is running:

<!-- Note: the default 'containerd-shim-spin' SpinAppExecutor CR needs to be present on the cluster before apps using this default can run. However, as of writing, it is a namespaced resource. As such, apps can only be deployed in the same namespace(s) that the CR is present. -->

```console
kubectl -n spin-operator apply -f config/samples/simple.yaml
```

<!-- TODO: Use spin-k8s-plugin here? -->

2. Forward a local port to the application pod so that it can be reached:

```console
kubectl -n spin-operator port-forward svc/simple-spinapp 8083:80
```

3. In a different terminal window, make a request to the application:

```console
curl localhost:8083/hello
```

You should see:

```console
Hello world from Spin!
```

<!-- TODO: guide the reader to the next relevant documentation section -->
