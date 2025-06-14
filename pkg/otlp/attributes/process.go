// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attributes

import (
	"fmt"

	semconv16 "go.opentelemetry.io/otel/semconv/v1.6.1"
)

type processAttributes struct {
	ExecutableName string
	ExecutablePath string
	Command        string
	CommandLine    string
	PID            int64
	Owner          string
}

func (pattrs *processAttributes) extractTags() []string {
	tags := make([]string, 0, 1)

	// According to OTel conventions: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/process.md,
	// a process can be defined by any of the 4 following attributes: process.executable.name, process.executable.path, process.command or process.command_line
	// (process.command_args isn't in the current attribute conventions: https://github.com/open-telemetry/opentelemetry-collector/blob/ecb27f49d4e26ae42d82e6ea18d57b08e252452d/model/semconv/opentelemetry.go#L58-L63)
	// We go through them, and add the first available one as a tag to identify the process.
	// We don't want to add all of them to avoid unnecessarily increasing the number of tags attached to a metric.

	// TODO: check if this order should be changed.
	switch {
	case pattrs.ExecutableName != "": // otelcol
		tags = append(tags, fmt.Sprintf("%s:%s", string(semconv16.ProcessExecutableNameKey), pattrs.ExecutableName))
	case pattrs.ExecutablePath != "": // /usr/bin/cmd/otelcol
		tags = append(tags, fmt.Sprintf("%s:%s", string(semconv16.ProcessExecutablePathKey), pattrs.ExecutablePath))
	case pattrs.Command != "": // cmd/otelcol
		tags = append(tags, fmt.Sprintf("%s:%s", string(semconv16.ProcessCommandKey), pattrs.Command))
	case pattrs.CommandLine != "": // cmd/otelcol --config="/path/to/config.yaml"
		tags = append(tags, fmt.Sprintf("%s:%s", string(semconv16.ProcessCommandLineKey), pattrs.CommandLine))
	}

	// For now, we don't care about the process ID nor the process owner.

	return tags
}
