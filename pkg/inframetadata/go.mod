module github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata

go 1.19

require (
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.7.1
	github.com/stretchr/testify v1.8.4
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0014
	go.opentelemetry.io/collector/semconv v0.82.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/DataDog/datadog-agent/pkg/proto v0.48.0-beta.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tinylib/msgp v1.1.8 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.56.2 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes => ../otlp/attributes
