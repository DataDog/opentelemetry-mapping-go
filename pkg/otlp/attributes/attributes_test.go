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
	"testing"

	semconv127 "go.opentelemetry.io/collector/semconv/v1.27.0"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

func TestTagsFromAttributes(t *testing.T) {
	attributeMap := map[string]interface{}{
		semconv127.AttributeProcessExecutableName:  "otelcol",
		semconv127.AttributeProcessExecutablePath:  "/usr/bin/cmd/otelcol",
		semconv127.AttributeProcessCommand:         "cmd/otelcol",
		semconv127.AttributeProcessCommandLine:     "cmd/otelcol --config=\"/path/to/config.yaml\"",
		semconv127.AttributeProcessPID:             1,
		semconv127.AttributeProcessOwner:           "root",
		semconv127.AttributeOSType:                 "linux",
		semconv127.AttributeK8SDaemonSetName:       "daemon_set_name",
		semconv127.AttributeAWSECSClusterARN:       "cluster_arn",
		semconv127.AttributeContainerRuntime:       "cro",
		"tags.datadoghq.com/service":               "service_name",
		conventions.AttributeDeploymentEnvironment: "prod",
		semconv127.AttributeContainerName:          "custom",
		"datadog.container.tag.custom.team":        "otel",
		"kube_cronjob":                             "cron",
	}
	attrs := pcommon.NewMap()
	attrs.FromRaw(attributeMap)

	assert.ElementsMatch(t, []string{
		fmt.Sprintf("%s:%s", semconv127.AttributeProcessExecutableName, "otelcol"),
		fmt.Sprintf("%s:%s", semconv127.AttributeOSType, "linux"),
		fmt.Sprintf("%s:%s", "kube_daemon_set", "daemon_set_name"),
		fmt.Sprintf("%s:%s", "ecs_cluster_name", "cluster_arn"),
		fmt.Sprintf("%s:%s", "service", "service_name"),
		fmt.Sprintf("%s:%s", "runtime", "cro"),
		fmt.Sprintf("%s:%s", "env", "prod"),
		fmt.Sprintf("%s:%s", "container_name", "custom"),
		fmt.Sprintf("%s:%s", "custom.team", "otel"),
		fmt.Sprintf("%s:%s", "kube_cronjob", "cron"),
	}, TagsFromAttributes(attrs))
}

func TestNewDeploymentEnvironmentNameConvention_TagsFromAttributes(t *testing.T) {
	attrs := pcommon.NewMap()
	attrs.PutStr("deployment.environment.name", "staging")

	assert.Equal(t, []string{"env:staging"}, TagsFromAttributes(attrs))
}

func TestTagsFromAttributesEmpty_TagsFromAttributes(t *testing.T) {
	attrs := pcommon.NewMap()

	assert.Equal(t, []string{}, TagsFromAttributes(attrs))
}

func TestGetTagsFromAttributesPreferringDatadogNamespace(t *testing.T) {
	attributeMap := map[string]interface{}{
		semconv127.AttributeProcessExecutableName:  "otelcol",
		semconv127.AttributeProcessExecutablePath:  "/usr/bin/cmd/otelcol",
		semconv127.AttributeProcessCommand:         "cmd/otelcol",
		semconv127.AttributeProcessCommandLine:     "cmd/otelcol --config=\"/path/to/config.yaml\"",
		semconv127.AttributeProcessPID:             1,
		semconv127.AttributeProcessOwner:           "root",
		semconv127.AttributeOSType:                 "linux",
		semconv127.AttributeK8SDaemonSetName:       "daemon_set_name",
		semconv127.AttributeAWSECSClusterARN:       "cluster_arn",
		semconv127.AttributeContainerRuntime:       "cro",
		"tags.datadoghq.com/service":               "service_name",
		conventions.AttributeDeploymentEnvironment: "prod",
		semconv127.AttributeContainerName:          "custom",
		"datadog.container.tag.custom.team":        "otel",
		"kube_cronjob":                             "cron",
	}
	attrs := pcommon.NewMap()
	attrs.FromRaw(attributeMap)

	expected := map[string]string{
		semconv127.AttributeProcessExecutableName: "otelcol",
		conventions.AttributeOSType:               "linux",
		"kube_daemon_set":                         "daemon_set_name",
		"ecs_cluster_name":                        "cluster_arn",
		"service":                                 "service_name",
		"runtime":                                 "cro",
		"env":                                     "prod",
		"container_name":                          "custom",
		"custom.team":                             "otel",
		"kube_cronjob":                            "cron",
	}
	assert.Equal(t, expected, GetTagsFromAttributesPreferringDatadogNamespace(attrs, false))
}

func TestNewDeploymentEnvironmentNameConvention(t *testing.T) {
	attrs := pcommon.NewMap()
	attrs.PutStr("deployment.environment.name", "staging")

	expected := map[string]string{"env": "staging"}
	assert.Equal(t, expected, GetTagsFromAttributesPreferringDatadogNamespace(attrs, false))
}

func TestTagsFromAttributesEmpty(t *testing.T) {
	attrs := pcommon.NewMap()
	assert.Equal(t, map[string]string{}, GetTagsFromAttributesPreferringDatadogNamespace(attrs, false))
}

func TestContainerTagFromResourceAttributes(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		attributes := pcommon.NewMap()
		err := attributes.FromRaw(map[string]interface{}{
			semconv127.AttributeContainerName:         "sample_app",
			conventions.AttributeContainerImageTag:    "sample_app_image_tag",
			semconv127.AttributeContainerRuntime:      "cro",
			semconv127.AttributeK8SContainerName:      "kube_sample_app",
			semconv127.AttributeK8SReplicaSetName:     "sample_replica_set",
			semconv127.AttributeK8SDaemonSetName:      "sample_daemonset_name",
			semconv127.AttributeK8SPodName:            "sample_pod_name",
			semconv127.AttributeCloudProvider:         "sample_cloud_provider",
			semconv127.AttributeCloudRegion:           "sample_region",
			semconv127.AttributeCloudAvailabilityZone: "sample_zone",
			semconv127.AttributeAWSECSTaskFamily:      "sample_task_family",
			semconv127.AttributeAWSECSClusterARN:      "sample_ecs_cluster_name",
			semconv127.AttributeAWSECSContainerARN:    "sample_ecs_container_name",
			"datadog.container.tag.custom.team":       "otel",
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"container_name":      "sample_app",
			"image_tag":           "sample_app_image_tag",
			"runtime":             "cro",
			"kube_container_name": "kube_sample_app",
			"kube_replica_set":    "sample_replica_set",
			"kube_daemon_set":     "sample_daemonset_name",
			"pod_name":            "sample_pod_name",
			"cloud_provider":      "sample_cloud_provider",
			"region":              "sample_region",
			"zone":                "sample_zone",
			"task_family":         "sample_task_family",
			"ecs_cluster_name":    "sample_ecs_cluster_name",
			"ecs_container_name":  "sample_ecs_container_name",
			"custom.team":         "otel",
		}, ContainerTagsFromResourceAttributes(attributes))
		fmt.Println(ContainerTagsFromResourceAttributes(attributes))
	})

	t.Run("conventions vs custom", func(t *testing.T) {
		attributes := pcommon.NewMap()
		err := attributes.FromRaw(map[string]interface{}{
			semconv127.AttributeContainerName:      "ok",
			"datadog.container.tag.container_name": "nok",
		})
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"container_name": "ok",
		}, ContainerTagsFromResourceAttributes(attributes))
	})

	t.Run("invalid", func(t *testing.T) {
		attributes := pcommon.NewMap()
		err := attributes.FromRaw(map[string]interface{}{
			"empty_string_val": "",
			"":                 "empty_string_key",
			"custom_tag":       "example_custom_tag",
		})
		assert.NoError(t, err)
		slice := attributes.PutEmptySlice("datadog.container.tag.slice")
		slice.AppendEmpty().SetStr("value1")
		slice.AppendEmpty().SetStr("value2")
		assert.Equal(t, map[string]string{}, ContainerTagsFromResourceAttributes(attributes))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Empty(t, ContainerTagsFromResourceAttributes(pcommon.NewMap()))
	})
}

func TestContainerTagFromAttributes(t *testing.T) {
	attributeMap := map[string]string{
		semconv127.AttributeContainerName:         "sample_app",
		conventions.AttributeContainerImageTag:    "sample_app_image_tag",
		semconv127.AttributeContainerRuntime:      "cro",
		semconv127.AttributeK8SContainerName:      "kube_sample_app",
		semconv127.AttributeK8SReplicaSetName:     "sample_replica_set",
		semconv127.AttributeK8SDaemonSetName:      "sample_daemonset_name",
		semconv127.AttributeK8SPodName:            "sample_pod_name",
		semconv127.AttributeCloudProvider:         "sample_cloud_provider",
		semconv127.AttributeCloudRegion:           "sample_region",
		semconv127.AttributeCloudAvailabilityZone: "sample_zone",
		semconv127.AttributeAWSECSTaskFamily:      "sample_task_family",
		semconv127.AttributeAWSECSClusterARN:      "sample_ecs_cluster_name",
		semconv127.AttributeAWSECSContainerARN:    "sample_ecs_container_name",
		"custom_tag":                              "example_custom_tag",
		"":                                        "empty_string_key",
		"empty_string_val":                        "",
	}

	assert.Equal(t, map[string]string{
		"container_name":      "sample_app",
		"image_tag":           "sample_app_image_tag",
		"runtime":             "cro",
		"kube_container_name": "kube_sample_app",
		"kube_replica_set":    "sample_replica_set",
		"kube_daemon_set":     "sample_daemonset_name",
		"pod_name":            "sample_pod_name",
		"cloud_provider":      "sample_cloud_provider",
		"region":              "sample_region",
		"zone":                "sample_zone",
		"task_family":         "sample_task_family",
		"ecs_cluster_name":    "sample_ecs_cluster_name",
		"ecs_container_name":  "sample_ecs_container_name",
	}, ContainerTagFromAttributes(attributeMap))
}

func TestContainerTagFromAttributesEmpty(t *testing.T) {
	assert.Empty(t, ContainerTagFromAttributes(map[string]string{}))
}

func TestOriginIDFromAttributes(t *testing.T) {
	tests := []struct {
		name     string
		attrs    pcommon.Map
		originID string
	}{
		{
			name: "pod UID and container ID",
			attrs: func() pcommon.Map {
				attributes := pcommon.NewMap()
				attributes.FromRaw(map[string]interface{}{
					semconv127.AttributeContainerID: "container_id_goes_here",
					semconv127.AttributeK8SPodUID:   "k8s_pod_uid_goes_here",
				})
				return attributes
			}(),
			originID: "container_id://container_id_goes_here",
		},
		{
			name: "only container ID",
			attrs: func() pcommon.Map {
				attributes := pcommon.NewMap()
				attributes.FromRaw(map[string]interface{}{
					semconv127.AttributeContainerID: "container_id_goes_here",
				})
				return attributes
			}(),
			originID: "container_id://container_id_goes_here",
		},
		{
			name: "only pod UID",
			attrs: func() pcommon.Map {
				attributes := pcommon.NewMap()
				attributes.FromRaw(map[string]interface{}{
					semconv127.AttributeK8SPodUID: "k8s_pod_uid_goes_here",
				})
				return attributes
			}(),
			originID: "kubernetes_pod_uid://k8s_pod_uid_goes_here",
		},
		{
			name:  "none",
			attrs: pcommon.NewMap(),
		},
	}

	for _, testInstance := range tests {
		t.Run(testInstance.name, func(t *testing.T) {
			originID := OriginIDFromAttributes(testInstance.attrs)
			assert.Equal(t, testInstance.originID, originID)
		})
	}
}

func TestTagsFromAttributesIncludingDatadogNamespacedKeys(t *testing.T) {
	attributeMap := map[string]interface{}{
		// Core tags
		semconv127.AttributeServiceName:    "svc",
		semconv127.AttributeServiceVersion: "v1.2.3",
		"tags.datadoghq.com/env":           "prod",
		// Container tags
		semconv127.AttributeContainerID:        "cid",
		semconv127.AttributeContainerName:      "cname",
		semconv127.AttributeContainerImageName: "imgname",
		conventions.AttributeContainerImageTag: "imgtag",
		semconv127.AttributeContainerRuntime:   "docker",
		// Cloud tags
		semconv127.AttributeCloudProvider:         "aws",
		semconv127.AttributeCloudRegion:           "us-east-1",
		semconv127.AttributeCloudAvailabilityZone: "az1",
		// ECS tags
		semconv127.AttributeAWSECSTaskFamily:   "tfam",
		semconv127.AttributeAWSECSTaskARN:      "tarn",
		semconv127.AttributeAWSECSClusterARN:   "clarn",
		semconv127.AttributeAWSECSTaskRevision: "trev",
		semconv127.AttributeAWSECSContainerARN: "carn",
		// K8s tags
		semconv127.AttributeK8SContainerName:   "k8scont",
		semconv127.AttributeK8SClusterName:     "k8sclu",
		semconv127.AttributeK8SDeploymentName:  "k8sdep",
		semconv127.AttributeK8SReplicaSetName:  "k8srs",
		semconv127.AttributeK8SStatefulSetName: "k8ssts",
		semconv127.AttributeK8SDaemonSetName:   "k8sds",
		semconv127.AttributeK8SJobName:         "k8sj",
		semconv127.AttributeK8SCronJobName:     "k8scron",
		semconv127.AttributeK8SNamespaceName:   "k8sns",
		semconv127.AttributeK8SPodName:         "k8spod",
		// App tags
		"app.kubernetes.io/name":       "appname",
		"app.kubernetes.io/instance":   "appinst",
		"app.kubernetes.io/version":    "appver",
		"app.kuberenetes.io/component": "appcomp",
		"app.kubernetes.io/part-of":    "apppart",
		"app.kubernetes.io/managed-by": "appman",
		// All keys from keysToCheckForInDDNamespace as direct datadog.* tags
		DDNamespacePrefix + KeyEnv:               "directenv",
		DDNamespacePrefix + KeyService:           "directsvc",
		DDNamespacePrefix + KeyVersion:           "directver",
		DDNamespacePrefix + KeyContainerID:       "directcid",
		DDNamespacePrefix + KeyContainerName:     "directcname",
		DDNamespacePrefix + KeyImageName:         "directimgname",
		DDNamespacePrefix + KeyImageTag:          "directimgtag",
		DDNamespacePrefix + KeyRuntime:           "directruntime",
		DDNamespacePrefix + KeyCloudProvider:     "directcloud",
		DDNamespacePrefix + KeyRegion:            "directregion",
		DDNamespacePrefix + KeyAvailabilityZone:  "directzone",
		DDNamespacePrefix + KeyTaskFamily:        "directtfam",
		DDNamespacePrefix + KeyTaskARN:           "directtarn",
		DDNamespacePrefix + KeyTaskVersion:       "directtrev",
		DDNamespacePrefix + KeyECSClusterName:    "directclarn",
		DDNamespacePrefix + KeyECSContainerName:  "directcarn",
		DDNamespacePrefix + KeyKubeContainerName: "directk8scont",
		DDNamespacePrefix + KeyKubeClusterName:   "directk8sclu",
		DDNamespacePrefix + KeyKubeDeployment:    "directk8sdep",
		DDNamespacePrefix + KeyKubeReplicaSet:    "directk8srs",
		DDNamespacePrefix + KeyKubeStatefulSet:   "directk8ssts",
		DDNamespacePrefix + KeyKubeDaemonSet:     "directk8sds",
		DDNamespacePrefix + KeyKubeJob:           "directk8sj",
		DDNamespacePrefix + KeyKubeCronJob:       "directk8scron",
		DDNamespacePrefix + KeyKubeNamespace:     "directk8sns",
		DDNamespacePrefix + KeyPodName:           "directk8spod",
		DDNamespacePrefix + KeyKubeAppName:       "directappname",
		DDNamespacePrefix + KeyKubeAppInstance:   "directappinst",
		DDNamespacePrefix + KeyKubeAppVersion:    "directappver",
		DDNamespacePrefix + KeyKubeAppComponent:  "directappcomp",
		DDNamespacePrefix + KeyKubeAppPartOf:     "directapppart",
		DDNamespacePrefix + KeyKubeAppManagedBy:  "directappman",
	}
	attrs := pcommon.NewMap()
	attrs.FromRaw(attributeMap)

	attributeMapEmptyDD := map[string]interface{}{}
	for k, v := range attributeMap {
		if len(k) >= len(DDNamespacePrefix) && k[:len(DDNamespacePrefix)] == DDNamespacePrefix {
			attributeMapEmptyDD[k] = ""
		} else {
			attributeMapEmptyDD[k] = v
		}
	}
	attrsEmptyDD := pcommon.NewMap()
	attrsEmptyDD.FromRaw(attributeMapEmptyDD)

	t.Run("ignoreMissingDatadogFields=false (fallback enabled)", func(t *testing.T) {
		tags := GetTagsFromAttributesPreferringDatadogNamespace(attrs, false)
		// Should include only the direct datadog.* tags, not fallback
		expected := map[string]string{
			"env":                 "directenv",
			"service":             "directsvc",
			"version":             "directver",
			"container_id":        "directcid",
			"container_name":      "directcname",
			"image_name":          "directimgname",
			"image_tag":           "directimgtag",
			"runtime":             "directruntime",
			"cloud_provider":      "directcloud",
			"region":              "directregion",
			"zone":                "directzone",
			"task_family":         "directtfam",
			"task_arn":            "directtarn",
			"ecs_cluster_name":    "directclarn",
			"task_version":        "directtrev",
			"ecs_container_name":  "directcarn",
			"kube_container_name": "directk8scont",
			"kube_cluster_name":   "directk8sclu",
			"kube_deployment":     "directk8sdep",
			"kube_replica_set":    "directk8srs",
			"kube_stateful_set":   "directk8ssts",
			"kube_daemon_set":     "directk8sds",
			"kube_job":            "directk8sj",
			"kube_cronjob":        "directk8scron",
			"kube_namespace":      "directk8sns",
			"pod_name":            "directk8spod",
			"kube_app_name":       "directappname",
			"kube_app_instance":   "directappinst",
			"kube_app_version":    "directappver",
			"kube_app_component":  "directappcomp",
			"kube_app_part_of":    "directapppart",
			"kube_app_managed_by": "directappman",
		}
		for k, v := range expected {
			assert.Equal(t, v, tags[k], "expected key %s to have value %s", k, v)
		}
		// Should not include fallback tags
	})

	t.Run("ignoreMissingDatadogFields=true (fallback disabled)", func(t *testing.T) {
		tags := GetTagsFromAttributesPreferringDatadogNamespace(attrs, true)
		// Should only include direct datadog.* tags, not fallback
		expected := map[string]string{
			"env":                 "directenv",
			"service":             "directsvc",
			"version":             "directver",
			"container_id":        "directcid",
			"container_name":      "directcname",
			"image_name":          "directimgname",
			"image_tag":           "directimgtag",
			"runtime":             "directruntime",
			"cloud_provider":      "directcloud",
			"region":              "directregion",
			"zone":                "directzone",
			"task_family":         "directtfam",
			"task_arn":            "directtarn",
			"ecs_cluster_name":    "directclarn",
			"task_version":        "directtrev",
			"ecs_container_name":  "directcarn",
			"kube_container_name": "directk8scont",
			"kube_cluster_name":   "directk8sclu",
			"kube_deployment":     "directk8sdep",
			"kube_replica_set":    "directk8srs",
			"kube_stateful_set":   "directk8ssts",
			"kube_daemon_set":     "directk8sds",
			"kube_job":            "directk8sj",
			"kube_cronjob":        "directk8scron",
			"kube_namespace":      "directk8sns",
			"pod_name":            "directk8spod",
			"kube_app_name":       "directappname",
			"kube_app_instance":   "directappinst",
			"kube_app_version":    "directappver",
			"kube_app_component":  "directappcomp",
			"kube_app_part_of":    "directapppart",
			"kube_app_managed_by": "directappman",
		}
		for k, v := range expected {
			assert.Equal(t, v, tags[k], "expected key %s to have value %s", k, v)
		}
		// Should not include fallback tags
	})

	t.Run("empty datadog.* keys, ignoreMissingDatadogFields=true", func(t *testing.T) {
		tags := GetTagsFromAttributesPreferringDatadogNamespace(attrsEmptyDD, true)
		// All datadog.* keys should be present with empty values
		expected := map[string]string{
			"env":                 "",
			"service":             "",
			"version":             "",
			"container_id":        "",
			"container_name":      "",
			"image_name":          "",
			"image_tag":           "",
			"runtime":             "",
			"cloud_provider":      "",
			"region":              "",
			"zone":                "",
			"task_family":         "",
			"task_arn":            "",
			"ecs_cluster_name":    "",
			"task_version":        "",
			"ecs_container_name":  "",
			"kube_container_name": "",
			"kube_cluster_name":   "",
			"kube_deployment":     "",
			"kube_replica_set":    "",
			"kube_stateful_set":   "",
			"kube_daemon_set":     "",
			"kube_job":            "",
			"kube_cronjob":        "",
			"kube_namespace":      "",
			"pod_name":            "",
			"kube_app_name":       "",
			"kube_app_instance":   "",
			"kube_app_version":    "",
			"kube_app_component":  "",
			"kube_app_part_of":    "",
			"kube_app_managed_by": "",
		}
		for k, v := range expected {
			assert.Equal(t, v, tags[k], "expected key %s to have value %s", k, v)
		}
	})

	t.Run("empty datadog.* keys, ignoreMissingDatadogFields=false", func(t *testing.T) {
		tags := GetTagsFromAttributesPreferringDatadogNamespace(attrsEmptyDD, false)
		expected := map[string]string{
			"service":             "svc",
			"version":             "v1.2.3",
			"env":                 "prod",
			"container_id":        "cid",
			"container_name":      "cname",
			"image_name":          "imgname",
			"image_tag":           "imgtag",
			"runtime":             "docker",
			"cloud_provider":      "aws",
			"region":              "us-east-1",
			"zone":                "az1",
			"task_family":         "tfam",
			"task_arn":            "tarn",
			"ecs_cluster_name":    "clarn",
			"task_version":        "trev",
			"ecs_container_name":  "carn",
			"kube_container_name": "k8scont",
			"kube_cluster_name":   "k8sclu",
			"kube_deployment":     "k8sdep",
			"kube_replica_set":    "k8srs",
			"kube_stateful_set":   "k8ssts",
			"kube_daemon_set":     "k8sds",
			"kube_job":            "k8sj",
			"kube_cronjob":        "k8scron",
			"kube_namespace":      "k8sns",
			"pod_name":            "k8spod",
			"kube_app_name":       "appname",
			"kube_app_instance":   "appinst",
			"kube_app_version":    "appver",
			"kube_app_component":  "appcomp",
			"kube_app_part_of":    "apppart",
			"kube_app_managed_by": "appman",
		}
		assert.Equal(t, expected, tags)
	})
}
