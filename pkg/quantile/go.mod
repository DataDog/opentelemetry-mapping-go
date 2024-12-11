module github.com/DataDog/opentelemetry-mapping-go/pkg/quantile

go 1.22

require (
	github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest v0.22.0
	github.com/DataDog/sketches-go v1.4.4
	github.com/dustin/go-humanize v1.0.1
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/DataDog/opentelemetry-mapping-go/pkg/internal/sketchtest => ../internal/sketchtest

retract v0.4.0 // see #107
