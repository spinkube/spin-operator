# spin-operator

spin-operator is a Kubernetes operator in charge of handling the lifecycle of Spin applications based on their SpinApp resources.

## Prerequisites

- Kubernetes v1.11.3+

## Prepare the cluster

Prior to installing the chart, you'll need to ensure the following:

- [Cert Manager](https://github.com/cert-manager/cert-manager) to automatically provision and manage TLS certificates (used by spin-operator's admission webhook system). Cert Manager must be running and the corresponding CRDs must be present on the cluster before installing the spin-operator chart.

- [Kwasm Operator](https://github.com/kwasm/kwasm-operator) to install WebAssembly support on Kubernetes nodes. See the [project README.md](https://github.com/KWasm/kwasm-operator/blob/main/README.md) for installation and configuration steps, including annotating nodes to run Spin/wasm workloads.

- spin-operator CustomResourceDefinition (CRD) resources are installed. This includes the SpinApp CRD representing Spin applications to be scheduled on the cluster.

  ```console
  $ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinapps.yaml
  $ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinappexecutors.yaml
  ```

## Installing the chart

The following installs the chart with the release name `spin-operator`:

```console
$ helm install spin-operator \
  --namespace spin-operator \
  --create-namespace \
  --version {{ CHART_VERSION }} \
  oci://ghcr.io/spinkube/spin-operator
```

## Post-installation

spin-operator depends on the following resources. If not already present on the cluster, install them now:

- An application executor is installed. This is the executor that spin-operator uses to run Spin applications.

  ```console
  $ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/samples/spin-shim-executor.yaml
  ```

- A RuntimeClass resource for the `wasmtime-spin-v2` container runtime is installed. This is the runtime that Spin applications use.

  ```console
  $ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/samples/spin-runtime-class.yaml
  ```

## Upgrading the chart

Note that you may also need to upgrade the spin-operator CRDs in tandem with upgrading the Helm release:

```console
$ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinapps.yaml
$ kubectl apply -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinappexecutors.yaml
```

To upgrade the `spin-operator` release, run the following:

```console
$ helm upgrade spin-operator \
  --namespace spin-operator \
  --version {{ CHART_VERSION }} \
  oci://ghcr.io/spinkube/spin-operator
```

## Uninstalling the chart

To delete the `spin-operator` release, run:

```console
$ helm delete spin-operator --namespace spin-operator
```

This will remove all Kubernetes resources associated with the chart and deletes the Helm release.

To completely uninstall all resources related to spin-operator, you may want to delete the corresponding CRD resources and, optionally, the RuntimeClass:

```console
$ kubectl delete -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/samples/spin-runtime-class.yaml
$ kubectl delete -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/samples/spin-shim-executor.yaml
$ kubectl delete -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinapps.yaml
$ kubectl delete -f https://raw.githubusercontent.com/spinkube/spin-operator/main/config/crd/bases/core.spinoperator.dev_spinappexecutors.yaml
```
