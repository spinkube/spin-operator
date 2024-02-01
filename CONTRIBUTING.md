# CONTRIBUTING

We are delighted that you are interested in making spin-operator better! Thank you! This document will guide you through
making your first contribution to the project. We welcome and appreciate contributions of all types - opening issues,
fixing typos, adding examples, one-liner code fixes, tests, or complete features.

First, any contribution and interaction on any Fermyon project MUST follow our [Code of
Conduct](https://www.fermyon.com/code-of-conduct). Thank you for being part of an inclusive and open community!

If you plan on contributing anything complex, please go through the [open
issues](https://github.com/spinkube/spin-operator/issues) and [PR queue](https://github.com/spinkube/spin-operator/pulls)
first to make sure someone else has not started working on it. If it doesn't exist already, please [open an
issue](https://github.com/spinkube/spin-operator/issues/new) so you have a chance to get feedback from the community and
the maintainers before you start working on your feature.

## Making Code Contributions to spin-operator

The following guide is intended to make sure your contribution can get merged as soon as possible. First, make sure you
have the following prerequisites configured:

- `go` version v1.20.0+
- `docker` version 17.03+.
- `kubectl` version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster (This project is being developed using [`k3d`](https://k3d.io/v5.6.0/))
- `make`
- please ensure you [configure adding a GPG signature to your
  commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification)
  as well as appending a sign-off message (`git commit -S -s`)

Once you have set up the prerequisites and identified the contribution you want to make to spin-operator, make sure you
can correctly build the project:

```console
# clone the repository
$ git clone https://github.com/spinkube/spin-operator && cd spin-operator
# add a new remote pointing to your fork of the project
$ git remote add fork https://github.com/<your-username>/spin-operator
# create a new branch for your work
$ git checkout -b <your-branch>

# build spin-operator
$ make

# make sure compilation is successful
$ ./bin/manager --help

# run the tests and make sure they pass
$ make test
```

Now you should be ready to start making your contribution. To familiarize yourself with the spin-operator project,
please read the [README](https://github.com/spinkube/spin-operator). Since most of spin-operator is written in Go, we try
to follow the common Go coding conventions. If applicable, add unit or integration tests to ensure your contribution is
correct.

## Before You Commit

- Format the code (`go fmt ./...`)
- Run Clippy (`go vet ./...`)
- Run the lint task (`make lint` or `make lint-fix`)
- Build the project and run the tests (`make test`)

spin-operator enforces lints and tests as part of continuous integration - running them locally will save you a
round-trip to your pull request!

If everything works locally, you're ready to commit your changes.

## Committing and Pushing Your Changes

We require commits to be signed both with an email address and with a GPG signature.

> Because of the way GitHub runs enforcement, the GPG signature isn't checked until after all tests have run. Be sure to
> GPG sign up front, as it can be a bit frustrating to wait for all the tests and then get blocked on the signature!

```console
$ git commit -S -s -m "<your commit message>"
```

Some contributors like to follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) convention
for commit messages.

We try to only keep useful changes as separate commits - if you prefer to commit often, please cleanup the commit
history before opening a pull request.

Once you are happy with your changes you can push the branch to your fork:

```console
# "fork" is the name of the git remote pointing to your fork
$ git push fork
```

Now you are ready to create a pull request. Thank you for your contribution!
