# Overview

This is an OCI-compliant package that can be used to demonstrate how a Spin app interacts with Redis in a Kubernetes cluster.

# Usage

## Deploying the Spin app

Create a Kubernetes manifest file named `redis_client.yaml` with the following code:

```yaml
apiVersion: core.spinkube.dev/v1alpha1
kind: SpinApp
metadata:
  name: redis-spinapp
spec:
  image: "ghcr.io/spinkube/redis-sample"
  replicas: 1
  executor: containerd-shim-spin
  variables:
    - name: redis_endpoint
      value: redis://redis.default.svc.cluster.local:6379
```

Once created, run `kubectl apply -f redis_client.yaml`.

## Deploying Redis

Create a Kubernetes manifest file named `redis_db.yaml` with the following code:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  labels:
    app: redis
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
        - name: redis
          image: redis
          ports:
            - containerPort: 6379

---
apiVersion: v1
kind: Service
metadata:
  name: redis
spec:
  selector:
    app: redis
  ports:
    - protocol: TCP
      port: 6379
      targetPort: 6379
```

Once created, run `kubectl apply -f redis_db.yaml`.

## Interacting with the Spin app

In your terminal run `kubectl port-forward svc/redis-spinapp 3000:80`, then in a different terminal window, try the below commands:

### Place a key-value pair in Redis

```bash
curl --request PUT --data-binary "Hello, world\!" -H 'x-key: helloKey' localhost:3000
```

### Retrieve a value from Redis

```bash
curl -H 'x-key: helloKey' localhost:3000
```

### Delete a value from Redis

```bash
curl --request DELETE -H 'x-key: helloKey' localhost:3000
```
