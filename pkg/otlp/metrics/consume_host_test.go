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

package metrics

import (
	"testing"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestHostsConsumed(t *testing.T) {
	tests := []struct {
		name     string
		otlpfile string
		ddogfile string
	}{
		{
			name:     "only stats metrics",
			otlpfile: "testdata/otlpdata/hosts/stats.json",
			ddogfile: "testdata/datadogdata/hosts/stats.json",
		},
		{
			name:     "stats metrics and other metrics",
			otlpfile: "testdata/otlpdata/hosts/stats_other.json",
			ddogfile: "testdata/datadogdata/hosts/stats_other.json",
		},
		{
			name:     "only runtime metrics",
			otlpfile: "testdata/otlpdata/hosts/runtime.json",
			ddogfile: "testdata/datadogdata/hosts/runtime.json",
		},
		{
			name:     "runtime metrics and other metrics",
			otlpfile: "testdata/otlpdata/hosts/runtime_other.json",
			ddogfile: "testdata/datadogdata/hosts/runtime_other.json",
		},
		{
			name:     "only stats and runtime metrics",
			otlpfile: "testdata/otlpdata/hosts/stats_and_runtime.json",
			ddogfile: "testdata/datadogdata/hosts/stats_and_runtime.json",
		},
	}

	for _, testinstance := range tests {
		t.Run(testinstance.name, func(t *testing.T) {
			set := componenttest.NewNopTelemetrySettings()
			attributesTranslator, err := attributes.NewTranslator(set)
			require.NoError(t, err)
			translator, err := NewTranslator(set, attributesTranslator,
				WithOriginProduct(OriginProductDatadogAgent),
			)
			require.NoError(t, err)
			AssertTranslatorMap(t, translator, testinstance.otlpfile, testinstance.ddogfile)
		})
	}
}
