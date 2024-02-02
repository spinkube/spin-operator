- [Running On a Cluster](#running-on-a-cluster)
  - [Prerequisites](#prerequisites)
    - [Set up Your Kubernetes Cluster](#set-up-your-kubernetes-cluster)
  - [Install Spin Operator](#install-spin-operator)
  - [Using an Example Application](#using-an-example-application)
  - [Building the Example Application](#building-the-example-application)

# Running On a Cluster

## Prerequisites

Ensure necessary [prerequisites](./prerequisites.md) are installed.

For this Quickstart in particular, you will need:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [k3d](./prerequisites.md#k3d) - a lightweight Kubernetes distribution that runs on Docker
- [Docker](./prerequisites.md#docker) - for running k3d
- [Helm](./prerequisites.md#helm) - the package manager for Kubernetes

<!-- NOTE: remove this prerequisite when the runtime-class and CRDs can be applied from their release artifacts, i.e. when repo and release are public -->

Also, ensure you have cloned this repository and have navigated to the root of the project:

```
git clone git@github.com:spinkube/spin-operator.git
cd spin-operator
```

### Set up Your Kubernetes Cluster

1. Create a Kubernetes cluster with a k3d image that includes the [containerd-shim-spin](https://github.com/spinkube/containerd-shim-spin) prerequisite already installed:

<!-- TODO: update below with ghcr.io/spinkube/containerd-shim-spin/examples/k3d:<tag> -->

```
k3d cluster create wasm-cluster \
  --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.10.0 \
  --port "8081:80@loadbalancer" \
  --agents 2
```

> Note: Spin Operator requires a few Kubernetes resources that are installed globally to the cluster. We create these directly through `kubectl` as a best practice, since their lifetimes are usually managed separately from a given Spin Operator installation.

2. Apply the [Runtime Class](../../spin-runtime-class.yaml) used for scheduling Spin apps onto nodes running the shim:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml' -->

```
kubectl apply -f spin-runtime-class.yaml
```

3. Apply the [Custom Resource Definitions](./glossary-of-terms.md#custom-resource-definition-crd) used by the Spin Operator:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->

```
make install
```

## Install Spin Operator

Now that your Kubernetes cluster is prepared, you can install the Spin Operator via its Helm chart and test running Spin Operator on a Kubernetes cluster. This is harder than [running Spin Operator on your local machine](./running-locally.md), but deploying Spin Operator into your cluster lets you test things like webhooks. 

> Note: The following `helm install` command automatically installs sub-charts such as [Cert Manager](https://github.com/cert-manager/cert-manager) (used by spin-operator's admission webhook system). There is more information available at [the chart dependencies section](./deploying-with-helm.md#chart-dependencies) on the [deploying with Helm page](./deploying-with-helm.md).

<!-- TODO: remove '--devel' flag once we have our first non-prerelease chart available, e.g. when v0.1.0 of this project is out -->

```
helm install spin-operator \
  --namespace spin-operator \
  --create-namespace \
  --devel \
  --wait \
  oci://ghcr.io/spinkube/spin-operator
```

This will create all of the Kubernetes resources required by Spin Operator under the Kubernetes namespace `spin-operator`. It may take a moment for the installation to complete as dependencies are installed and pods are spinning up.

## Using an Example Application

Change into the `spin-operator/apps/hello-world` directory:

```bash
cd apps/hello-world
```

## Building the Example Application

```bash
spin build
```

pushes the Spin application to the appropriate repository:

```sh
spin registry push ghrc.io/user-name/application-name:latest
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

To create instances of your solution, apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

> **NOTE**: Ensure that the samples has default values to test it out.
