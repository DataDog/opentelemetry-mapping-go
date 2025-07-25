module github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics

go 1.23.0

require (
	github.com/DataDog/datadog-agent/pkg/proto v0.70.0-devel
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest v0.30.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.30.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.30.0
	github.com/DataDog/sketches-go v1.4.7
	github.com/golang/protobuf v1.5.4
	github.com/lightstep/go-expohisto v1.0.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest v0.130.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/component v1.36.0
	go.opentelemetry.io/collector/component/componenttest v0.130.0
	go.opentelemetry.io/collector/pdata v1.36.0
	go.opentelemetry.io/otel v1.37.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.130.0 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.36.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.130.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.12.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../../internal/sketchtest
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes => ../attributes
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile => ../../quantile
)

retract v0.4.0 // see #107
