- [Uninstall](#uninstall)
  - [Delete (CRs)](#delete-crs)
  - [Delete APIs(CRDs)](#delete-apiscrds)
  - [UnDeploy](#undeploy)

# Uninstall

These are commands to delete, uninstall and undeploy resources.

## Delete (CRs)

The following command will delete the instances (CRs) from the cluster:

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
