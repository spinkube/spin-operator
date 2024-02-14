This chart directory currently holds files related to the
[Spin Operator Helm chart](../../charts/spin-operator/) that we inject after
[helmify](../../README.md#packaging-and-deployment-via-helm) performs chart
(re-)generation.

As an example, helmify produces an auto-generated
`values.yaml` but we'd like to provide a more ergonomic and descriptive
version for users, eg with ample comments.

If we no longer use helmify, this configuration directory can be removed.
