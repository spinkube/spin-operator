**Update - Monday 19th February 2024:**
The SpinKube documentation has been moved to [the documentation repository](https://github.com/spinkube/documentation/tree/main/content/en/docs) at the surface of the SpinKube GitHub organization. Please ensure that all new content is written in the new location. Content in this repository is currently being moved to the new location. The following README file and content herein is in the process of being migrated to [the new location](https://github.com/spinkube/documentation/tree/main/content/en/docs)

---

- [Spin Operator](#spin-operator)
  - [Getting Started](#getting-started)
  - [Deploying Spin Operator On Your Cluster](#deploying-spin-operator-on-your-cluster)
  - [Scaling SpinApps](#scaling-spinapps)
  - [Feedback](#feedback)
  - [Contributing](#contributing)
  - [Official Documentation](#official-documentation)

# Spin Operator

The Spin Operator enables deploying Spin applications to Kubernetes. It watches [SpinApp Custom Resources](./documentation/content/custom-resource-definition-reference.md) and realizes desired state in the Kubernetes cluster. This project was built using the Kubebuilder framework and contains a Spin App CRD and controller. To learn more about the SpinKube organization, visit our [project overivew documentation](./documentation/content/project-overview.md)

At this point in the priview, we recommend testing Spin Operator on a local k3d cluster via `make install`. The [quickstart guide](./documentation/content/quickstart.md) will walk you through prequisites and the installation workflow.

> > Spin Operator installation via Helm chart for remote clusters while in private preview is WIP and can tracked [here](https://github.com/spinkube/spin-operator/issues/54). In the meantime, please use the guidance from our quickstart guide.

## Official Documentation

Our content is under developement as markdown source files located at the [documentation](./documentation/) section. The following articles are ready for review

**Tutorials**

- [Quickstart](./documentation/content/quickstart.md)
- [Scale Spin Apps with Horizontal Pod Autoscaler](./documentation/content/scaling-spinapp-on-k8s-with-hpa.md)
- [Scale Spin Apps with Kubernetes Event Driver Autoscaler](./documentation/content/scaling-spinapp-on-k8s-with-keda.md)

**Glossary, Reference, & Misc.**

- [Glossary](./documentation/content/glossary-of-terms.md)
- [Spin App Custom Resource Definition](./documentation/content/custom-resource-definition-reference.md)
- [Troubleshooting](./documentation/content/troubleshooting.md)

Remaining articles are under construction. You're welcome to view and [open issues](https://github.com/spinkube/spin-operator/issues/new), but please proceed with caution as they are subject to change.

## Feedback

For questions or support, please visit our [Discord channel](https://discord.com/channels/926888690310053918/1200012610196738208). If you would like to file a feature, bug report, or documentation request, please [open a new issue](https://github.com/spinkube/spin-operator/issues/new).

## Contributing

Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for a guide on how to contribute to this project.

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
