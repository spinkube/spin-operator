- [Quickstart](#quickstart)
  - [Prerequisites](#prerequisites)
  - [Set up your Kubernetes cluster](#set-up-your-kubernetes-cluster)
  - [Deploy the Spin Operator](#deploy-the-spin-operator)
  - [Run the sample application](#run-the-sample-application)

# Quickstart

This Quickstart guide demonstrates how to set up a new Kubernetes cluster, install the Spin Operator and deploy your first Spin application.

## Prerequisites

Ensure necessary [prerequisites](./prerequisites.md) are installed.

For this Quickstart in particular, you will need:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [k3d](./prerequisites.md#k3d) - a lightweight Kubernetes distribution that runs on Docker
- [Docker](./prerequisites.md#docker) - for running k3d

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

> > For now in private preview, the installation workflow uses `make`. In the future, we will add Helm chart support

2. Build the Spin Operator image.

```console
make docker-build IMG=ghcr.io/spinkube/spin-operator:dev
```

3. Import Spin Operator to your k3d cluster.

```console
k3d image import -c wasm-cluster ghcr.io/spinkube/spin-operator:dev
```

4. Apply the [Runtime Class](../../spin-runtime-class.yaml) used for scheduling Spin apps onto nodes running the shim:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml' -->

```console
kubectl apply -f spin-runtime-class.yaml
```

5. Apply the [Custom Resource Definitions](./glossary-of-terms.md#custom-resource-definition-crd) used by the Spin Operator:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->

```console
make install
```

6. Install cert manager

```console
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.2/cert-manager.yaml
```

## Deploy the Spin Operator

Run the following command to run the Spin Operator locally. This will create all of the Kubernetes resources required by Spin Operator under the Kubernetes namespace spin-operator. It may take a moment for the installation to complete as dependencies are installed and pods are spinning up.

```console
make deploy IMG=ghcr.io/spinkube/spin-operator:dev
```

Lastly, create the shim executor:

```console
kubectl apply -f config/samples/shim-executor.yaml
```

## Run the Sample Application

You are now ready to deploy Spin applications onto the cluster!

<!-- TODO: if/when we have the option and if we wanted to, we could mention that the kwasm operator isn't needed when using k3d, as the containerd-shim-spin is already present. Installation could be skipped via --set kwasm-operator.enabled=false -->

1. Create your first application in the same `spin-operator` namespace that the operator is running:

<!-- Note: the default 'containerd-shim-spin' SpinAppExecutor CR needs to be present on the cluster before apps using this default can run. However, as of writing, it is a namespaced resource. As such, apps can only be deployed in the same namespace(s) that the CR is present. -->

```console
kubectl apply -f config/samples/simple.yaml
```

2. Forward a local port to the application pod so that it can be reached:

```console
kubectl port-forward svc/simple-spinapp 8083:80
```

3. In a different terminal window, make a request to the application:

```console
curl localhost:8083/hello
```

You should see:

```bash
Hello world from Spin!
```

## Next Steps

Congrats on deploying your first SpinApp! Recommended next steps:

- Scale your [Spin Apps with Horizontal Pod Autoscaler](./scaling-spinapp-on-k8s-with-hpa.md)
- Scale your [Spin Apps with Kubernetes Event Driven Autoscaler](./scaling-spinapp-on-k8s-with-keda.md)
