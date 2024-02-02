- [Troubleshooting](#troubleshooting)
  - [Cluster Already Exists](#cluster-already-exists)
    - [Cluster Information](#cluster-information)
    - [Cluster Delete](#cluster-delete)
  - [Too long: must have at most 262144 bytes](#too-long-must-have-at-most-262144-bytes)
  - [Redis Operator](#redis-operator)
  - [error: requires go version](#error-requires-go-version)

# Troubleshooting

The following is a list of common error messages and potential troubleshooting suggestions that might assist you with your work.

## Cluster Already Exists

When trying to create a cluster (e.g. a cluster named `wasm-cluster`) you may recieve an error message similar to the following:

```bash
FATA[0000] Failed to create cluster 'wasm-cluster' because a cluster with that name already exists
```

### Cluster Information

With `k3d` installed, you can use the following command to get a cluster list:

```bash
$ k3d cluster list
NAME           SERVERS   AGENTS   LOADBALANCER
wasm-cluster   1/1       2/2      true
```

With `kubectl` installed, you can use the following command to dump cluster information (this is much more verbose):

```bash
$ kubectl cluster-info dump
```

### Cluster Delete

With `k3d` installed, you can delete the cluster by name, as shown in the command below:

```bash
$ k3d cluster delete wasm-cluster
INFO[0000] Deleting cluster 'wasm-cluster'
INFO[0002] Deleting cluster network 'k3d-wasm-cluster'
INFO[0002] Deleting 1 attached volumes...
INFO[0002] Removing cluster details from default kubeconfig...
INFO[0002] Removing standalone kubeconfig file (if there is one)...
INFO[0002] Successfully deleted cluster wasm-cluster!
```

## Too long: must have at most 262144 bytes

When running `kubectl apply -f my-file.yaml`, the following error can occur of the `.yam.` file is too large:

```bash
Too long: must have at most 262144 bytes
```

Using the `--server-side=true` option resolves this issue:

```bash
kubectl apply --server-side=true -f my-file.yaml
```

## Redis Operator

Noted an error when installing Redis Operator:

```bash
$ helm repo add redis-operator https://spotahome.github.io/redis-operator
"redis-operator" has been added to your repositories
$ helm repo update
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "redis-operator" chart repository
Update Complete. ⎈Happy Helming!⎈
$ helm install redis-operator redis-operator/redis-operator
Error: INSTALLATION FAILED: failed to install CRD crds/databases.spotahome.com_redisfailovers.yaml: error parsing : error converting YAML to JSON: yaml: line 4: did not find expected node content
```

Used the following commands to enforce using a different version of Redis Operator (whilst waiting on [this PR fix](https://github.com/spotahome/redis-operator/pull/685) to be merged).

```bash
$ helm install redis-operator redis-operator/redis-operator --version 3.2.9
NAME: redis-operator
LAST DEPLOYED: Mon Jan 22 12:33:54 2024
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

## error: requires go version

When building apps like the [cpu-load-gen](../../apps/cpu-load-gen/) Spin app, you may get the following error if your TinyGo is not up to date. The error requires go version `1.18` through `1.20` but this is not necessarily the case. It **is** recommended that you have the latest go installed e.g. `1.21` and downgrading is not necessary. Instead please go ahead and [install the latest version of TinyGo](./prerequisites.md#tinygo) to resolve this error:

```bash
user@user:~/spin-operator/apps/cpu-load-gen$ spin build
Building component cpu-load-gen with `tinygo build -target=wasi -gc=leaking -no-debug -o main.wasm main.go`
error: requires go version 1.18 through 1.20, got go1.21
```
