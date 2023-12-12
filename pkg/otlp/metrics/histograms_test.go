// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestDeltaHistogramTranslatorOptions(t *testing.T) {
	tests := []struct {
		name     string
		otlpfile string
		ddogfile string
		options  []TranslatorOption
		err      string
	}{
		{
			name:     "distributions",
			otlpfile: "testdata/otlpdata/histogram/simple-delta.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta_dist-nocs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeDistributions),
			},
		},
		{
			name:     "distributions-test-min-max",
			otlpfile: "testdata/otlpdata/histogram/simple-delta-min-max.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta-min-max_dist-nocs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeDistributions),
			},
		},
		{
			name:     "distributions-count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-delta.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta_dist-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeDistributions),
				WithHistogramAggregations(),
			},
		},
		{
			name:     "buckets",
			otlpfile: "testdata/otlpdata/histogram/simple-delta.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta_counters-nocs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeCounters),
			},
		},
		{
			name:     "buckets-count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-delta.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta_counters-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeCounters),
				WithHistogramAggregations(),
			},
		},
		{
			name:     "count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-delta.json",
			ddogfile: "testdata/datadogdata/histogram/simple-delta_nobuckets-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeNoBuckets),
				WithHistogramAggregations(),
			},
		},
		{
			name: "no-count-sum-no-buckets",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeNoBuckets),
			},
			err: errNoBucketsNoSumCount,
		},
	}

	for _, testinstance := range tests {
		t.Run(testinstance.name, func(t *testing.T) {
			translator, err := NewTranslator(zap.NewNop(), testinstance.options...)
			if testinstance.err != "" {
				assert.EqualError(t, err, testinstance.err)
				return
			}
			require.NoError(t, err)
			AssertTranslatorMap(t, translator, testinstance.otlpfile, testinstance.ddogfile)
		})
	}
}

func TestCumulativeHistogramTranslatorOptions(t *testing.T) {
	tests := []struct {
		name     string
		otlpfile string
		ddogfile string
		options  []TranslatorOption
	}{
		{
			name:     "distributions",
			otlpfile: "testdata/otlpdata/histogram/simple-cumulative.json",
			ddogfile: "testdata/datadogdata/histogram/simple-cumulative_dist-nocs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeDistributions),
			},
		},
		{
			name:     "distributions-count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-cumulative.json",
			ddogfile: "testdata/datadogdata/histogram/simple-cumulative_dist-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeDistributions),
				WithHistogramAggregations(),
			},
		},
		{
			name:     "buckets",
			otlpfile: "testdata/otlpdata/histogram/simple-cumulative.json",
			ddogfile: "testdata/datadogdata/histogram/simple-cumulative_counters-nocs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeCounters),
			},
		},
		{
			name:     "buckets-count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-cumulative.json",
			ddogfile: "testdata/datadogdata/histogram/simple-cumulative_counters-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeCounters),
				WithHistogramAggregations(),
			},
		},
		{
			name:     "count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-cumulative.json",
			ddogfile: "testdata/datadogdata/histogram/simple-cumulative_nobuckets-cs.json",
			options: []TranslatorOption{
				WithHistogramMode(HistogramModeNoBuckets),
				WithHistogramAggregations(),
			},
		},
	}

	for _, testinstance := range tests {
		t.Run(testinstance.name, func(t *testing.T) {
			translator, err := NewTranslator(zap.NewNop(), testinstance.options...)
			require.NoError(t, err)
			AssertTranslatorMap(t, translator, testinstance.otlpfile, testinstance.ddogfile)
		})
	}
}

func TestExponentialHistogramTranslatorOptions(t *testing.T) {
	tests := []struct {
		name                                      string
		otlpfile                                  string
		ddogfile                                  string
		options                                   []TranslatorOption
		expectedUnknownMetricType                 int
		expectedUnsupportedAggregationTemporality int
	}{
		{
			name:                      "no-options",
			otlpfile:                  "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile:                  "testdata/datadogdata/histogram/simple-exponential.json",
			expectedUnknownMetricType: 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			// https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/26103
			name:     "empty-delta-issue-26103",
			otlpfile: "testdata/otlpdata/histogram/empty-delta-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/empty-delta-exponential.json",
		},
		{
			// https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/26103
			name:     "empty-cumulative-issue-26103",
			otlpfile: "testdata/otlpdata/histogram/empty-cumulative-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/empty-cumulative-exponential.json",
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:                      "resource-attributes-as-tags",
			otlpfile:                  "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile:                  "testdata/datadogdata/histogram/simple-exponential_res-tags.json",
			options:                   []TranslatorOption{},
			expectedUnknownMetricType: 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "count-sum",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_cs.json",
			options: []TranslatorOption{
				WithHistogramAggregations(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "instrumentation-library-metadata-as-tags",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_ilmd-tags.json",
			options: []TranslatorOption{
				WithInstrumentationLibraryMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "instrumentation-scope-metadata-as-tags",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_ismd-tags.json",
			options: []TranslatorOption{
				WithInstrumentationScopeMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "count-sum-instrumentation-library-metadata-as-tags",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_cs-ilmd-tags.json",
			options: []TranslatorOption{
				WithHistogramAggregations(),
				WithInstrumentationLibraryMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "resource-tags-instrumentation-library-metadata-as-tags",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_res-ilmd-tags.json",
			options: []TranslatorOption{
				WithInstrumentationLibraryMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "count-sum-resource-tags-instrumentation-library-metadata-as-tags",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_cs-both-tags.json",
			options: []TranslatorOption{
				WithHistogramAggregations(),
				WithInstrumentationLibraryMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
		{
			name:     "with-all",
			otlpfile: "testdata/otlpdata/histogram/simple-exponential.json",
			ddogfile: "testdata/datadogdata/histogram/simple-exponential_all.json",
			options: []TranslatorOption{
				WithHistogramAggregations(),
				WithInstrumentationLibraryMetadataAsTags(),
				WithInstrumentationScopeMetadataAsTags(),
			},
			expectedUnknownMetricType:                 1,
			expectedUnsupportedAggregationTemporality: 1,
		},
	}

	for _, testinstance := range tests {
		t.Run(testinstance.name, func(t *testing.T) {
			core, observed := observer.New(zapcore.DebugLevel)
			testLogger := zap.New(core)
			translator, err := NewTranslator(testLogger, testinstance.options...)
			require.NoError(t, err)
			AssertTranslatorMap(t, translator, testinstance.otlpfile, testinstance.ddogfile)
			assert.Equal(t, testinstance.expectedUnknownMetricType, observed.FilterMessage("Unknown or unsupported metric type").Len())
			assert.Equal(t, testinstance.expectedUnsupportedAggregationTemporality, observed.FilterMessage("Unknown or unsupported aggregation temporality").Len())
		})
	}
}
