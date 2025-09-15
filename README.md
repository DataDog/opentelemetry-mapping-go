## Repository archive

This repository has been archived as of 15 September, 2025. The modules have been migrated to the github.com/DataDog/datadog-agent monorepo.
- pkg/inframetadata/* -> datadog-agent/pkg/opentelemetry-mapping-go/inframetadata/*
- pkg/internal/sketchtest -> datadog-agent/pkg/util/quantile/sketchtest
- pkg/otlp/* -> datadog-agent/pkg/opentelemetry-mapping-go/otlp/*
- pkg/quantile -> datadog-agent/pkg/util/quantile

# opentelemetry-mapping-go

This repository contains Go modules that implement [OpenTelemetry][1]-to-Datadog mapping for all telemetry signals as well as for semantic conventions.

These modules are used internally by Datadog in the [Datadog Agent OTLP ingest][2] and [OpenTelemetry Collector Datadog Exporter][3] implementations as well as related features, to ensure a consistent mapping between the two formats on all Datadog products. If building a new Datadog product that accepts telemetry in the [OTLP format][5], use the modules on this repository to convert to the Datadog public API format.

## Getting started

To get started contributing, clone this repository locally and check [CONTRIBUTING.md][4] for instructions on how to set up your development environment and send patches. You will need a supported Go compiler and the `make` utility for local testing and development.

[1]: https://opentelemetry.io
[2]: https://docs.datadoghq.com/opentelemetry/otlp_ingest_in_the_agent
[3]: https://docs.datadoghq.com/opentelemetry/otel_collector_datadog_exporter
[4]: CONTRIBUTING.md
[5]: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md
