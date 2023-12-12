// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.

package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestSummaryMetrics(t *testing.T) {
	tests := []struct {
		name     string
		otlpfile string
		ddogfile string
		options  []TranslatorOption
		tags     []string
	}{
		{
			name:     "summary",
			otlpfile: "testdata/otlpdata/summary/simple.json",
			ddogfile: "testdata/datadogdata/summary/simple_summary.json",
			options:  []TranslatorOption{WithFallbackSourceProvider(testProvider("fallbackHostname"))},
		},
		{
			name:     "summary-with-quantiles",
			otlpfile: "testdata/otlpdata/summary/simple.json",
			ddogfile: "testdata/datadogdata/summary/simple_summary-with-quantile.json",
			options: []TranslatorOption{
				WithFallbackSourceProvider(testProvider("fallbackHostname")),
				WithQuantiles(),
			},
		},
		{
			name:     "summary-with-attributes",
			otlpfile: "testdata/otlpdata/summary/with-attributes.json",
			ddogfile: "testdata/datadogdata/summary/with-attributes_summary.json",
			options:  []TranslatorOption{WithFallbackSourceProvider(testProvider("fallbackHostname"))},
			tags:     []string{"attribute_tag:attribute_value"},
		},
		{
			name:     "summary-with-attributes-quantiles",
			otlpfile: "testdata/otlpdata/summary/with-attributes.json",
			ddogfile: "testdata/datadogdata/summary/with-attributes-quantile_summary.json",
			options: []TranslatorOption{
				WithFallbackSourceProvider(testProvider("fallbackHostname")),
				WithQuantiles(),
			},
			tags: []string{"attribute_tag:attribute_value"},
		},
	}

	for _, testinstance := range tests {
		t.Run(testinstance.name, func(t *testing.T) {
			translator, err := NewTranslator(componenttest.NewNopTelemetrySettings(), testinstance.options...)
			require.NoError(t, err)
			AssertTranslatorMap(t, translator, testinstance.otlpfile, testinstance.ddogfile)
		})
	}
}
