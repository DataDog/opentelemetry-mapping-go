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

	"github.com/stretchr/testify/assert"
)

func TestOriginServiceFromScopeName(t *testing.T) {
	tests := []struct {
		scopeName string
		expected  OriginService
	}{
		{
			scopeName: "otelcol/notsupportedreceiver",
			expected:  OriginServiceUnknown,
		},
		{
			scopeName: "otelcol/kubeletstatsreceiver",
			expected:  OriginServiceKubeletStatsReceiver,
		},
		{
			scopeName: "otelcol/hostmetricsreceiver/memory",
			expected:  OriginServiceHostMetricsReceiver,
		},
		{
			scopeName: "go.opentelemetry.io/otel/metric/example",
			expected:  OriginServiceUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.scopeName, func(t *testing.T) {
			service := originServiceFromScopeName(tt.scopeName)
			assert.Equal(t, tt.expected, service)
		})
	}
}

func TestOriginFull(t *testing.T) {
	translator := NewTestTranslator(t, WithOriginProduct(OriginProduct(42)))
	AssertTranslatorMap(t, translator,
		"testdata/otlpdata/origin/origin.json",
		"testdata/datadogdata/origin/origin.json",
	)
}
