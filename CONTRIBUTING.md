# Contributing

All Go modules in this repository are officially supported for usage in Datadog products only.

## Submitting issues

This repository contains Go modules used in the implementation of the Datadog Agent and the OpenTelemetry Collector Datadog components. Whenever possible, prefer reporting the issue directly on the upstream issue trackers instead.

## Pull requests

To submit a pull request on this repository, we assume you have installed a supported Go compiler and the `make` utility. Additional tooling can be installed by running `make install-tools`.

When submitting a pull request, take into account the following guidelines:

- Ensure your contribution is properly tested and that tests pass locally.
- Open your PR against the `main` branch.
- If your PR results in user-facing changes, add a changelog note with `make chlog-new`.
