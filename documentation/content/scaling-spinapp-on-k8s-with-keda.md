- [Scaling Spinapp on Kubernetes (k8s) With Kubernetes Event-Driven Autoscaling (KEDA)](#scaling-spinapp-on-kubernetes-k8s-with-kubernetes-event-driven-autoscaling-keda)
  - [Prerequisites](#prerequisites)
  - [Fetch Spin Operator (Source Code)](#fetch-spin-operator-source-code)
  - [Setting Up k8s Cluster](#setting-up-k8s-cluster)
  - [Setting Up Ingress](#setting-up-ingress)
  - [Setting Up KEDA](#setting-up-keda)
  - [Build and Store SpinApp in a TTL Registry](#build-and-store-spinapp-in-a-ttl-registry)
  - [Deploy SpinApp and the KEDA ScaledObject](#deploy-spinapp-and-the-keda-scaledobject)
  - [Generate Load to Test Autoscale](#generate-load-to-test-autoscale)

# Scaling Spinapp on Kubernetes (k8s) With Kubernetes Event-Driven Autoscaling (KEDA)

[KEDA](https://keda.sh) extends Kubernetes to provide event-driven scaling capabilities, allowing it to react to events from Kubernetes internal and external sources using [KEDA scalers](https://keda.sh/docs/2.13/scalers/). KEDA provides a wide variety of scalers to define scaling behavior base on sources like CPU, Memory, Azure Event Hubs, Kafka, RabbitMQ, and more. We use a `ScaledObject` to dynamically scale the instance count of our SpinApp to meet the demand.

## Prerequisites

> We use k3d to run a k8s cluster locally as part of this tutorial, but you can follow these steps to configure KEDA autoscaling on your desired k8s environment.

Please see the following sections in the [Prerequisites](./prerequisites.md) page and fulfil those prerequisite requirements before continuing:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [k3d](./prerequisites.md#k3d) - a lightweight Kubernetes distribution that runs on Docker
- [Docker](./prerequisites.md#docker) - for running k3d
- [Helm](./prerequisites.md#helm) - the package manager for Kubernetes
- [Bombardier](#prerequisites#bombardier) - cross-platform HTTP benchmarking CLI

## Fetch Spin Operator (Source Code)

If you haven't already, please go ahead and clone the Spin Operator repository:

```console
git clone https://github.com/spinkube/spin-operator.git
```

Change into the Spin Operator directory:

```console
cd spin-operator
```

## Setting Up k8s Cluster

Run the following command to create a k8s k3d cluster that has [the containerd-wasm-shims](https://github.com/deislabs/containerd-wasm-shims) pre-requisites installed: If you have a k3d cluster already, please feel free to use it:

```console
k3d cluster create wasm-cluster-scale --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.11.0 -p "8081:80@loadbalancer" --agents 2
```

Next, from within the `spin-operator` directory, run the following commands to install the Spin runtime class and Spin Operator:

```console
kubectl apply -f spin-runtime-class.yaml
make install
```

Lastly, start the operator locally with the following command:

```console
make run
```

Great, now you have Spin Operator up and running on your cluster. This means you’re set to create and deploy SpinApps later on in the tutorial.

## Setting Up Ingress

Use the following command to set up ingress on your k8s cluster. This ensures traffic can reach your SpinApp once we’ve created it in future steps:

```console
# Setup ingress following this tutorial https://k3d.io/v5.4.6/usage/exposing_services/
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx
  annotations:
    ingress.kubernetes.io/ssl-redirect: "false"
spec:
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: keda-spinapp
            port:
              number: 80
EOF
```

Hit enter to create the ingress resource.

## Setting Up KEDA

Use the following command to setup KEDA on your k8s cluster using Helm. Different deployment methods are described at [Deploying KEDA on keda.sh](https://keda.sh/docs/2.13/deploy/):

```console
# Add the Helm repository
helm repo add kedacore https://kedacore.github.io/charts

# Update your Helm repositories
helm repo update

# Install the keda Helm chart into the keda namespace
helm install keda kedacore/keda --namespace keda --create-namespace
```

## Build and Store Spinapp in an OCI Registry

Next up we’re going to build the SpinApp we will be scaling and storing inside of a [ttl.sh](http://ttl.sh) registry. We've chosen TTL for ease of set-up, but you're welcome to use any OCI registry of your choosing, Change into the [apps/cpu-load-gen](https://github.com/spinkube/spin-operator/tree/hpa-tutorial/apps/cpu-load-gen) directory and build the SpinApp we’ve provided:

```console
# Build and publish the sample app
cd apps/cpu-load-gen
spin build
spin registry push ttl.sh/cpu-load-gen:1h
```

Note that the tag at the end of [ttl.sh/cpu-load-gen:1h](http://ttl.sh/cpu-load-gen:1h) indicates how long the image will last e.g. `1h` (1 hour). The maximum is `24h` and you will need to repush if ttl exceeds 24 hours.

## Deploy SpinApp and the KEDA ScaledObject

We can take a look at the SpinApp and the KEDA ScaledObject definitions in our deployment files below. As you can see, we have explicitly specified resource limits to `500m` of `cpu` (`spec.resources.limits.cpu`) and `500Mi` of `memory` (`spec.resources.limits.memory`) per SpinApp:

```yaml
# config/samples/keda-app.yaml
apiVapiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: keda-spinapp
spec:
  # TODO: Depend on a ghcr.io version of this image
  image: "ttl.sh/cpu-load-gen:1h"
  executor: containerd-shim-spin
  enableAutoscaling: true
  replicas: 1
  resources:
    limits:
      cpu: 500m
      memory: 500Mi
    requests:
      cpu: 100m
      memory: 400Mi
---
```

We will scale the instance count when we’ve reached a 50% utilization in `cpu` (`spec.triggers[cpu].metadata.value`). We’ve also instructed KEDA to scale our SpinApp horizontally within the range of 1 (`spec.minReplicaCount`) and 20 (`spec.maxReplicaCount`).:

```yaml
# config/samples/keda-scaledobject.yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: cpu-scaling
spec:
  scaleTargetRef:
    name: keda-spinapp
  minReplicaCount: 1
  maxReplicaCount: 20
  triggers:
    - type: cpu
      metricType: Utilization
      metadata:
        value: "50"
```

> The k8s documentation is the place to learn more about [limits and requests](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits). Consult the KEDA documentation to learn more about [ScaledObject](https://keda.sh/docs/2.13/concepts/scaling-deployments/#scaledobject-spec) and [KEDA's built-in scalers](https://keda.sh/docs/2.13/scalers/).

Let’s deploy the SpinApp and the KEDA ScaledObject instance onto our cluster with the following command:

```console
# Deploy the SpinApp
kubectl apply -f config/samples/keda-app.yaml
spinapp.core.spinoperator.dev/keda-spinapp created

# Deploy the ScaledObject
kubectl apply -f config/samples/keda-scaledobject.yaml
scaledobject.keda.sh/cpu-scaling created
```

You can see your running Spin application by running the following command:

```console
kubectl get spinapps

NAME          READY REPLICAS   EXECUTOR
keda-spinapp  1                containerd-shim-spin
```

You can also see your KEDA ScaledObject instance with the following command:

```console
kubectl get scaledobject

NAME          SCALETARGETKIND      SCALETARGETNAME   MIN   MAX   TRIGGERS   READY   ACTIVE   AGE
cpu-scaling   apps/v1.Deployment   keda-spinapp      1     20    cpu        True    True     7m
```

## Generate Load to Test Autoscale

Now let’s use Bombardier to generate traffic to test how well KEDA scales our SpinApp. The following Bombardier command will attempt to establish 40 connections during a period of 3 minutes (or less). If a request is not responded to within 5 seconds that request will timeout:

```console
# Generate a bunch of load
bombardier -c 40 -t 5s -d 3m http://localhost:8081
```

To watch the load, we can run the following command to get the status of our deployment:

```console
kubectl describe deploy keda-spinapp
...
---

Available      True    MinimumReplicasAvailable
Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   keda-spinapp-76db5d7f9f (1/1 replicas created)
Events:
  Type    Reason             Age   From                   Message
  ----    ------             ----  ----                   -------
  Normal  ScalingReplicaSet  84s   deployment-controller  Scaled up replica set hpa-spinapp-76db5d7f9f  to 2 from 1
  Normal  ScalingReplicaSet  69s   deployment-controller  Scaled up replica set hpa-spinapp-76db5d7f9f  to 4 from 2
  Normal  ScalingReplicaSet  54s   deployment-controller  Scaled up replica set hpa-spinapp-76db5d7f9f  to 8 from 4
  Normal  ScalingReplicaSet  39s   deployment-controller  Scaled up replica set hpa-spinapp-76db5d7f9f  to 16 from 8
  Normal  ScalingReplicaSet  24s   deployment-controller  Scaled up replica set hpa-spinapp-76db5d7f9f  to 20 from 16
```
