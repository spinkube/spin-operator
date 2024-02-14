- [Spin Operator](#spin-operator)

# Project Overview

[SpinKube](https://github.com/spinkube) is a new open source project that streamlines the experience of developing, deploying, and operating Wasm workloads on Kubernetes, using [Spin](https://github.com/fermyon/spin) in tandem with the `[runwasi](https://github.com/containerd/runwasi)` and [KWasm](https://kwasm.sh/) open source projects.

With SpinKube, you can leverage the advantages of using WebAssembly (Wasm) for your workloads:

- Artifacts are significantly smaller in size compared to container images.
- Artifacts can be quickly fetched over the network and started much faster (\*Note: We are aware of several optimizations that still need to be implemented to enhance the startup time for workloads).
- Substantially fewer resources are required during idle times.

All of this while being able to integrate with with Kubernetes primitives (DNS, probes, autoscaling, metrics, and a lot more cloud native and CNCF projects) thanks to Spin Operator.

Spin operator watches Spin App Custom Resources and realizes the desired state in the Kubernetes cluster. ThiThe foundation of this project was built using the kubebuilder framework and contains a Spin App Custom Resource Definition (CRD) and controller.
