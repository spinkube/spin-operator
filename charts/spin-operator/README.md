# spin-operator

spin-operator is a Kubernetes operator in charge of handling the lifecycle of Spin applications based on their SpinApp resources.

## Prerequisites

- Kubernetes v1.11.3+

## Prepare the cluster

Prior to installing the chart, you'll need to ensure the following:

- spin-operator CustomResourceDefinition (CRD) resources are installed. This includes the SpinApp CRD representing Spin applications to be scheduled on the cluster.

  <!-- TODO: templatize with release version corresponding to chart's appVersion -->

  ```console
  $ kubectl apply -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml
  ```

- A RuntimeClass resource for the `wasmtime-spin-v2` container runtime is installed. This is the runtime that Spin applications use.

  <!-- TODO: point to GH release artifact and templatize with release version corresponding to chart's appVersion? Or static code link? -->

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

## Installing the chart

The following installs the chart with the release name `spin-operator`:

<!-- TODO: templatize with release version corresponding to chart's appVersion -->

```console
$ helm install spin-operator --namespace spin-operator oci://ghcr.io/spinkube/spin-operator
```

## Upgrading the chart

Note that you may also need to upgrade the spin-operator CRDs in tandem with upgrading the Helm release:

```console
$ kubectl apply -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml
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

```console
$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/latest/download/spin-operator.crds.yaml

$ kubectl delete runtimeclass wasmtime-spin-v2
```

<!-- TODO: list out configuration options? -->
