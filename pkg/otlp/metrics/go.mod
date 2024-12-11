module github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics

go 1.22.0

require (
	github.com/DataDog/datadog-agent/pkg/proto v0.52.0-devel
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest v0.21.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.21.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.21.0
	github.com/DataDog/sketches-go v1.4.4
	github.com/golang/protobuf v1.5.3
	github.com/lightstep/go-expohisto v1.0.0
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest v0.115.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/component v0.115.0
	go.opentelemetry.io/collector/component/componenttest v0.115.0
	go.opentelemetry.io/collector/pdata v1.21.0
	go.opentelemetry.io/otel v1.32.0
	go.uber.org/zap v1.27.0
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.115.0 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.115.0 // indirect
	go.opentelemetry.io/collector/semconv v0.115.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk v1.32.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.32.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../../internal/sketchtest
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes => ../attributes
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile => ../../quantile
)

retract v0.4.0 // see #107
