- [Prerequisites](#prerequisites)
  - [Go](#go)
    - [TinyGo](#tinygo)
  - [Docker](#docker)
  - [Kubectl](#kubectl)
  - [K3d](#k3d)
  - [Helm](#helm)
  - [Bombardier](#bombardier)

# Prerequisites

The following prerequisites are required.

## Go

If building the Spin Operator from source or contributing to the development of Spin Operator then you will require [Go](https://go.dev/doc/install) version v1.21.0+ to be installed on your machine. Otherwise, please ignore this section, and move to the next prerequisite.

### TinyGo

Please also install the latest version of [TinyGo](https://tinygo.org/getting-started/install/)

## Docker

If you'd like to run Spin Operator locally, then please install [Docker](https://docs.docker.com/get-docker/) version 17.03+.

## Kubectl

If you'd like to manage your Spin applications with `kubectl`, then Spin Operator requires that you have [kubectl](https://kubernetes.io/docs/tasks/tools/) version v1.27.0+ installed.

## K3d

If running/deploying your Spin application involves the use of k3d, then the Spin Operator requires that you have [k3d](https://k3d.io/v5.6.0/?h=installation#installation) installed and that you have access to a Kubernetes v1.27.0 cluster.

## Helm

If running/deploying your Spin application involves the use of Helm, then the Spin Operator requires that you have [Helm](https://helm.sh/docs/intro/install/#helm) installed on your system.

## Bombardier

Installing [Bombardier](https://pkg.go.dev/github.com/codesenberg/bombardier) is **not** required to use Spin Operator. Bombardier is used in tutorials like [Scaling Spinapps on k8s With HPA](./scaling-spinapp-on-k8s-with-hpa.md) to generate load to test autoscaling.
