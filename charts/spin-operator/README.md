# spin-operator

spin-operator is a Kubernetes operator in charge of handling the lifecycle of Spin applications based on their SpinApp resources.

## Prerequisites

- Kubernetes v1.11.3+

## Prepare the cluster

Prior to installing the chart, you'll need to ensure the following:

- [Cert Manager](https://github.com/cert-manager/cert-manager) to automatically provision and manage TLS certificates (used by spin-operator's admission webhook system). Cert Manager must be running and the corresponding CRDs must be present on the cluster before installing the spin-operator chart.

- spin-operator CustomResourceDefinition (CRD) resources are installed. This includes the SpinApp CRD representing Spin applications to be scheduled on the cluster.

  <!-- TODO: templatize with release version corresponding to chart's appVersion -->

  ```console
  $ kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml
  ```

## Chart dependencies

The spin-operator chart currently includes the following sub-charts:

- [Kwasm Operator](https://github.com/kwasm/kwasm-operator) to install WebAssembly support on Kubernetes nodes

## Installing the chart

The following installs the chart with the release name `spin-operator`:

<!-- TODO: templatize with release version corresponding to chart's appVersion -->

```console
$ helm install spin-operator --namespace spin-operator --create-namespace oci://ghcr.io/spinkube/spin-operator
```

## Post-installation

After installing the chart, you'll need to ensure the following:

- An application executor is installed. This is the executor that spin-operator uses to run Spin applications.

  <!-- TODO: templatize with release version corresponding to chart's appVersion -->

  ```console
  $ kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.shim-executor.yaml
  ```

- A RuntimeClass resource for the `wasmtime-spin-v2` container runtime is installed. This is the runtime that Spin applications use.

  <!-- TODO: templatize with release version corresponding to chart's appVersion -->

  ```console
  $ kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml
  ```

## Upgrading the chart

Note that you may also need to upgrade the spin-operator CRDs in tandem with upgrading the Helm release:

<!-- TODO: templatize with release version corresponding to chart's appVersion -->

```console
$ kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml
```

To upgrade the `spin-operator` release, run the following:

```console
$ helm upgrade spin-operator --namespace spin-operator oci://ghcr.io/spinkube/spin-operator
```

## Uninstalling the chart

To delete the `spin-operator` release, run:

```console
$ helm delete spin-operator --namespace spin-operator
```

This will remove all Kubernetes resources associated with the chart and deletes the Helm release.

To completely uninstall all resources related to spin-operator, you may want to delete the corresponding CRD resources and, optionally, the RuntimeClass:

<!-- TODO: templatize with release version corresponding to chart's appVersion -->

```console
$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml

$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml
```

<!-- TODO: list out configuration options? -->
