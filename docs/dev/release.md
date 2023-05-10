# How to release a new version of all modules

On a release, all modules will be tagged with the same version. We use [`multimod`][1] for releases.
As a prerequisite, run `make install-tools` to ensure the necessary tooling is available.

To make a new release, follow these steps:

0. Make sure CI passes on [main](https://github.com/DataDog/opentelemetry-mapping-go/actions/workflows/test.yaml?query=branch%3Amain) and there are no release blockers.
1. Choose the new version number, `${VERSION}`. We follow semantic versioning and are currently doing `v0.x.y` releases.
2. Checkout to a new branch.
3. Update the version number on `versions.yaml` and commit the changes.
4. Run `chloggen update -v ${VERSION}` to update the changelog and commit the changes.
5. Run `make prerelease` and checkout to the branch created by this step. Open a PR to `main` from this branch and get it merged.
6. Checkout and pull the main branch locally. Run `git show HEAD` and make sure that it points to the commit from the previously merged PR.
7. Run `make push-tags` to push the tags.
8. Check that the new version is available on the [Github repository](https://github.com/DataDog/opentelemetry-mapping-go/tags).

## What to do if something goes wrong

If something goes wrong, it is important that you **do not remove or modify a tag once it has been pushed to Github**.

Instead, follow these steps:

1. Add a [`retract` directive][2] to all affected `go.mod` files, open a PR and commit it.
2. Fix the issue(s) in the release and commit the changes.
3. Follow the usual release process.

This will make tooling and bots ignore the retracted version. 
If a project has already updated to the new version, you may want to notify them directly.

[1]: https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod
[2]: https://go.dev/ref/mod#go-mod-file-retract
