package hostmap

import (
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
)

// Custom attributes not in the OpenTelemetry specification
const (
	attributeKernelName    = "os.kernel.name"
	attributeKernelRelease = "os.kernel.release"
	attributeKernelVersion = "os.kernel.version"
)

// gohai platform subsection fields.
const (
	platformHostname         = "hostname"
	platformOS               = "os"
	platformProcessor        = "processor"
	platformMachine          = "machine"
	platformHardwarePlatform = "hardware_platform"
	platformKernelName       = "kernel_name"
	platformKernelRelease    = "kernel_release"
	platformKernelVersion    = "kernel_version"
)

// platformAttributesMap defines the mapping between Gohai platform fields
// and resource attribute names (semantic conventions or not).
var platformAttributesMap map[string]string = map[string]string{
	platformOS:               conventions.AttributeOSDescription,
	platformProcessor:        conventions.AttributeHostArch,
	platformMachine:          conventions.AttributeHostArch,
	platformHardwarePlatform: conventions.AttributeHostArch,
	platformKernelName:       attributeKernelName,
	platformKernelRelease:    attributeKernelRelease,
	platformKernelVersion:    attributeKernelVersion,
}
