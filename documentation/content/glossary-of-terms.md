- [Glossary of Terms](#glossary-of-terms)
  - [Chart](#chart)
  - [Cluster](#cluster)
  - [Container Runtime](#container-runtime)
  - [Controller](#controller)
  - [Custom Resource (CR)](#custom-resource-cr)
  - [Custom Resource Definition (CRD)](#custom-resource-definition-crd)
  - [Helm](#helm)
  - [Image](#image)
  - [Kubernetes (K8s)](#kubernetes-k8s)
  - [Open Container Initiative (OCI)](#open-container-initiative-oci)
  - [Pod](#pod)
  - [Role Based Access Control (RBAC)](#role-based-access-control-rbac)
  - [Runtime Class](#runtime-class)
  - [Scheduler](#scheduler)
  - [Service](#service)
  - [Spin](#spin)
  - [Spin Operator](#spin-operator)

# Glossary of Terms

The following glossary of terms is in the context of deploying, scaling, automating and managing Spin applications in containerized environments.

## Chart

Helm packages are known as "charts". The main Spin Operator chart does not include the SpinApp CRD or any non-namespace or cluster-level resources.

## Cluster

TODO

## Container Runtime

TODO

## Controller

TODO

## Custom Resource (CR)

A CR facilitates the storage and retrieval of your own API Objects (as structured data). A Spin application can be described as a CR.

## Custom Resource Definition (CRD)

A CRD defines your Custom Resources (CR). For example, the following `.yaml` file describes a `SpinApp` using CRD syntax:

```yaml
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: simple-spinapp
spec:
  image: "ghcr.io/deislabs/containerd-wasm-shims/examples/spin-rust-hello:v0.10.0"
  replicas: 1
  runtime: "containerd-shim-spin"
```

> SpinApp CRDs are kept separate from Helm. If using Helm, CustomResourceDefinition (CRD) resources will need to be installed prior to installing the Heml chart.

## Helm

TODO

## Image

TODO

## Kubernetes (K8s)

TODO

## Open Container Initiative (OCI)

TODO

## Pod

A pod is a group of containers that can share resources.

## Role Based Access Control (RBAC)

TODO

## Runtime Class

A RuntimeClass isn't a namespaced resource. A RuntimeClass is not part of a Helm chart.

## Scheduler

TODO

## Service

In Kubernetes, a Service is an abstraction that defines a logical set of Pods that enables clients to interact with a consistent set of Pods, regardless of whether the code is designed for a cloud-native environment or a containerized legacy application.

## Spin

Spin is a framework designed for building and running event-driven microservice applications using WebAssembly (Wasm) components.

## Spin Operator

Spin Operator is a Kubernetes (K8s) operator in charge of handling the lifecycle of Spin applications based on their SpinApp resources.
