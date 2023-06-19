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

// gohai fieldPlatform subsection fields.
const (
	fieldPlatformHostname         = "hostname"
	fieldPlatformOS               = "os"
	fieldPlatformProcessor        = "processor"
	fieldPlatformMachine          = "machine"
	fieldPlatformHardwarePlatform = "hardware_fieldPlatform"
	fieldPlatformKernelName       = "kernel_name"
	fieldPlatformKernelRelease    = "kernel_release"
	fieldPlatformKernelVersion    = "kernel_version"
)

// platformAttributesMap defines the mapping between Gohai fieldPlatform fields
// and resource attribute names (semantic conventions or not).
var platformAttributesMap map[string]string = map[string]string{
	fieldPlatformOS:               conventions.AttributeOSDescription,
	fieldPlatformProcessor:        conventions.AttributeHostArch,
	fieldPlatformMachine:          conventions.AttributeHostArch,
	fieldPlatformHardwarePlatform: conventions.AttributeHostArch,
	fieldPlatformKernelName:       attributeKernelName,
	fieldPlatformKernelRelease:    attributeKernelRelease,
	fieldPlatformKernelVersion:    attributeKernelVersion,
}
