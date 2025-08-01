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

package azure

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv16 "go.opentelemetry.io/otel/semconv/v1.6.1"
)

const (
	// AttributeResourceGroupName is the Azure resource group name attribute
	AttributeResourceGroupName = "azure.resourcegroup.name"
)

// HostInfo has the Azure host information
type HostInfo struct {
	HostAliases []string
}

// HostnameFromAttrs gets the Azure hostname from attributes.
func HostnameFromAttrs(attrs pcommon.Map) (string, bool) {
	if vmID, ok := attrs.Get(string(semconv16.HostIDKey)); ok {
		return vmID.Str(), true
	}

	if hostname, ok := attrs.Get(string(semconv16.HostNameKey)); ok {
		return hostname.Str(), true
	}

	return "", false
}

// ClusterNameFromAttributes gets the Azure cluster name from attributes
func ClusterNameFromAttributes(attrs pcommon.Map) (string, bool) {
	// Get cluster name from resource group from pkg/util/cloudprovider/azure:GetClusterName
	if resourceGroup, ok := attrs.Get(AttributeResourceGroupName); ok {
		splitAll := strings.Split(resourceGroup.Str(), "_")
		if len(splitAll) < 4 || strings.ToLower(splitAll[0]) != "mc" {
			return "", false // Failed to parse
		}
		return splitAll[len(splitAll)-2], true
	}

	return "", false
}
