# Spin Operator Development

- [To Deploy on the Cluster](#to-deploy-on-the-cluster)
- [To Uninstall](#to-uninstall)
- [Packaging and Deployment via Helm](#packaging-and-deployment-via-helm)

## To Deploy on the Cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/spin-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Apply the Runtime Class to the cluster:**

```sh
k apply -f config/samples/runtime-class.yaml
```

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy cert-manager to the cluster:**

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml
```

> **NOTE**: Cert-manager is required to manage the TLS certificates for the admission webhooks.

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/spin-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

> **NOTE**: Ensure that the samples has default values to test it out.

## To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**Delete the Runtime Class from the cluster:**

```sh
k delete -f config/samples/runtime-class.yaml
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

**UnDeploy cert-manager from the cluster:**

```sh
kubectl remove -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml
```

### Packaging and deployment via Helm

The [Spin Operator chart](./charts/spin-operator) is assembled via a combination of
[helmify](https://github.com/arttor/helmify) using the kustomize manifests from the
[config](./config/) directory as well as other non-kustomize items such as the
[NOTES.txt](./charts/spin-operator/templates/NOTES.txt) and [Chart.yaml](./charts/spin-operator/Chart.yaml).

> **NOTE**: Manual changes to helmify-generated resources, including the
> [values.yml](./charts/spin-operator/values.yaml) file and applicable resources in
> [templates](./charts/spin-operator/templates/) are not persisted across helmify
> invocations.

**Generate the Helm chart:**

```sh
make helm-generate
```

**Install the Helm chart onto the cluster:**

> **Note**: [CRDs](./config/crd/bases/) and the [wasm-spin-v2](./config/samples/runtime-class.yaml)
> RuntimeClass are currently _not_ installed as part of the chart. You'll need to ensure these are
> present via the method(s) mentioned above.

```sh
make helm-install
```

Follow the release notes printed after helm installs the chart for next steps.

**Upgrade the Helm release on the cluster:**

```sh
make helm-upgrade
```

**Delete the Helm release from the cluster:**

```sh
make helm-uninstall
```
