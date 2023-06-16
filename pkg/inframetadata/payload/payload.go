// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Original sources for this file:
// - https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/a5cdd2/exporter/datadogexporter/internal/hostmetadata/metadata.go
// - https://github.com/DataDog/datadog-agent/blob/ab37437/pkg/metadata/host/payload.go
//
// This file defines the host metadata payload. This payload fills in information about the hosts in the Datadog's infrastructure list page.

package payload

import (
	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/gohai"
)

// HostMetadata includes metadata about the host tags,
// host aliases and identifies the host as an OpenTelemetry host
type HostMetadata struct {
	// Meta includes metadata about the host.
	Meta *Meta `json:"meta"`

	// InternalHostname is the canonical hostname
	InternalHostname string `json:"internalHostname"`

	// Version is the OpenTelemetry Collector version.
	// This is used for correctly identifying the Collector in the backend,
	// and for telemetry purposes.
	Version string `json:"otel_version"`

	// Flavor is always set to "opentelemetry-collector".
	// It is used for telemetry purposes in the backend.
	Flavor string `json:"agent-flavor"`

	// Tags includes the host tags
	Tags *HostTags `json:"host-tags"`

	// Payload contains inventory of system information provided by gohai
	// this is embedded because of special serialization requirements
	// the field `gohai` is JSON-formatted string
	gohai.Payload

	// Processes contains the process payload devired by gohai
	// Because of legacy reasons this is called resources in datadog intake
	Processes *gohai.ProcessesPayload `json:"resources"`
}

// HostTags are the host tags.
// Currently only system (configuration) tags are considered.
type HostTags struct {
	// OTel are host tags set in the configuration
	OTel []string `json:"otel,omitempty"`

	// GCP are Google Cloud Platform tags
	GCP []string `json:"google cloud platform,omitempty"`
}

// Meta includes metadata about the host aliases
type Meta struct {
	// InstanceID is the EC2 instance id the Collector is running on, if available
	InstanceID string `json:"instance-id,omitempty"`

	// EC2Hostname is the hostname from the EC2 metadata API
	EC2Hostname string `json:"ec2-hostname,omitempty"`

	// Hostname is the canonical hostname
	Hostname string `json:"hostname"`

	// SocketHostname is the OS hostname
	SocketHostname string `json:"socket-hostname,omitempty"`

	// SocketFqdn is the FQDN hostname
	SocketFqdn string `json:"socket-fqdn,omitempty"`

	// HostAliases are other available host names
	HostAliases []string `json:"host_aliases,omitempty"`
}
