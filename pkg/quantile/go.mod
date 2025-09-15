// Deprecated: use github.com/DataDog/datadog-agent/pkg/util/quantile instead.
module github.com/DataDog/opentelemetry-mapping-go/pkg/quantile

go 1.24.0

toolchain go1.24.5

require (
	github.com/DataDog/datadog-agent/pkg/util/quantile v0.71.0-devel.0.20250820180704-be0d2d237646
	github.com/DataDog/datadog-agent/pkg/util/quantile/sketchtest v0.72.0-devel
	github.com/DataDog/sketches-go v1.4.7
	github.com/dustin/go-humanize v1.0.1
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../internal/sketchtest

retract v0.4.0 // see #107
