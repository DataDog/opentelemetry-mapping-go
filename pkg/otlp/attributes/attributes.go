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
	semconv127 "go.opentelemetry.io/otel/semconv/v1.27.0"
	semconv16 "go.opentelemetry.io/otel/semconv/v1.6.1"
)

// customContainerTagPrefix defines the prefix for custom container tags.
const customContainerTagPrefix = "datadog.container.tag."
const DDNamespacePrefix = "datadog."

var (
	// coreMapping defines the mapping between OpenTelemetry semantic conventions
	// and Datadog Agent conventions for env, service and version.
	coreMapping = map[string]string{
		// Datadog conventions
		// https://docs.datadoghq.com/getting_started/tagging/unified_service_tagging/
		string(semconv16.DeploymentEnvironmentKey):      KeyEnv,
		string(semconv127.ServiceNameKey):               KeyService,
		string(semconv127.ServiceVersionKey):            KeyVersion,
		string(semconv127.DeploymentEnvironmentNameKey): KeyEnv,
	}

	// ContainerMappings defines the mapping between OpenTelemetry semantic conventions
	// and Datadog Agent conventions for containers.
	ContainerMappings = map[string]string{
		// Containers
		string(semconv127.ContainerIDKey):        KeyContainerID,
		string(semconv127.ContainerNameKey):      KeyContainerName,
		string(semconv127.ContainerImageNameKey): KeyImageName,
		string(semconv16.ContainerImageTagKey):   KeyImageTag,
		string(semconv127.ContainerRuntimeKey):   KeyRuntime,

		// Cloud conventions
		// https://www.datadoghq.com/blog/tagging-best-practices/
		string(semconv127.CloudProviderKey):         KeyCloudProvider,
		string(semconv127.CloudRegionKey):           KeyRegion,
		string(semconv127.CloudAvailabilityZoneKey): KeyAvailabilityZone,

		// ECS conventions
		// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/tagger/collectors/ecs_extract.go
		string(semconv127.AWSECSTaskFamilyKey):   KeyTaskFamily,
		string(semconv127.AWSECSTaskARNKey):      KeyTaskARN,
		string(semconv127.AWSECSClusterARNKey):   KeyECSClusterName,
		string(semconv127.AWSECSTaskRevisionKey): KeyTaskVersion,
		string(semconv127.AWSECSContainerARNKey): KeyECSContainerName,

		// Kubernetes resource name (via semantic conventions)
		// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/util/kubernetes/const.go
		string(semconv127.K8SContainerNameKey):   KeyKubeContainerName,
		string(semconv127.K8SClusterNameKey):     KeyKubeClusterName,
		string(semconv127.K8SDeploymentNameKey):  KeyKubeDeployment,
		string(semconv127.K8SReplicaSetNameKey):  KeyKubeReplicaSet,
		string(semconv127.K8SStatefulSetNameKey): KeyKubeStatefulSet,
		string(semconv127.K8SDaemonSetNameKey):   KeyKubeDaemonSet,
		string(semconv127.K8SJobNameKey):         KeyKubeJob,
		string(semconv127.K8SCronJobNameKey):     KeyKubeCronJob,
		string(semconv127.K8SNamespaceNameKey):   KeyKubeNamespace,
		string(semconv127.K8SPodNameKey):         KeyPodName,
	}

	// Kubernetes mappings defines the mapping between Kubernetes conventions (both general and Datadog specific)
	// and Datadog Agent conventions. The Datadog Agent conventions can be found at
	// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/tagger/collectors/const.go and
	// https://github.com/DataDog/datadog-agent/blob/e081bed/pkg/util/kubernetes/const.go
	kubernetesMapping = map[string]string{
		// Standard Datadog labels
		"tags.datadoghq.com/env":     KeyEnv,
		"tags.datadoghq.com/service": KeyService,
		"tags.datadoghq.com/version": KeyVersion,

		// Standard Kubernetes labels
		"app.kubernetes.io/name":       KeyKubeAppName,
		"app.kubernetes.io/instance":   KeyKubeAppInstance,
		"app.kubernetes.io/version":    KeyKubeAppVersion,
		"app.kuberenetes.io/component": KeyKubeAppComponent,
		"app.kubernetes.io/part-of":    KeyKubeAppPartOf,
		"app.kubernetes.io/managed-by": KeyKubeAppManagedBy,
	}

	// Kubernetes out of the box Datadog tags
	// https://docs.datadoghq.com/containers/kubernetes/tag/?tab=containerizedagent#out-of-the-box-tags
	// https://github.com/DataDog/datadog-agent/blob/d33d042d6786e8b85f72bb627fbf06ad8a658031/comp/core/tagger/taggerimpl/collectors/workloadmeta_extract.go
	// Note: if any OTel semantics happen to overlap with these tag names, they will also be added as Datadog tags.
	kubernetesDDTags = map[string]struct{}{
		"architecture":                {},
		"availability-zone":           {},
		"chronos_job":                 {},
		"chronos_job_owner":           {},
		"cluster_name":                {},
		"container_id":                {},
		"container_name":              {},
		"dd_remote_config_id":         {},
		"dd_remote_config_rev":        {},
		"display_container_name":      {},
		"docker_image":                {},
		"ecs_cluster_name":            {},
		"ecs_container_name":          {},
		"eks_fargate_node":            {},
		"env":                         {},
		"git.commit.sha":              {},
		"git.repository_url":          {},
		"image_id":                    {},
		"image_name":                  {},
		"image_tag":                   {},
		"kube_app_component":          {},
		"kube_app_instance":           {},
		"kube_app_managed_by":         {},
		"kube_app_name":               {},
		"kube_app_part_of":            {},
		"kube_app_version":            {},
		"kube_container_name":         {},
		"kube_cronjob":                {},
		"kube_daemon_set":             {},
		"kube_deployment":             {},
		"kube_job":                    {},
		"kube_namespace":              {},
		"kube_ownerref_kind":          {},
		"kube_ownerref_name":          {},
		"kube_priority_class":         {},
		"kube_qos":                    {},
		"kube_replica_set":            {},
		"kube_replication_controller": {},
		"kube_service":                {},
		"kube_stateful_set":           {},
		"language":                    {},
		"marathon_app":                {},
		"mesos_task":                  {},
		"nomad_dc":                    {},
		"nomad_group":                 {},
		"nomad_job":                   {},
		"nomad_namespace":             {},
		"nomad_task":                  {},
		"oshift_deployment":           {},
		"oshift_deployment_config":    {},
		"os_name":                     {},
		"os_version":                  {},
		"persistentvolumeclaim":       {},
		"pod_name":                    {},
		"pod_phase":                   {},
		"rancher_container":           {},
		"rancher_service":             {},
		"rancher_stack":               {},
		"region":                      {},
		"service":                     {},
		"short_image":                 {},
		"swarm_namespace":             {},
		"swarm_service":               {},
		"task_name":                   {},
		"task_family":                 {},
		"task_version":                {},
		"task_arn":                    {},
		"version":                     {},
	}

	// HTTPMappings defines the mapping between OpenTelemetry semantic conventions
	// and Datadog Agent conventions for HTTP attributes.
	HTTPMappings = map[string]string{
		string(semconv127.ClientAddressKey):          "http.client_ip",
		string(semconv127.HTTPResponseBodySizeKey):   "http.response.content_length",
		string(semconv127.HTTPResponseStatusCodeKey): "http.status_code",
		string(semconv127.HTTPRequestBodySizeKey):    "http.request.content_length",
		"http.request.header.referrer":               "http.referrer",
		string(semconv127.HTTPRequestMethodKey):      "http.method",
		string(semconv127.HTTPRouteKey):              "http.route",
		string(semconv127.NetworkProtocolVersionKey): "http.version",
		string(semconv127.ServerAddressKey):          "http.server_name",
		string(semconv127.URLFullKey):                "http.url",
		string(semconv127.UserAgentOriginalKey):      "http.useragent",
	}

	KeyDatadogHostname              = DDNamespacePrefix + "host.name"
	KeyDatadogProcessExecutableName = DDNamespacePrefix + "process.executable.name"
	KeyDatadogProcessExecutablePath = DDNamespacePrefix + "process.executable.path"
	KeyDatadogProcessCommand        = DDNamespacePrefix + "process.command"
	KeyDatadogProcessCommandLine    = DDNamespacePrefix + "process.command_line"
	KeyDatadogProcessPID            = DDNamespacePrefix + "process.pid"
	KeyDatadogProcessOwner          = DDNamespacePrefix + "process.owner"
	KeyDatadogOSType                = DDNamespacePrefix + "ostype"

	KeyDatadogOriginID         = DDNamespacePrefix + "origin.id"
	KeyDatadogSourceKind       = DDNamespacePrefix + "source.kind"
	KeyDatadogSourceIdentifier = DDNamespacePrefix + "source.identifier"
)

const (
	KeyEnv     = "env"
	KeyService = "service"
	KeyVersion = "version"

	KeyContainerID       = "container_id"
	KeyContainerName     = "container_name"
	KeyImageName         = "image_name"
	KeyImageTag          = "image_tag"
	KeyRuntime           = "runtime"
	KeyCloudProvider     = "cloud_provider"
	KeyRegion            = "region"
	KeyAvailabilityZone  = "zone"
	KeyTaskFamily        = "task_family"
	KeyTaskARN           = "task_arn"
	KeyTaskVersion       = "task_version"
	KeyECSClusterName    = "ecs_cluster_name"
	KeyECSContainerName  = "ecs_container_name"
	KeyKubeContainerName = "kube_container_name"
	KeyKubeClusterName   = "kube_cluster_name"
	KeyKubeDeployment    = "kube_deployment"
	KeyKubeReplicaSet    = "kube_replica_set"
	KeyKubeStatefulSet   = "kube_stateful_set"
	KeyKubeDaemonSet     = "kube_daemon_set"
	KeyKubeJob           = "kube_job"
	KeyKubeCronJob       = "kube_cronjob"
	KeyKubeNamespace     = "kube_namespace"
	KeyPodName           = "pod_name"

	KeyKubeAppName      = "kube_app_name"
	KeyKubeAppInstance  = "kube_app_instance"
	KeyKubeAppVersion   = "kube_app_version"
	KeyKubeAppComponent = "kube_app_component"
	KeyKubeAppPartOf    = "kube_app_part_of"
	KeyKubeAppManagedBy = "kube_app_managed_by"
)

var keysToCheckForInDDNamespace = map[string]struct{}{
	KeyEnv:               {},
	KeyService:           {},
	KeyVersion:           {},
	KeyContainerID:       {},
	KeyContainerName:     {},
	KeyImageName:         {},
	KeyImageTag:          {},
	KeyRuntime:           {},
	KeyCloudProvider:     {},
	KeyRegion:            {},
	KeyAvailabilityZone:  {},
	KeyTaskFamily:        {},
	KeyTaskARN:           {},
	KeyTaskVersion:       {},
	KeyECSClusterName:    {},
	KeyECSContainerName:  {},
	KeyKubeContainerName: {},
	KeyKubeClusterName:   {},
	KeyKubeDeployment:    {},
	KeyKubeReplicaSet:    {},
	KeyKubeStatefulSet:   {},
	KeyKubeDaemonSet:     {},
	KeyKubeJob:           {},
	KeyKubeCronJob:       {},
	KeyKubeNamespace:     {},
	KeyPodName:           {},
	KeyKubeAppName:       {},
	KeyKubeAppInstance:   {},
	KeyKubeAppVersion:    {},
	KeyKubeAppComponent:  {},
	KeyKubeAppPartOf:     {},
	KeyKubeAppManagedBy:  {},
}

func MergeTagMaps(signalTagsMap, resourceTagsMap map[string]string, ignoreMissingDatadogFields bool) map[string]string {
	tagsMap := make(map[string]string, len(signalTagsMap)+len(resourceTagsMap))

	for key, val := range resourceTagsMap {
		tagsMap[key] = val
	}
	// signal tags take precedence over resource tags
	for key, val := range signalTagsMap {
		tagsMap[key] = val
	}

	// Only keep the highest-precedence process key in tagsMap
	processKeys := []string{
		string(semconv127.ProcessExecutableNameKey),
		string(semconv127.ProcessExecutablePathKey),
		string(semconv127.ProcessCommandKey),
		string(semconv127.ProcessCommandLineKey),
	}
	for i, k := range processKeys {
		if v := tagsMap[k]; v != "" {
			// Delete all lower-precedence keys
			for _, lowerK := range processKeys[i+1:] {
				delete(tagsMap, lowerK)
			}
			break
		}
	}

	return tagsMap
}

// GetTagsFromAttributesPreferringDatadogNamespace converts a selected list of attributes
// to a tag list that can be added to metrics. It follows this order of precedence:
// 1. datadog.* span attributes
// 2. datadog.* resource attributes
// 3. standard span attributes
// 4. standard resource attributes
// If ignoreMissingDatadogFields is true, it will not add tags that are not present in the Datadog namespace.
func GetTagsFromAttributesPreferringDatadogNamespace(attrs pcommon.Map, ignoreMissingDatadogFields bool) map[string]string {
	tagsMap := make(map[string]string, attrs.Len())

	var processAttributes processAttributes
	var systemAttributes systemAttributes

	attrs.Range(func(key string, value pcommon.Value) bool {
		switch key {
		// Process attributes
		case KeyDatadogProcessExecutableName:
			if processAttributes.ExecutableName == "" {
				processAttributes.ExecutableName = value.Str()
			}
		case KeyDatadogProcessExecutablePath:
			if processAttributes.ExecutablePath == "" {
				processAttributes.ExecutablePath = value.Str()
			}
		case KeyDatadogProcessCommand:
			if processAttributes.Command == "" {
				processAttributes.Command = value.Str()
			}
		case KeyDatadogProcessCommandLine:
			if processAttributes.CommandLine == "" {
				processAttributes.CommandLine = value.Str()
			}
		case KeyDatadogProcessPID:
			if processAttributes.PID == 0 {
				processAttributes.PID = value.Int()
			}
		case KeyDatadogProcessOwner:
			if processAttributes.Owner == "" {
				processAttributes.Owner = value.Str()
			}

		// System attributes
		case KeyDatadogOSType:
			if systemAttributes.OSType == "" {
				systemAttributes.OSType = value.Str()
			}
		}

		if strings.HasPrefix(key, customContainerTagPrefix) {
			// Custom container tags are checked after all other semantic conventions
			return true
		}

		if strings.HasPrefix(key, DDNamespacePrefix) {
			key = strings.TrimPrefix(key, DDNamespacePrefix)
			// Kubernetes DD tags
			_, found1 := kubernetesDDTags[key]
			_, found2 := keysToCheckForInDDNamespace[key]
			if found1 || found2 {
				if tagsMap[key] == "" {
					tagsMap[key] = value.Str()
				}
				return true
			}
		}

		return true
	})

	if !ignoreMissingDatadogFields {
		attrs.Range(func(key string, value pcommon.Value) bool {
			switch key {
			// Process attributes
			case string(semconv127.ProcessExecutableNameKey):
				if processAttributes.ExecutableName == "" {
					processAttributes.ExecutableName = value.Str()
				}
			case string(semconv127.ProcessExecutablePathKey):
				if processAttributes.ExecutablePath == "" {
					processAttributes.ExecutablePath = value.Str()
				}
			case string(semconv127.ProcessCommandKey):
				if processAttributes.Command == "" {
					processAttributes.Command = value.Str()
				}
			case string(semconv127.ProcessCommandLineKey):
				if processAttributes.CommandLine == "" {
					processAttributes.CommandLine = value.Str()
				}
			case string(semconv127.ProcessPIDKey):
				if processAttributes.PID == 0 {
					processAttributes.PID = value.Int()
				}
			case string(semconv127.ProcessOwnerKey):
				if processAttributes.Owner == "" {
					processAttributes.Owner = value.Str()
				}

			// System attributes
			case string(semconv127.OSTypeKey):
				if systemAttributes.OSType == "" {
					systemAttributes.OSType = value.Str()
				}
			}

			// core attributes mapping
			if datadogKey, found := coreMapping[key]; found && value.Str() != "" {
				if tagsMap[datadogKey] == "" {
					tagsMap[datadogKey] = value.Str()
				}
			}

			// Kubernetes labels mapping
			if datadogKey, found := kubernetesMapping[key]; found && value.Str() != "" {
				if tagsMap[datadogKey] == "" {
					tagsMap[datadogKey] = value.Str()
				}
			}

			// Kubernetes DD tags
			if _, found := kubernetesDDTags[key]; found {
				if tagsMap[key] == "" {
					tagsMap[key] = value.Str()
				}
			}

			// Container Tag mappings
			ctags := ContainerTagsFromResourceAttributes(attrs)
			for key, val := range ctags {
				if tagsMap[key] == "" {
					tagsMap[key] = val
				}
			}
			return true
		})
	}

	for k, v := range processAttributes.extractTags() {
		if tagsMap[k] == "" {
			tagsMap[k] = v
		}
	}
	for k, v := range systemAttributes.extractTags() {
		if tagsMap[k] == "" {
			tagsMap[k] = v
		}
	}

	return tagsMap
}

// TagsFromAttributes converts a selected list of attributes
// to a tag list that can be added to metrics.
// Deprecated: Use GetTagsFromAttributesPreferringDatadogNamespace instead.
func TagsFromAttributes(attrs pcommon.Map) []string {
	tags := make([]string, 0, attrs.Len())

	var pAttributes processAttributes
	var sAttributes systemAttributes

	attrs.Range(func(key string, value pcommon.Value) bool {
		switch key {
		// Process attributes
		case string(semconv127.ProcessExecutableNameKey):
			pAttributes.ExecutableName = value.Str()
		case string(semconv127.ProcessExecutablePathKey):
			pAttributes.ExecutablePath = value.Str()
		case string(semconv127.ProcessCommandKey):
			pAttributes.Command = value.Str()
		case string(semconv127.ProcessCommandLineKey):
			pAttributes.CommandLine = value.Str()
		case string(semconv127.ProcessPIDKey):
			pAttributes.PID = value.Int()
		case string(semconv127.ProcessOwnerKey):
			pAttributes.Owner = value.Str()

		// System attributes
		case string(semconv127.OSTypeKey):
			sAttributes.OSType = value.Str()
		}

		// core attributes mapping
		if datadogKey, found := coreMapping[key]; found && value.Str() != "" {
			tags = append(tags, fmt.Sprintf("%s:%s", datadogKey, value.Str()))
		}

		// Kubernetes labels mapping
		if datadogKey, found := kubernetesMapping[key]; found && value.Str() != "" {
			tags = append(tags, fmt.Sprintf("%s:%s", datadogKey, value.Str()))
		}

		// Kubernetes DD tags
		if _, found := kubernetesDDTags[key]; found {
			tags = append(tags, fmt.Sprintf("%s:%s", key, value.Str()))
		}
		return true
	})

	// Container Tag mappings
	ctags := ContainerTagsFromResourceAttributes(attrs)
	for key, val := range ctags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, val))
	}

	// Convert process and system attribute maps to tag strings
	for key, val := range pAttributes.extractTags() {
		tags = append(tags, fmt.Sprintf("%s:%s", key, val))
	}
	for key, val := range sAttributes.extractTags() {
		tags = append(tags, fmt.Sprintf("%s:%s", key, val))
	}

	return tags
}

// OriginIDFromAttributes gets the origin IDs from resource attributes.
// If not found, an empty string is returned for each of them.
func OriginIDFromAttributes(attrs pcommon.Map) (originID string) {
	// originID is always empty. Container ID is preferred over Kubernetes pod UID.
	// Prefixes come from pkg/util/kubernetes/kubelet and pkg/util/containers.
	if containerID, ok := attrs.Get(string(semconv16.ContainerIDKey)); ok {
		originID = "container_id://" + containerID.AsString()
	} else if podUID, ok := attrs.Get(string(semconv16.K8SPodUIDKey)); ok {
		originID = "kubernetes_pod_uid://" + podUID.AsString()
	}
	return
}

// ContainerTagFromResourceAttributes extracts container tags from the given
// set of resource attributes. Container tags are extracted via semantic
// conventions. Customer container tags are extracted via resource attributes
// prefixed by datadog.container.tag. Custom container tag values of a different type
// than ValueTypeStr will be ignored.
// In the case of duplicates between semantic conventions and custom resource attributes
// (e.g. container.id, datadog.container.tag.container_id) the semantic convention takes
// precedence.
func ContainerTagsFromResourceAttributes(attrs pcommon.Map) map[string]string {
	ddtags := make(map[string]string)
	attrs.Range(func(key string, value pcommon.Value) bool {
		// Semantic Conventions
		if datadogKey, found := ContainerMappings[key]; found && value.Str() != "" {
			ddtags[datadogKey] = value.Str()
		}
		// Custom (datadog.container.tag namespace)
		if strings.HasPrefix(key, customContainerTagPrefix) {
			customKey := strings.TrimPrefix(key, customContainerTagPrefix)
			if customKey != "" && value.Str() != "" {
				// Do not replace if set via semantic conventions mappings.
				if _, found := ddtags[customKey]; !found {
					ddtags[customKey] = value.Str()
				}
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
		datadogKey, found := ContainerMappings[key]
		if !found {
			continue
		}
		ddtags[datadogKey] = val
	}
	return ddtags
}
