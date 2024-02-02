- [Uninstall](#uninstall)
  - [Uninstalling the Helm Chart](#uninstalling-the-helm-chart)
  - [Delete (CRs)](#delete-crs)
  - [Delete APIs(CRDs)](#delete-apiscrds)
  - [UnDeploy](#undeploy)

# Uninstall

## Uninstalling the Helm Chart

If you [installed Spin Operator using Helm](./deploying-with-helm.md
#install-spin-operator-using-helm), the following steps will uninstall the Helm chart.

To delete the `spin-operator` release, run:

```console
$ helm delete spin-operator --namespace spin-operator
```

This will remove all Kubernetes resources associated with the chart and delete the Helm release.

To completely uninstall all resources related to spin-operator, you may want to delete the corresponding CRD resources and, optionally, the RuntimeClass:

<!-- TODO: templatize with release version corresponding to chart's appVersion -->

```console
$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml

$ kubectl delete -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml
```

## Delete (CRs)

These are commands to delete un-deploy resources. The following command will delete the instances (CRs) from the cluster:

```bash
kubectl delete -k config/samples/
```

## Delete APIs(CRDs)

The following command will uninstall CRDs from the K8s cluster specified in `~/.kube/config`:

```sh
make uninstall
```

> Call with `ignore-not-found=true` to ignore resource not found errors during deletion.

## UnDeploy

The following command will undeploy the controller from the K8s cluster specified in `~/.kube/config`:

```sh
make undeploy
```

> Call with `ignore-not-found=true` to ignore resource not found errors during deletion.
