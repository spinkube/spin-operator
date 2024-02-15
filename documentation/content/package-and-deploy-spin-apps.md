- [Package and Deploy Spin Apps](#package-and-deploy-spin-apps)
- [Prerequisites](#prerequisites)
- [Creating a new Spin App](#creating-a-new-spin-app)
- [Packaging and Distributing Spin Apps](#packaging-and-distributing-spin-apps)
- [Deploying Spin Apps](#deploying-spin-apps)
- [Distributing and Deploying Spin Apps via private registries](#distributing-and-deploying-spin-apps-via-private-registries)

# Package and Deploy Spin Apps

## Prerequisites

Ensure the necessary [prerequisites](./prerequisites.md) are installed.

For this tutorial in particular, you should either have the Spin Operator [running locally](./running-locally.md) or [running on your Kubernetes cluster](./running-on-a-cluster.md).

## Creating a new Spin App

You use the `spin` CLI, to create a new Spin App. The `spin` CLI provides different templates, which you can use to quickly create different kinds of Spin Apps. For demonstration purposes, you will use the `http-go` template to create a simple Spin App.

```shell
# Create a new Spin App using the http-go template
spin new -t http-go hello-spin
Description: A simple Spin App for demonstration purposes
HTTP path: /...

# Navigate into the hello-spin directory
cd hello-spin
```

The `spin` CLI created all necessary files within `hello-spin`. Besides the Spin Manifest (`spin.toml`), you can find the actual implementation of the app in `main.go`:

```go
package main

import (
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Hello Fermyon!")
	})
}

func main() {}
```

This implementation will respond to any incoming HTTP request, and return an HTTP response with a status code of 200 (`Ok`) and send `Hello Fermyon` as the response body.

You can test the app on your local machine by invoking the `spin up` command from within the `hello-spin` folder.

## Packaging and Distributing Spin Apps

Spin Apps are packaged and distributed as OCI artifacts. By leveraging OCI artifacts, Spin Apps can be distributed using any registry that implements the [Open Container Initiative Distribution Specification](https://github.com/opencontainers/distribution-spec/tree/main) (a.k.a. "OCI Distribution Spec").

The `spin` CLI simplifies packaging and distribution of Spin Apps and provides an atomic command for this (`spin registry push`). You can package and distribute the `hello-spin` app that you created as part of the previous section like this:

```shell
# Package and Distribute the hello-spin app
spin registry push --build ttl.sh/hello-spin:24h
```

> It is a good practice to add the `--build` flag to `spin registry push`. It prevents you from accidentally pushing an outdated version of your Spin App to your registry of choice.

## Deploying Spin Apps

To deploy Spin Apps to a Kubernetes cluster which has Spin Operator running, you create a custom resource of type `SpinApp`:

```yaml
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: hello-spin
spec:
  image: ttl.sh/hello-spin:24h
  replicas: 1
  executor: containerd-shim-spin
```

You can deploy the manifest using `kubectl` as shown here:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: hello-spin
spec:
  image: ttl.sh/hello-spin:24h
  replicas: 1
  executor: containerd-shim-spin
EOF
```

> The `SpinApp` Custom Resource Definition (CRD) contains additional properties that you use to configure the Spin App according to your needs. However, for the sake of this article, we kept it as easy as possible.

## Distributing and Deploying Spin Apps via private registries

It is quite common to distribute Spin Apps through private registries that require some sort of authentication. To publish a Spin App to a private registry, you have to authenticate using the `spin registry login` command.

For demonstration purposes, you will now distribute the Spin App via GitHub Container Registry (GHCR). You can follow [this guide by GitHub](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-with-a-personal-access-token-classic) to create a new personal access token (PAT), which is required for authentication.

```shell
# Store PAT and GitHub username as environment variables
export CR_PAT=YOUR_TOKEN
export GH_USER=YOUR_GITHUB_USERNAME

# Authenticate spin CLI with GHCR
echo $GH_PAT | spin registry login ghcr.io -u $GH_USER --password-stdin

Successfully logged in as YOUR_GITHUB_USERNAME to registry ghcr.io
```

Once authentication succeeded, you can use `spin registry push` to push your Spin App to GHCR:

```shell
# Push hello-spin to GHCR
spin registry push --build ghcr.io/$GH_USER/hello-spin:0.0.1

Pushing app to the Registry...
Pushed with digest sha256:1611d51b296574f74b99df1391e2dc65f210e9ea695fbbce34d770ecfcfba581
```

In Kubernetes you store authentication information as secret of type `docker-registry`. The following snippet shows how to create such a secret with `kubectl` leveraging the environment variables, you specified in the previous section:

```shell
# Create Secret in Kubernetes
kubectl create secret docker-registry ghcr \
    --docker-server ghcr.io \
    --docker-username $GH_USER \
    --docker-password $CR_PAT

secret/ghcr created
```

Finally, you create a new `SpinApp` and link the secret using the `imagePullSecrets` property:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: core.spinoperator.dev/v1
kind: SpinApp
metadata:
  name: hello-spin
spec:
  image: ghcr.io/$GH_USER/hello-spin:0.0.1
  imagePullSecrets:
    - name: ghcr
  replicas: 1
  executor: containerd-shim-spin
EOF
```
