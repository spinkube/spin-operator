- [Custom Resource Definition (CRD) Reference](#custom-resource-definition-crd-reference)
  - [Runtime (Deprecated)](#runtime-deprecated)
  - [Executor (Optional)](#executor-optional)
  - [Scheduler](#scheduler)
  - [Image (Required)](#image-required)
  - [ImagePullSecrets (Optional)](#imagepullsecrets-optional)
  - [Replicas (Required)](#replicas-required)
  - [RuntimeConfig (Optional)](#runtimeconfig-optional)
  - [ServiceAnnotations (Optional)](#serviceannotations-optional)
  - [DeploymentAnnotations (Optional)](#deploymentannotations-optional)
  - [PodAnnotations (Optional)](#podannotations-optional)
  - [ResourceRequirements](#resourcerequirements)
  - [VolumeMounts](#volumemounts)
  - [Status](#status)

# Custom Resource Definition (CRD) Reference

## Runtime (Deprecated)

## Executor (Optional)

Executor configures what will execute the `SpinApp`. Currently we support `containerd-shim-spin`. However, this field is configurable so Spin Operator can manage other runtimes in the future.

If you choose `containerd-shim-spin` (it is also the default) the operator will create a deployment and a service. The service will point at pods managed by the deployment. The deployment will start and manage the pods with `containerd-shim-spin` for you. It is expected that you already have the `wasmtime-spin-v2` runtime class installed on the cluster to use this executor.

See [this sample application](https://github.com/spinkube/spin-operator/blob/main/config/samples/cyclotron.yaml).

## Scheduler

The Spin Operator Scheduler is currently operating based on one instance per scheduler. For now, Spin Operator is treated as being necessary to control SpinApps. Spin Operator can manage a Service regardless of the scheduler.

We have the option of making a future feature whereby the Service management is configurable on a per-scheduler basis. This allows external schedulers to define their own.

## Image (Required)

Points to the image of the Spin app you want to run. For example:

```yaml
image: "ghcr.io/deislabs/containerd-wasm-shims/examples/spin-rust-hello:v0.10.0"
```

See [this sample application](https://github.com/spinkube/spin-operator/blob/main/config/samples/simple.yaml)

## ImagePullSecrets (Optional)

In some cases, your image might be coming from a private registry. Lets you reference a k8s secret that has credentials for you to pull an image.

For example, a secret which is created with the following command:

```bash
kubectl create secret docker-registry spin-image-secret --docker-server=https://ghcr.io --docker-username=$YOUR_GITHUB_USERNAME --docker-password=$YOUR_GITHUB_PERSONAL_ACCESS_TOKEN --docker-email=$YOUR_EMAIL
```

See [this sample application](https://github.com/spinkube/spin-operator/blob/main/config/samples/private-image.yaml)

## Replicas (Required)

Replicas is a field in the `SpinApp` Custom Resource Definition (CRD). Configures how many replicas of a spin app you want to run. If containerd-shim-spin is the executor that is the number of pods. This definition may very depending on how other executors choose to define replica.

## RuntimeConfig (Optional)

Lets you define Spin runtime config for your app. You must base64 encode your runtime config and put it in a secret with the right key. See the sample app for an example. It will then put the runtime config as a volume mount in your pod in the right place for the shim to pick it up and use it.

Converts a runtime-config.toml file to a Kubernetes secret `runtime-config-to-secret [PATH_TO_RUNTIME_CONFIG] [SECRET_NAME]`

See [this sample application](https://github.com/spinkube/spin-operator/blob/main/config/samples/runtime-config.yaml)

## ServiceAnnotations (Optional)

Passing annotations through to the deployment is supported. Lets you set specific annotations on the underlying service that is created. For example:

```yaml
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: annotations-spinapp
spec:
  image: "ghcr.io/deislabs/containerd-wasm-shims/examples/spin-rust-hello:v0.10.0"
  replicas: 1
  serviceAnnotations:
    key: value
  deploymentAnnotations:
    key: value
    multiple-keys: are-supported
  podAnnotations:
    key: value
```

See [this example application](https://github.com/spinkube/spin-operator/blob/main/config/samples/annotations.yaml)

## DeploymentAnnotations (Optional)

Lets you set specific annotations on the underlying deployment that is created.

See [this example application](https://github.com/spinkube/spin-operator/blob/main/config/samples/annotations.yaml)

## PodAnnotations (Optional)

Lets you set specific annotations on the underlying pods that are created.

See [this example application](https://github.com/spinkube/spin-operator/blob/main/config/samples/annotations.yaml)

## ResourceRequirements

Still being developed TODO

## VolumeMounts

Still being developed TODO

## Status

Still being designed TODO

The status field indicates the state of a `SpinApp` in the Kubernetes cluster which includes the state of the resources the `SpinApp` controller has created. Below are examples of status:

```
TODO
```
