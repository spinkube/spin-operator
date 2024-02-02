- [Deploying with Helm](#deploying-with-helm)
  - [Prerequisites](#prerequisites)
  - [Install Spin Operator Using Helm](#install-spin-operator-using-helm)
    - [Prepare the Cluster](#prepare-the-cluster)
  - [Chart dependencies](#chart-dependencies)
    - [Installing the Chart](#installing-the-chart)
    - [Upgrading the Chart](#upgrading-the-chart)
    - [Uninstalling the Chart](#uninstalling-the-chart)

# Deploying with Helm

## Prerequisites

Please ensure that your system has all of the [./prerequisites.md](prerequisites) installed before continuing.

## Install Spin Operator Using Helm

The following instructions are for installing Spin Operator as a chart (using helm install).

### Prepare the Cluster

Before installing the chart, you'll need to ensure the following:

The [Custom Resource Definition (CRD)](glossary-of-terms#custom-resource-definition-crd) resources are installed. This includes the SpinApp CRD representing Spin applications to be scheduled on the cluster.

<!-- TODO: templatize with release version corresponding to chart -->

```console
$ kubectl apply -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml
```

A [RuntimeClass](glossary-of-terms/#runtime-class) resource for the `wasmtime-spin-v2` container runtime is installed. This is the runtime that Spin applications use.

<!-- TODO: point to GH release artifact and templatize with release version corresponding to chart? Or static code link? -->

```console
$ kubectl apply -f - <<EOF
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: wasmtime-spin-v2
handler: spin
EOF
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

<!-- TODO: templatize with release version corresponding to chart -->

```console
$ helm install spin-operator --namespace spin-operator oci://ghcr.io/spinkube/spin-operator
```

### Upgrading the Chart

Note that you may also need to upgrade the spin-operator CRDs in tandem with upgrading the Helm release:

```console
$ kubectl apply -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml
```

To upgrade the `spin-operator` release, run the following:

```console
$ helm upgrade spin-operator --namespace spin-operator oci://ghcr.io/spinkube/spin-operator
```

### Uninstalling the Chart

To delete the `spin-operator` release, run:

```console
$ helm delete spin-operator --namespace spin-operator
```

This will remove all Kubernetes resources associated with the chart and deletes the Helm release.

To completely uninstall all resources related to spin-operator, you may want to delete the corresponding CRD resources and, optionally, the RuntimeClass:

```console
$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml

$ kubectl delete runtimeclass wasmtime-spin-v2
```

<!-- TODO: list out configuration options? -->
