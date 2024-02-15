- [Deploying with Helm](#deploying-with-helm)
  - [Prerequisites](#prerequisites)
    - [Install Helm](#install-helm)
  - [Install Spin Operator Using Helm](#install-spin-operator-using-helm)
    - [Prepare the Cluster](#prepare-the-cluster)
    - [Installing the Chart](#installing-the-chart)
    - [Upgrading the Chart](#upgrading-the-chart)
    - [Uninstalling the Chart](#uninstalling-the-chart)

# Deploying with Helm

## Prerequisites

Please ensure that your system has all of the [./prerequisites.md](prerequisites) installed before continuing.

For this guide in particular, you will need:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [Helm](./prerequisites.md#helm) - the package manager for Kubernetes

<!-- NOTE: remove this prerequisite when the runtime-class and CRDs can be applied from their release artifacts, i.e. when repo and release are public -->

Also, ensure you have cloned this repository and have navigated to the root of the project:

```console
git clone git@github.com:spinkube/spin-operator.git
cd spin-operator
```

## Install Spin Operator Using Helm

The following instructions are for installing Spin Operator as a chart (using helm install).

### Prepare the Cluster

Before installing the chart, you'll need to ensure the following:

The [Custom Resource Definition (CRD)](glossary-of-terms#custom-resource-definition-crd) resources are installed. This includes the SpinApp CRD representing Spin applications to be scheduled on the cluster.

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->

```console
make install
```

A [RuntimeClass](glossary-of-terms/#runtime-class) resource for the `wasmtime-spin-v2` container runtime is installed. This is the runtime that Spin applications use.

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml' -->

```console
kubectl apply -f spin-runtime-class.yaml
```

## Chart dependencies

The spin-operator chart currently includes the following sub-charts:

- [Kwasm Operator](https://github.com/kwasm/kwasm-operator) to install WebAssembly support on Kubernetes nodes
- [Cert Manager](https://github.com/cert-manager/cert-manager) to automatically provision and manage TLS certificates (used by spin-operator's admission webhook system)
  - If you'd like to manage Cert Manager completely separate from spin-operator, you can disable installation via:
    `--set certmanager.enabled=false` on `helm install`.
  - Or, if you'd like to install Cert Manager separate from its CRDs, you can opt-out of installing the CRDs via:
    `--set certmanager.installCRDs=false` on `helm install`.
  - In either case, Cert Manager must be running and the corresponding CRDs must be present on the cluster before installing the spin-operator chart.

### Installing the Chart

The following installs the chart with the release name `spin-operator`:

<!-- TODO: remove '--devel' flag once we have our first non-prerelease chart available, e.g. when v0.1.0 of this project is released and public -->

```console
helm install spin-operator \
  --namespace spin-operator-system \
  --create-namespace \
  --devel \
  --wait \
  charts/spin-operator
```

### Upgrading the Chart

Note that you may also need to upgrade the spin-operator CRDs in tandem with upgrading the Helm release:

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->

```
make install
```

To upgrade the `spin-operator` release, run the following:

<!-- TODO: remove '--devel' flag once we have our first non-prerelease chart available, e.g. when v0.1.0 of this project is released and public -->

```console
helm upgrade spin-operator \
  --namespace spin-operator-system \
  --devel \
  --wait \
  charts/spin-operator
```

### Uninstalling the Chart

To delete the `spin-operator` release, run:

```console
helm delete spin-operator --namespace spin-operator-system
```

This will remove all Kubernetes resources associated with the chart and deletes the Helm release.

To completely uninstall all resources related to spin-operator, you may want to delete the corresponding CRD resources and, optionally, the RuntimeClass:

<!-- TODO: replace with:
```console
kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml

kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml
```
-->

```console
make uninstall
kubectl delete -f spin-runtime-class.yaml
```

<!-- TODO: list out configuration options? -->
