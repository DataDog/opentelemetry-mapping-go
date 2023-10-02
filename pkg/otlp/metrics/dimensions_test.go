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
	"testing"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/metricscommon"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestWithAttributeMap(t *testing.T) {
	attributes := pcommon.NewMap()
	attributes.FromRaw(map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
		"key3": "",
	})

	dims := metricscommon.Dimensions{}
	assert.ElementsMatch(t,
		WithAttributeMap(&dims, attributes).Tags(),
		[...]string{"key1:val1", "key2:val2", "key3:n/a"},
	)
}

func TestAllFieldsAreCopied(t *testing.T) {
	dims := metricscommon.NewDimensions(
		"example.name",
		[]string{"tagOne:a", "tagTwo:b"},
		"hostname",
		"origin_id",
	)

	attributes := pcommon.NewMap()
	attributes.FromRaw(map[string]interface{}{
		"tagFour": "d",
	})
	newDims := dims.AddTags("tagThree:c").WithSuffix("suffix")
	newDims = WithAttributeMap(newDims, attributes)

	assert.Equal(t, "example.name.suffix", newDims.Name())
	assert.Equal(t, "hostname", newDims.Host())
	assert.ElementsMatch(t, []string{"tagOne:a", "tagTwo:b", "tagThree:c", "tagFour:d"}, newDims.Tags())
	assert.Equal(t, "origin_id", newDims.OriginID())
}
