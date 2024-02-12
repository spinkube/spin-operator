- [Running on Azure Kubernetes Service (AKS)](#running-on-azure-kubernetes-service-aks)
  - [Prerequisites](#prerequisites)
  - [Provisioning the necessary Azure infrastructure](#provisioning-the-necessary-azure-infrastructure)
  - [Deploying the Spin Operator](#deploying-the-spin-operator)
  - [Deploying a Spin App to AKS](#deploying-a-spin-app-to-aks)
  - [Verifying the Spin App](#verifying-the-spin-app)
  - [Removing the Azure Infrastructure](#removing-the-azure-infrastructure)

# Running on Azure Kubernetes Service (AKS)

## Prerequisites

Please see the following sections in the [Prerequisites](./prerequisites.md) page and fulfill those prerequisite requirements before continuing:

- [kubectl](./prerequisites.md#kubectl) - the Kubernetes CLI
- [Helm](./prerequisites.md#helm) - the package manager for Kubernetes
- [Azure CLI](./prerequisites.md#azure-cli) - cross-platform CLI for managing Azure resources

## Provisioning the necessary Azure Infrastructure

Before you dive into deploying Spin Operator on Azure Kubernetes Service (AKS), the underlying cloud infrastructure must be provisioned. For the sake of this article, you will provision a simple AKS cluster. (Alternatively, you can setup the AKS cluster following [this guide from Microsoft](https://learn.microsoft.com/en-us/azure/aks/tutorial-kubernetes-deploy-cluster?tabs=azure-cli).)

```shell
# Login with Azure CLI
az login

# Select the desired Azure Subscription
az account set --subscription <YOUR_SUBSCRIPTION>

# Create an Azure Resource Group
az group create --name rg-spin-operator \
    --location germanywestcentral

# Create an AKS cluster
az aks create --name aks-spin-operator \
    --resource-group rg-spin-operator \
    --node-count 1 \
    --tier free \
    --generate-ssh-keys
```

Once the AKS cluster has been provisioned, use the `aks get-credentials` command to download credentials for `kubectl`:

```shell
# Download credentials for kubectl
az aks get-credentials --name aks-spin-operator \
    --resource-group rg-spin-operator
```

For verification, you can use `kubectl` to browse common resources inside of the AKS cluster:

```shell
# Browse namespaces in the AKS cluster
kubectl get namespaces

NAME              STATUS   AGE
default           Active   3m
kube-node-lease   Active   3m
kube-public       Active   3m
kube-system       Active   3m
```

## Deploying the Spin Operator

First, the [Custom Resource Definition (CRD)](glossary-of-terms#custom-resource-definition-crd) and the [Runtime Class](./glossary-of-terms.md#runtime-class) for `wasmtime-spin-v2` must be installed.

<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.crds.yaml' -->
<!-- TODO: replace with e.g. 'kubectl apply -f https://github.com/spinkube/spin-operator/releases/download/v0.1.0-rc.1/spin-operator.runtime-class.yaml' -->

```shell
# Install the CRD
make install

# Install the RuntimeClass
kubectl apply -f spin-runtime-class.yaml
```

The following installs the chart with the release name `spin-operator`:

<!-- TODO: remove '--devel' flag once we have our first non-prerelease chart available, e.g. when v0.1.0 of this project is released and public -->

```shell
helm install spin-operator \
  --namespace spin-operator \
  --create-namespace \
  --devel \
  --wait \
  oci://ghcr.io/spinkube/spin-operator
```

The Spin Operator chart has a dependency on [Kwasm](https://kwasm.sh/), which you use to install `containerd-wasm-shim` on the Kubernetes node(s):

```shell
# Trigger containerd-wasm-shim installation on all nodes
kubectl annotate node --all kwasm.sh/kwasm-node=true
```

To verify `containerd-wasm-shim` installation, you can inspect the logs from the Kwasm Operator:

```shell
# Ispect logs from the Kwasm Operator
kubectl logs -n spin-operator -l app.kubernetes.io/name=kwasm-operator
{"level":"info","node":"aks-nodepool1-31687461-vmss000000","time":"2024-02-12T11:23:43Z","message":"Trying to Deploy on aks-nodepool1-31687461-vmss000000"}
{"level":"info","time":"2024-02-12T11:23:43Z","message":"Job aks-nodepool1-31687461-vmss000000-provision-kwasm is still Ongoing"}
{"level":"info","time":"2024-02-12T11:24:00Z","message":"Job aks-nodepool1-31687461-vmss000000-provision-kwasm is Completed. Happy WASMing"}
```

The final step for setting up Spin Operator is creating a `SpinAppExecutor`:

```shell
# Create a SpinAppExecutor
kubectl apply -f https://github.com/spinkube/spin-operator/blob/main/config/samples/shim-executor.yaml
```

## Deploying a Spin App to AKS

To validate the Spin Operator deployment, you will deploy a simple Spin App to the AKS cluster. The following command will install a simple Spin App using the `SpinApp` CRD you provisioned in the previous section:

```shell
# Deploy a sample Spin app
kubectl apply -f https://github.com/spinkube/spin-operator/blob/main/config/samples/simple.yaml
```

## Verifying the Spin App

Configure port forwarding from port `8080` of your local machine to port `80` of the Kubernetes service which points to the Spin App you installed in the previous section:

```shell
kubectl port-forward services/simple-spinapp 8080:80
Forwarding from 127.0.0.1:8080 -> 80
Forwarding from [::1]:8080 -> 80
```

Send a HTTP request to [http://127.0.0.1:8080/hello](http://127.0.0.1:8080/hello) using [`curl`](https://curl.se/):

```shell
# Send an HTTP GET request to the Spin App
curl -iX GET http://localhost:8080/hello
HTTP/1.1 200 OK
transfer-encoding: chunked
date: Mon, 12 Feb 2024 12:23:52 GMT

Hello world from Spin!%
```

## Removing the Azure infrastructure

To delete the Azure infrastructure created as part of this article, use the following command:

```shell
# Remove all Azure resources
az group delete --name rg-spin-operator \
    --no-wait \
    --yes
```
