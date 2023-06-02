module github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metrics

go 1.19

require (
	github.com/DataDog/datadog-agent/pkg/trace v0.45.0-rc.6
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest v0.2.3
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.2.3
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile v0.2.3
	github.com/DataDog/sketches-go v1.4.2
	github.com/golang/protobuf v1.5.3
	github.com/lightstep/go-expohisto v1.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/stretchr/testify v1.8.3
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0012
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	go.opentelemetry.io/collector/semconv v0.78.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../../internal/sketchtest
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes => ../attributes
	github.com/DataDog/opentelemetry-mapping-go/pkg/quantile => ../../quantile
)
