// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metricscommon"
)

// getTags maps an attributeMap into a slice of Datadog tags
func getTags(labels pcommon.Map) []string {
	tags := make([]string, 0, labels.Len())
	labels.Range(func(key string, value pcommon.Value) bool {
		v := value.AsString()
		tags = append(tags, metricscommon.FormatKeyValueTag(key, v))
		return true
	})
	return tags
}

// WithAttributeMap creates a new metricDimensions struct with additional tags from attributes.
func WithAttributeMap(d *metricscommon.Dimensions, labels pcommon.Map) *metricscommon.Dimensions {
	return d.AddTags(getTags(labels)...)
}
