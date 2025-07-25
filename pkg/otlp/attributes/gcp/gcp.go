// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv16 "go.opentelemetry.io/otel/semconv/v1.6.1"
)

// HostInfo holds the GCP host information.
type HostInfo struct {
	HostAliases []string
	GCPTags     []string
}

// HostnameFromAttrs gets the GCP Integration hostname from attributes
// if available.
func HostnameFromAttrs(attrs pcommon.Map) (string, bool) {
	hostName, ok := attrs.Get(string(semconv16.HostNameKey))
	if !ok {
		// We need the hostname.
		return "", false
	}

	name := hostName.Str()
	if strings.Count(name, ".") >= 3 {
		// Unless the host.name attribute has been tampered with, use the same logic as the Agent to
		// extract the hostname: https://github.com/DataDog/datadog-agent/blob/7.36.0/pkg/util/cloudproviders/gce/gce.go#L106
		name = strings.SplitN(name, ".", 2)[0]
	}

	cloudAccount, ok := attrs.Get(string(semconv16.CloudAccountIDKey))
	if !ok {
		// We need the project ID.
		return "", false
	}

	alias := fmt.Sprintf("%s.%s", name, cloudAccount.Str())
	return alias, true
}

// HostInfoFromAttrs gets GCP host info from attributes following
// OpenTelemetry semantic conventions
func HostInfoFromAttrs(attrs pcommon.Map) (hostInfo *HostInfo) {
	hostInfo = &HostInfo{}

	if hostID, ok := attrs.Get(string(semconv16.HostIDKey)); ok {
		hostInfo.GCPTags = append(hostInfo.GCPTags, fmt.Sprintf("instance-id:%s", hostID.Str()))
	}

	if cloudZone, ok := attrs.Get(string(semconv16.CloudAvailabilityZoneKey)); ok {
		hostInfo.GCPTags = append(hostInfo.GCPTags, fmt.Sprintf("zone:%s", cloudZone.Str()))
	}

	if hostType, ok := attrs.Get(string(semconv16.HostTypeKey)); ok {
		hostInfo.GCPTags = append(hostInfo.GCPTags, fmt.Sprintf("instance-type:%s", hostType.Str()))
	}

	if cloudAccount, ok := attrs.Get(string(semconv16.CloudAccountIDKey)); ok {
		hostInfo.GCPTags = append(hostInfo.GCPTags, fmt.Sprintf("project:%s", cloudAccount.Str()))
	}

	return
}
