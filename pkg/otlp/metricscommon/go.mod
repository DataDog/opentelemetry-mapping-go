module github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metricscommon

go 1.20

require (
	github.com/DataDog/datadog-agent/pkg/proto v0.49.0-devel
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.8.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.8.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/DataDog/sketches-go v1.4.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../../internal/sketchtest
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes => ../attributes
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile => ../../quantile
)

retract v0.4.0 // see #107
