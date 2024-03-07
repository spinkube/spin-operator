# Spin Operator

The Spin Operator enables deploying Spin applications to Kubernetes. It watches [SpinApp Custom Resources](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/reference/custom-resource-definition.md) and realizes the desired state in the Kubernetes cluster. This project was built using the Kubebuilder framework and contains a Spin App CRD and controller.

![spin-operator diagram](https://github.com/spinkube/spin-operator/assets/686194/bf07365f-1d07-421a-864f-d77c0a27a764)

## Documentation

To learn more about the Spin Operator and the SpinKube organization, please visit [the official Spin Operator documentation](https://github.com/spinkube/documentation/tree/main/content/en/docs/spin-operator) which is housed inside the [the official SpinKube documentation](https://github.com/spinkube/documentation/tree/main/content/en/docs).

At this point in the preview, we recommend testing Spin Operator on a local k3d cluster via `make install`. The [quickstart guide](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/quickstart/_index.md) will walk you through prequisites and the installation workflow.

> > Spin Operator installation via Helm chart for remote clusters while in private preview is WIP and can tracked [here](https://github.com/spinkube/spin-operator/issues/54). In the meantime, please use the guidance from our quickstart guide.

## Tutorials

There are a host of tutorials in the [Spin Operator tutorials](https://github.com/spinkube/documentation/tree/main/content/en/docs/spin-operator/tutorials) directory of the documentation. For example:

- [Quickstart](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/quickstart/_index.md)
- [Running Locally](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/tutorials/running-locally.md)
- [Running on a remote (non-local) K8s cluster](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/tutorials/running-on-a-cluster.md)
- [Deploying on Azure k8s service](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/tutorials/deploy-on-azure-kubernetes-service.md)
- [Scaling Spin Apps with Horizontal Pod Autoscaler (HPA)](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/tutorials/scaling-with-hpa.md)
- [Scaling Spin Apps with Kubernetes Event Driver Autoscaler (KEDA)](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/tutorials/scaling-with-keda.md)

## Feedback

The remaining articles are under construction. You're welcome to view and open both [Spin Operator](https://github.com/spinkube/spin-operator/issues) and [documentation](https://github.com/spinkube/documentation/issues) issues and feature requests. As this work is under development, please note that current features, functionality and supporting documentation are likely to change as the projects evolve and improvements are made.

For questions or support, please visit our [Discord channel](https://discord.com/channels/926888690310053918/1200012610196738208).

## Contributing (Spin Operator)

If you would like to contribute, please visit this [contributing](https://github.com/spinkube/documentation/blob/main/content/en/docs/spin-operator/contributing/_index.md) page.

## Contributing (Documentation)

If you would like to contribute to SpinKube and Spin Operator, please visit this [contributing](https://github.com/spinkube/documentation/blob/main/content/en/docs/contribution-guidelines/_index.md) page.

The documentation is written using Hugo (as the static site generator), Docsy (as the technical documentation template) and GitHub pages (for hosting). However, during construction (prior to the website being rendered and publicly available) you are welcome to run a local copy of the documentation using the `hugo server` command. You can do so by following [these instructions](https://github.com/spinkube/documentation/blob/main/content/en/docs/contribution-guidelines/_index.md#previewing-your-changes-locally).
