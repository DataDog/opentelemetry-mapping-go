# Deprecation policy

Module versioning follows [semantic versioning][1] with a definition of breaking changes aligned with the [Go compatibility promise][2]. Currently, all modules are published under a `0.x` versioning scheme, which means that breaking changes may happen at any time. To ensure upgrades are easy for dependent projects, we ask contributors to abide by a simple deprecation policy, which is inspired by the [OpenTelemetry Collector breaking changes policy][3].

When making a breaking change to the public API follow these recommendations:

- When removing, replacing or refactoring a symbol or field, first mark it as [deprecated][4]. On the deprecation directive, note both the **deprecation version** as well as **any alternatives** available. To deprecate a module, use a [Go module deprecation directive][5] instead.
- Wait at least one version after marking something as deprecated to remove it.
- Add a changelog note both when marking something as deprecated as well as when removing it.

See the [OpenTelemetry Collector breaking changes policy][3] for examples on concrete refactor situations.

To minimize the number of breaking changes needed, carefully consider the public API exposed by your contributions and reduce its exposure by the usage of `internal` modules and by unexporting symbols and fields.

[1]: https://semver.org
[2]: https://go.dev/doc/go1compat
[3]: https://github.com/open-telemetry/opentelemetry-collector/blob/v0.69.0/CONTRIBUTING.md#breaking-changes
[4]: https://github.com/golang/go/wiki/Deprecated
[5]: https://go.dev/ref/mod#go-mod-file-module-deprecation
