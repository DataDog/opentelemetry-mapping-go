# How to release a new version

On a release, all modules will be tagged with the same version. We use [`multimod`][1] for releases.

As a prerequisite, run `make install-tools` to ensure the necessary tooling is available.

To make a new release, follow these steps:

0. Make sure CI passes and there are no release blockers.
1. Choose the new version number, `${VERSION}`. We follow semantic versioning and are currently doing `0.x` releases.
2. Checkout to a new branch.
3. Modify the version number on the `pkgs` modset on `versions.yaml` and commit the changes.
4. If there are user-facing changes, run `chloggen update -v ${VERSION}` to update the changelog and commit the changes.
5. Run `make prerelease` and open a PR with the changes from all previous steps. Get it reviewed and merge it to the `main` branch.
6. Checkout locally to the main branch and make sure your repository's HEAD points to the commit you made on the previous step.
7. Run `make push-tags` to push the tags.
8. Check that the new version is available on the Github repository.

## If something goes wrong

If something goes wrong, it is important that you **do not remove or modify a tag once it has been pushed to Github**.

Instead, follow these steps:

1. Add a [`retract` directive][2] to all affected `go.mod` files, open a PR and commit it.
2. Fix the issue(s) in the release and commit the changes.
3. Follow the normal steps on the release process,

This will make tooling and bots ignore the retracted version. 
If a project has already updated to the new version, you may want to notify them directly.

[1]: https://github.com/open-telemetry/opentelemetry-go-build-tools/tree/main/multimod
[2]: https://go.dev/ref/mod#go-mod-file-retract
