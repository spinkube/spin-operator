- [Scaling Spinapp on Kubernetes (k8s) With Horizontal Pod Autoscaling (HPA)](#scaling-spinapp-on-kubernetes-k8s-with-horizontal-pod-autoscaling-hpa)
  - [Prerequisites](#prerequisites)
  - [Fetch Spin Operator (Source Code)](#fetch-spin-operator-source-code)
  - [Setting Up k8s Cluster](#setting-up-k8s-cluster)
  - [Set Up Ingress](#set-up-ingress)
  - [Build and Store Spinapp in a TTL Registry](#build-and-store-spinapp-in-a-ttl-registry)
  - [Deploy SpinApp and HPA](#deploy-spinapp-and-hpa)
  - [Generate Load to Test Autoscale](#generate-load-to-test-autoscale)

# Scaling Spinapp on Kubernetes (k8s) With Horizontal Pod Autoscaling (HPA)

Horizontal scaling, in the k8s sense, means deploying more pods to meet demand (different from vertical scaling whereby more memory and CPU resources are assigned to already running pods). In this tutorial, we configure [HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) to dynamically scale the instance count of our SpinApps to meet the demand.

## Prerequisites

> We use k3d to run a k8s cluster locally as part of this tutorial, but you can follow these steps to configure HPA autoscaling on your desired k8s environment.

Please see the [Go](./prerequisites.md#go), [Docker](./prerequisites.md#docker), [Kubectl](./prerequisites.md#kubectl), [k3d](./prerequisites.md#k3d) and [Bombardier](#prerequisites#bombardier) sections in the [Prerequisites](./prerequisites.md) page and fulfill those prerequisite requirements before continuing.

## Fetch Spin Operator (Source Code)

If you haven't already, please go ahead and clone the Spin Operator repository:

```bash
git clone https://github.com/spinkube/spin-operator.git
```

Change into the Spin Operator directory:

```bash
cd spin-operator
```

## Setting Up k8s Cluster

Run the following command to create a k8s k3d cluster that has [the containerd-wasm-shims](https://github.com/deislabs/containerd-wasm-shims) pre-requisites installed: If you have a k3d cluster already, please feel free to use it:

```sql
k3d cluster create wasm-cluster-scale --image ghcr.io/deislabs/containerd-wasm-shims/examples/k3d:v0.10.0 -p "8081:80@loadbalancer" --agents 2
```

Next, from within the `spin-operator` directory, run the following commands to install the Spin runtime class and Spin Operator:

```sql
kubectl apply -f spin-runtime-class.yaml
make install
```

Lastly, start the operator with the following command:

```sql
make run
```

Great, now you have Spin Operator up and running on your cluster. This means you’re set to create and deploy SpinApps later on in the tutorial.

## Set Up Ingress

Use the following command to set up ingress on your k8s cluster. This ensures traffic can reach your SpinApp once we’ve created it in future steps:

```bash
# Setup ingress following this tutorialhttps://k3d.io/v5.4.6/usage/exposing_services/
cat <<EOF >nginx-ingress.yaml
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
            name: hpa-spinapp
            port:
              number: 80
EOF
```

Hit enter to create the ingress resource. It can live inside the `config/samples` directory alongside other sample applications.

## Build and Store Spinapp in an OCI Registry

Next up we’re going to build the SpinApp we will be scaling and storing inside of a [ttl.sh](http://ttl.sh) registry. Change into the [apps/cpu-load-gen](https://github.com/spinkube/spin-operator/tree/hpa-tutorial/apps/cpu-load-gen) directory and build the SpinApp we’ve provided:

```bash
# Build and publish the sample app
cd apps/cpu-load-gen
spin build
spin registry push ttl.sh/cpu-load-gen:1h
```

Note that the tag at the end of [ttl.sh/cpu-load-gen:1h](http://ttl.sh/cpu-load-gen:1h) indicates how long the image will last e.g. `1h` (1 hour). The maximum is `24h` and you will need to repush if ttl exceeds 24 hours.

<aside>
☁️ In the future, we will be storing this application in an OCI registry (like [ghcr.io](https://docs.github.com/en/packages/learn-github-packages)) for more permanent persistence. For now, we’re making the recommendation to use [ttl.sh](http://ttl.sh) for convenience.
</aside>

## Deploy SpinApp and HPA

We can take a look at the SpinApp and HPA definitions in our deployment file below/. As you can see, we have set our `resources` -> `limits` to `500m` of `cpu` and `500Mi` of `memory` per Spin application and we will scale the instance count when we’ve reached a 50% utilization in `cpu` and `memory`. We’ve also defined support a maximum [replica](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#replicas) count of 10 and a minimum replica count of 1:

```yaml
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: hpa-spinapp
spec:
  # TODO: Depend on a ghcr.io version of this image
  image: "ttl.sh/cpu-load-gen:1h"
  enableAutoscaling: true
  resources:
    limits:
      cpu: 500m
      memory: 500Mi
    requests:
      cpu: 100m
      memory: 400Mi
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: spinapp-autoscaler
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: hpa-spinapp
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
```

<aside>
☁️ TODO - we need to define which metrics we’re supporting with HPA. Is it the [base set used by containers](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#container-resource-metrics)?

</aside>

The k8s documentation is the place to learn more about [limits and requests](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits) and [other metrics supported by HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#container-resource-metrics).

Let’s deploy the SpinApp and the HPA instance onto our cluster with the following command:

```bash
kubectl apply -f config/samples/hpa.yaml
```

You can see your running Spin application by running the following command:

```bash
kubectl get spinapps
NAME          AGE
hpa-spinapp   92m
```

You can also see your HPA instance with the following command:

```bash
NAME                 REFERENCE                TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
spinapp-autoscaler   Deployment/hpa-spinapp   6%/50%    1         10        1          97m
```

## Generate Load to Test Autoscale

Now let’s use Bombardier to generate traffic to test how well HPA scales our SpinApp. The following Bombardier command will attempt to establish 40 connections during a period of 3 minutes (or less). If a request is not responded to within 5 seconds that request will timeout:

```bash
# Generate a bunch of load
bombardier -c 40 -t 5s -d 3m http://localhost:8081
```

To watch the load, we can run the following command to get the status of our deployment:

```bash
kubectl describe deploy hpa-spinapp
...
---

Available      True    MinimumReplicasAvailable
Progressing    True    NewReplicaSetAvailable
OldReplicaSets:  <none>
NewReplicaSet:   hpa-spinapp-544c649cf4 (1/1 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  11m    deployment-controller  Scaled up replica set hpa-spinapp-544c649cf4 to 1
  Normal  ScalingReplicaSet  9m45s  deployment-controller  Scaled up replica set hpa-spinapp-544c649cf4  to 4
  Normal  ScalingReplicaSet  9m30s  deployment-controller  Scaled up replica set hpa-spinapp-544c649cf4  to 8
  Normal  ScalingReplicaSet  9m15s  deployment-controller  Scaled up replica set hpa-spinapp-544c649cf4  to 10
```
