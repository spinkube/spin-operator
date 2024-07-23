# Spin Operator

The Spin Operator enables deploying Spin applications to Kubernetes. It watches [SpinApp Custom Resources](https://www.spinkube.dev/docs/glossary/#spinapp-crd) and realizes the desired state in the Kubernetes cluster. This project was built using the Kubebuilder framework and contains a Spin App CRD and controller.

## Documentation

To learn more about the Spin Operator and the SpinKube organization, please visit [the official SpinKube documentation](https://www.spinkube.dev/docs/).

At this point in the preview, we recommend testing Spin Operator on a local k3d cluster via `make install`. The [quickstart guide](https://www.spinkube.dev/docs/spin-operator/quickstart/) will walk you through prequisites and the installation workflow.

> > Spin Operator installation via Helm chart for remote clusters while in private preview is WIP and can tracked [here](https://github.com/spinkube/spin-operator/issues/54). In the meantime, please use the guidance from our quickstart guide.

## Tutorials

There are a host of tutorials in the [SpinKube documentation](https://www.spinkube.dev/docs/install/). For example:

- [Quickstart](https://www.spinkube.dev/docs/install/quickstart/)
- [Deploying on Azure k8s service](https://www.spinkube.dev/docs/install/azure-kubernetes-service/)
- [Scaling Spin Apps with Horizontal Pod Autoscaler (HPA)](https://www.spinkube.dev/docs/topics/autoscaling/scaling-with-hpa/)
- [Scaling Spin Apps with Kubernetes Event Driver Autoscaler (KEDA)](https://www.spinkube.dev/docs/topics/autoscaling/scaling-with-keda/)
- [Running spin-operator locally](https://www.spinkube.dev/docs/contrib/running-locally/)
- [Running on a remote (non-local) K8s cluster](https://www.spinkube.dev/docs/contrib/running-on-a-cluster/)

## Feedback

The remaining articles are under construction. You're welcome to view and open both [Spin Operator](https://github.com/spinkube/spin-operator/issues) and [documentation](https://github.com/spinkube/documentation/issues) issues and feature requests. As this work is under development, please note that current features, functionality and supporting documentation are likely to change as the projects evolve and improvements are made.

For questions or support, please visit our [Discord channel](https://discord.com/channels/926888690310053918/1200012610196738208).

## Contributing (Spin Operator)

If you would like to contribute, please visit this [contributing](https://www.spinkube.dev/docs/contrib/) page.

## Contributing (Documentation)

If you would like to contribute to SpinKube and Spin Operator, please visit this [contributing](https://www.spinkube.dev/docs/contrib/) page.

The documentation is written using Hugo (as the static site generator), Docsy (as the technical documentation template) and GitHub pages (for hosting). However, during construction (prior to the website being rendered and publicly available) you are welcome to run a local copy of the documentation using the `hugo server` command. You can do so by following [these instructions](https://www.spinkube.dev/docs/contrib/#previewing-your-changes-locally).
