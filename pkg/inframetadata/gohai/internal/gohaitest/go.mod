// Deprecated: use github.com/DataDog/datadog-agent/pkg/opentelemetry-mapping-go/inframetadata instead.
module github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/gohai/internal/gohaitest

go 1.23.0

require (
	github.com/DataDog/datadog-agent/pkg/opentelemetry-mapping-go/inframetadata v0.71.0-devel.0.20250820180704-be0d2d237646
	github.com/DataDog/gohai v0.0.0-20230524154621-4316413895ee
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata => ../../../
