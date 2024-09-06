# spin-operator release process

The vast majority of the release process is handled after you push a new git tag.

First, let's start by setting some environment variables we can reference later.

```console
export TAG=v0.x.0 #CHANGEME
```

To push a tag, do the following:

```console
git checkout main
git remote add upstream git@github.com:spinkube/spin-operator
git pull upstream main
git tag --sign $TAG --message "Release $TAG"
git push upstream $TAG
```

Observe that the [CI run for that tag](https://github.com/spinkube/spin-operator/actions) completed.

Bump the Helm chart versions. See #311 for an example.

Next, you'll need to update the documentation:

```console
git clone git@github.com:spinkube/documentation
cd documentation
```

Change all references from the previous version to the new version.

Contribute those changes and open a PR.

As an optional step, you can run a set of smoke tests to ensure the latest release works as expected.

Finally, announce the new release on the #spinkube CNCF Slack channel.
