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
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

// customContainerTagPrefix defines the prefix for custom container tags.
const customContainerTagPrefix = "datadog.container.tag."

var (
	// coreMapping defines the mapping between OpenTelemetry semantic conventions
	// and Datadog Agent conventions for env, service and version.
	coreMapping = map[string]string{
		// Datadog conventions
		// https://docs.datadoghq.com/getting_started/tagging/unified_service_tagging/
		conventions.AttributeDeploymentEnvironment: "env",
		conventions.AttributeServiceName:           "service",
		conventions.AttributeServiceVersion:        "version",
	}

	// containerMappings defines the mapping between OpenTelemetry semantic conventions
	// and Datadog Agent conventions for containers.
	containerMappings = map[string]string{
		// Containers
		conventions.AttributeContainerID:        "container_id",
		conventions.AttributeContainerName:      "container_name",
		conventions.AttributeContainerImageName: "image_name",
		conventions.AttributeContainerImageTag:  "image_tag",
		conventions.AttributeContainerRuntime:   "runtime",

		// Cloud conventions
		// https://www.datadoghq.com/blog/tagging-best-practices/
		conventions.AttributeCloudProvider:         "cloud_provider",
		conventions.AttributeCloudRegion:           "region",
		conventions.AttributeCloudAvailabilityZone: "zone",

		// ECS conventions
		// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/tagger/collectors/ecs_extract.go
		conventions.AttributeAWSECSTaskFamily:   "task_family",
		conventions.AttributeAWSECSTaskARN:      "task_arn",
		conventions.AttributeAWSECSClusterARN:   "ecs_cluster_name",
		conventions.AttributeAWSECSTaskRevision: "task_version",
		conventions.AttributeAWSECSContainerARN: "ecs_container_name",

		// Kubernetes resource name (via semantic conventions)
		// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/util/kubernetes/const.go
		conventions.AttributeK8SContainerName:   "kube_container_name",
		conventions.AttributeK8SClusterName:     "kube_cluster_name",
		conventions.AttributeK8SDeploymentName:  "kube_deployment",
		conventions.AttributeK8SReplicaSetName:  "kube_replica_set",
		conventions.AttributeK8SStatefulSetName: "kube_stateful_set",
		conventions.AttributeK8SDaemonSetName:   "kube_daemon_set",
		conventions.AttributeK8SJobName:         "kube_job",
		conventions.AttributeK8SCronJobName:     "kube_cronjob",
		conventions.AttributeK8SNamespaceName:   "kube_namespace",
		conventions.AttributeK8SPodName:         "pod_name",
	}

	// Kubernetes mappings defines the mapping between Kubernetes conventions (both general and Datadog specific)
	// and Datadog Agent conventions. The Datadog Agent conventions can be found at
	// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/tagger/collectors/const.go and
	// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/util/kubernetes/const.go
	kubernetesMapping = map[string]string{
		// Standard Datadog labels
		"tags.datadoghq.com/env":     "env",
		"tags.datadoghq.com/service": "service",
		"tags.datadoghq.com/version": "version",

		// Standard Kubernetes labels
		"app.kubernetes.io/name":       "kube_app_name",
		"app.kubernetes.io/instance":   "kube_app_instance",
		"app.kubernetes.io/version":    "kube_app_version",
		"app.kuberenetes.io/component": "kube_app_component",
		"app.kubernetes.io/part-of":    "kube_app_part_of",
		"app.kubernetes.io/managed-by": "kube_app_managed_by",
	}
)

// TagsFromAttributes converts a selected list of attributes
// to a tag list that can be added to metrics.
func TagsFromAttributes(attrs pcommon.Map) []string {
	tags := make([]string, 0, attrs.Len())

	var processAttributes processAttributes
	var systemAttributes systemAttributes

	attrs.Range(func(key string, value pcommon.Value) bool {
		switch key {
		// Process attributes
		case conventions.AttributeProcessExecutableName:
			processAttributes.ExecutableName = value.Str()
		case conventions.AttributeProcessExecutablePath:
			processAttributes.ExecutablePath = value.Str()
		case conventions.AttributeProcessCommand:
			processAttributes.Command = value.Str()
		case conventions.AttributeProcessCommandLine:
			processAttributes.CommandLine = value.Str()
		case conventions.AttributeProcessPID:
			processAttributes.PID = value.Int()
		case conventions.AttributeProcessOwner:
			processAttributes.Owner = value.Str()

		// System attributes
		case conventions.AttributeOSType:
			systemAttributes.OSType = value.Str()
		}

		// core attributes mapping
		if datadogKey, found := coreMapping[key]; found && value.Str() != "" {
			tags = append(tags, fmt.Sprintf("%s:%s", datadogKey, value.Str()))
		}

		// Kubernetes labels mapping
		if datadogKey, found := kubernetesMapping[key]; found && value.Str() != "" {
			tags = append(tags, fmt.Sprintf("%s:%s", datadogKey, value.Str()))
		}
		return true
	})

	// Container Tag mappings
	ctags := ContainerTagsFromResourceAttributes(attrs)
	for key, val := range ctags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, val))
	}

	tags = append(tags, processAttributes.extractTags()...)
	tags = append(tags, systemAttributes.extractTags()...)

	return tags
}

// OriginIDFromAttributes gets the origin IDs from resource attributes.
// If not found, an empty string is returned for each of them.
func OriginIDFromAttributes(attrs pcommon.Map) (originID string) {
	// originID is always empty. Container ID is preferred over Kubernetes pod UID.
	// Prefixes come from pkg/util/kubernetes/kubelet and pkg/util/containers.
	if containerID, ok := attrs.Get(conventions.AttributeContainerID); ok {
		originID = "container_id://" + containerID.AsString()
	} else if podUID, ok := attrs.Get(conventions.AttributeK8SPodUID); ok {
		originID = "kubernetes_pod_uid://" + podUID.AsString()
	}
	return
}

// ContainerTagFromResourceAttributes extracts container tags from the given
// set of resource attributes.
func ContainerTagsFromResourceAttributes(attrs pcommon.Map) map[string]string {
	ddtags := make(map[string]string)

	attrs.Range(func(key string, value pcommon.Value) bool {
		// Semantic Conventions
		if datadogKey, found := containerMappings[key]; found && value.Str() != "" {
			ddtags[datadogKey] = value.AsString()
		}
		// Custom (datadog.container.tag namespace)
		if strings.HasPrefix(key, customContainerTagPrefix) {
			if strings.TrimPrefix(key, customContainerTagPrefix) != "" && value.Str() != "" {
				ddtags[strings.TrimPrefix(key, customContainerTagPrefix)] = value.AsString()
			}
		}
		return true
	})

	return ddtags
}

// ContainerTagFromAttributes extracts the value of _dd.tags.container from the given
// set of attributes.
// Deprecated: Deprecated in favor of ContainerTagFromResourceAttributes.
func ContainerTagFromAttributes(attr map[string]string) map[string]string {
	ddtags := make(map[string]string)
	for key, val := range attr {
		datadogKey, found := containerMappings[key]
		if !found {
			continue
		}
		ddtags[datadogKey] = val
	}
	return ddtags
}
