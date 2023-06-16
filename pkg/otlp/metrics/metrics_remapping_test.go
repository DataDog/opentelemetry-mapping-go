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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestRemapMetrics(t *testing.T) {
	// point is a datapoint
	type point struct {
		// i defines a IntValue datapoint when non-zero
		i int64
		// f defines a DoubleValue datapoint when non-zero
		f float64
		// attrs specifies the raw attributes of the datapoint
		attrs map[string]any
	}
	// metric is a convenience function to create a new metric with the given name
	// and set of datapoints
	metric := func(name string, dps ...point) pmetric.Metric {
		out := pmetric.NewMetric()
		out.SetName(name)
		g := out.SetEmptyGauge()
		for _, d := range dps {
			p := g.DataPoints().AppendEmpty()
			if d.i != 0 {
				p.SetIntValue(d.i)
			} else {
				p.SetDoubleValue(d.f)
			}
			p.Attributes().FromRaw(d.attrs)
		}
		return out
	}
	chunit := func(m pmetric.Metric, typ string) pmetric.Metric {
		m.SetUnit(typ)
		return m
	}

	dest := pmetric.NewMetricSlice()
	for _, tt := range []struct {
		in  pmetric.Metric
		out []pmetric.Metric
	}{
		{
			in:  metric("system.cpu.load_average.1m", point{f: 1}),
			out: []pmetric.Metric{metric("system.load.1", point{f: 1})},
		},
		{
			in:  metric("system.cpu.load_average.5m", point{f: 5}),
			out: []pmetric.Metric{metric("system.load.5", point{f: 5})},
		},
		{
			in:  metric("system.cpu.load_average.15m", point{f: 15}),
			out: []pmetric.Metric{metric("system.load.15", point{f: 15})},
		},
		{
			in: metric("system.cpu.utilization",
				point{f: 1, attrs: map[string]any{"state": "idle"}},
				point{f: 2, attrs: map[string]any{"state": "user"}},
				point{f: 3, attrs: map[string]any{"state": "system"}},
				point{f: 5, attrs: map[string]any{"state": "wait"}},
				point{f: 8, attrs: map[string]any{"state": "steal"}},
				point{f: 13, attrs: map[string]any{"state": "other"}},
			),
			out: []pmetric.Metric{
				metric("system.cpu.idle",
					point{f: 100, attrs: map[string]any{"state": "idle"}}),
				metric("system.cpu.user",
					point{f: 200, attrs: map[string]any{"state": "user"}}),
				metric("system.cpu.system",
					point{f: 300, attrs: map[string]any{"state": "system"}}),
				metric("system.cpu.iowait",
					point{f: 500, attrs: map[string]any{"state": "wait"}}),
				metric("system.cpu.stolen",
					point{f: 800, attrs: map[string]any{"state": "steal"}}),
			},
		},
		{
			in:  metric("system.cpu.utilization", point{i: 5, attrs: map[string]any{"state": "idle"}}),
			out: []pmetric.Metric{metric("system.cpu.idle", point{i: 5, attrs: map[string]any{"state": "idle"}})},
		},
		{
			in: metric("system.memory.usage",
				point{f: divMebibytes * 1, attrs: map[string]any{"state": "free"}},
				point{f: divMebibytes * 2, attrs: map[string]any{"state": "cached"}},
				point{f: divMebibytes * 3, attrs: map[string]any{"state": "buffered"}},
				point{f: divMebibytes * 5, attrs: map[string]any{"state": "steal"}},
				point{f: divMebibytes * 8, attrs: map[string]any{"state": "system"}},
				point{f: divMebibytes * 13, attrs: map[string]any{"state": "user"}},
			),
			out: []pmetric.Metric{
				metric("system.mem.total",
					point{f: 1, attrs: map[string]any{"state": "free"}},
					point{f: 2, attrs: map[string]any{"state": "cached"}},
					point{f: 3, attrs: map[string]any{"state": "buffered"}},
					point{f: 5, attrs: map[string]any{"state": "steal"}},
					point{f: 8, attrs: map[string]any{"state": "system"}},
					point{f: 13, attrs: map[string]any{"state": "user"}},
				),
				metric("system.mem.usable",
					point{f: 1, attrs: map[string]any{"state": "free"}},
					point{f: 2, attrs: map[string]any{"state": "cached"}},
					point{f: 3, attrs: map[string]any{"state": "buffered"}},
				),
			},
		},
		{
			in:  metric("system.memory.usage", point{i: divMebibytes * 5}),
			out: []pmetric.Metric{metric("system.mem.total", point{i: 5})},
		},
		{
			in: metric("system.network.io",
				point{f: 1, attrs: map[string]any{"direction": "receive"}},
				point{f: 2, attrs: map[string]any{"direction": "transmit"}},
				point{f: 3, attrs: map[string]any{"state": "buffered"}},
			),
			out: []pmetric.Metric{
				metric("system.net.bytes_rcvd",
					point{f: 1, attrs: map[string]any{"direction": "receive"}},
				),
				metric("system.net.bytes_sent",
					point{f: 2, attrs: map[string]any{"direction": "transmit"}},
				),
			},
		},
		{
			in: metric("system.paging.usage",
				point{f: divMebibytes * 1, attrs: map[string]any{"state": "free"}},
				point{f: divMebibytes * 2, attrs: map[string]any{"state": "used"}},
				point{f: 3, attrs: map[string]any{"state": "buffered"}},
			),
			out: []pmetric.Metric{
				metric("system.swap.free",
					point{f: 1, attrs: map[string]any{"state": "free"}},
				),
				metric("system.swap.used",
					point{f: 2, attrs: map[string]any{"state": "used"}},
				),
			},
		},
		{
			in:  metric("system.filesystem.utilization", point{f: 15}),
			out: []pmetric.Metric{metric("system.disk.in_use", point{f: 15})},
		},
		{
			in:  metric("other.metric", point{f: 15}),
			out: []pmetric.Metric{},
		},
		{
			in: metric("container.cpu.usage.total", point{f: 15}),
			out: []pmetric.Metric{
				chunit(metric("container.cpu.usage", point{f: 15}), "nanocore"),
			},
		},
		{
			in: metric("container.cpu.usage.usermode", point{f: 15}),
			out: []pmetric.Metric{
				chunit(metric("container.cpu.user", point{f: 15}), "nanocore"),
			},
		},
		{
			in: metric("container.cpu.usage.system", point{f: 15}),
			out: []pmetric.Metric{
				chunit(metric("container.cpu.system", point{f: 15}), "nanocore"),
			},
		},
		{
			in:  metric("container.cpu.throttling_data.throttled_time", point{f: 15}),
			out: []pmetric.Metric{metric("container.cpu.throttled", point{f: 15})},
		},
		{
			in:  metric("container.cpu.throttling_data.throttled_periods", point{f: 15}),
			out: []pmetric.Metric{metric("container.cpu.throttled.periods", point{f: 15})},
		},
		{
			in:  metric("container.memory.usage.total", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.usage", point{f: 15})},
		},
		{
			in:  metric("container.memory.active_anon", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.kernel", point{f: 15})},
		},
		{
			in:  metric("container.memory.hierarchical_memory_limit", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.limit", point{f: 15})},
		},
		{
			in:  metric("container.memory.usage.limit", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.soft_limit", point{f: 15})},
		},
		{
			in:  metric("container.memory.total_cache", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.cache", point{f: 15})},
		},
		{
			in:  metric("container.memory.total_swap", point{f: 15}),
			out: []pmetric.Metric{metric("container.memory.swap", point{f: 15})},
		},
		{
			in: metric("container.blockio.io_service_bytes_recursive",
				point{f: 1, attrs: map[string]any{"operation": "write"}},
				point{f: 2, attrs: map[string]any{"operation": "read"}},
				point{f: 3, attrs: map[string]any{"state": "buffered"}},
			),
			out: []pmetric.Metric{
				metric("container.io.write",
					point{f: 1, attrs: map[string]any{"operation": "write"}}),
				metric("container.io.read",
					point{f: 2, attrs: map[string]any{"operation": "read"}}),
			},
		},
		{
			in: metric("container.blockio.io_service_bytes_recursive",
				point{f: 1, attrs: map[string]any{"operation": "write"}},
			),
			out: []pmetric.Metric{
				metric("container.io.write",
					point{f: 1, attrs: map[string]any{"operation": "write"}}),
			},
		},
		{
			in: metric("container.blockio.io_serviced_recursive",
				point{f: 1, attrs: map[string]any{"operation": "write"}},
				point{f: 2, attrs: map[string]any{"operation": "read"}},
				point{f: 3, attrs: map[string]any{"state": "buffered"}},
			),
			out: []pmetric.Metric{
				metric("container.io.write.operations",
					point{f: 1, attrs: map[string]any{"operation": "write"}}),
				metric("container.io.read.operations",
					point{f: 2, attrs: map[string]any{"operation": "read"}}),
			},
		},
		{
			in: metric("container.blockio.io_serviced_recursive",
				point{f: 1, attrs: map[string]any{"xoperation": "write"}},
				point{f: 2, attrs: map[string]any{"xoperation": "read"}},
				point{f: 3, attrs: map[string]any{"state": "buffered"}},
			),
			out: nil, // no datapoints match filter
		},
		{
			in:  metric("container.network.io.usage.tx_bytes", point{f: 15}),
			out: []pmetric.Metric{metric("container.net.sent", point{f: 15})},
		},
		{
			in:  metric("container.network.io.usage.tx_packets", point{f: 15}),
			out: []pmetric.Metric{metric("container.net.sent.packets", point{f: 15})},
		},
		{
			in:  metric("container.network.io.usage.rx_bytes", point{f: 15}),
			out: []pmetric.Metric{metric("container.net.rcvd", point{f: 15})},
		},
		{
			in:  metric("container.network.io.usage.rx_packets", point{f: 15}),
			out: []pmetric.Metric{metric("container.net.rcvd.packets", point{f: 15})},
		},
	} {
		lena := dest.Len()
		checkprefix := strings.HasPrefix(tt.in.Name(), "system.") || strings.HasPrefix(tt.in.Name(), "process.")
		remapMetrics(dest, tt.in)
		if checkprefix {
			require.True(t, strings.HasPrefix(tt.in.Name(), "otel."), "system.* and process.* metrics need to be prepended with the otel.* namespace")
		}
		require.Equal(t, dest.Len()-lena, len(tt.out), "unexpected number of metrics added")
		for i, out := range tt.out {
			require.Equal(t, out, dest.At(dest.Len()-len(tt.out)+i))
		}
	}

}

func TestCopyMetric(t *testing.T) {
	m := pmetric.NewMetric()
	m.SetName("test.metric")
	m.SetDescription("metric-description")
	m.SetUnit("cm")

	dest := pmetric.NewMetricSlice()
	t.Run("gauge", func(t *testing.T) {
		v := m.SetEmptyGauge()
		dp := v.DataPoints().AppendEmpty()
		dp.SetDoubleValue(12)
		dp.Attributes().FromRaw(map[string]any{"fruit": "apple", "count": 15})
		dp = v.DataPoints().AppendEmpty()
		dp.SetIntValue(24)
		dp.Attributes().FromRaw(map[string]any{"human": "Ann", "age": 25})

		t.Run("plain", func(t *testing.T) {
			out, ok := copyMetric(dest, m, "copied.test.metric", 1)
			require.True(t, ok)
			require.Equal(t, m.Name(), "test.metric")
			require.Equal(t, out.Name(), "copied.test.metric")
			sameExceptName(t, m, out)
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("div", func(t *testing.T) {
			out, ok := copyMetric(dest, m, "copied.test.metric", 2)
			require.True(t, ok)
			require.Equal(t, out.Name(), "copied.test.metric")
			require.Equal(t, out.Gauge().DataPoints().At(0).DoubleValue(), 6.)
			require.Equal(t, out.Gauge().DataPoints().At(1).IntValue(), int64(12))
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("filter", func(t *testing.T) {
			out, ok := copyMetric(dest, m, "copied.test.metric", 1, kv{"human", "Ann"})
			require.True(t, ok)
			require.Equal(t, out.Name(), "copied.test.metric")
			require.Equal(t, out.Gauge().DataPoints().Len(), 1)
			require.Equal(t, out.Gauge().DataPoints().At(0).IntValue(), int64(24))
			require.Equal(t, out.Gauge().DataPoints().At(0).Attributes().AsRaw(), map[string]any{"human": "Ann", "age": int64(25)})
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("none", func(t *testing.T) {
			_, ok := copyMetric(dest, m, "copied.test.metric", 1, kv{"human", "Paul"})
			require.False(t, ok)
		})
	})

	t.Run("sum", func(t *testing.T) {
		dp := m.SetEmptySum().DataPoints().AppendEmpty()
		dp.SetDoubleValue(12)
		dp.Attributes().FromRaw(map[string]any{"fruit": "apple", "count": 15})
		out, ok := copyMetric(dest, m, "copied.test.metric", 1)
		require.True(t, ok)
		require.Equal(t, out.Name(), "copied.test.metric")
		sameExceptName(t, m, out)
		require.Equal(t, dest.At(dest.Len()-1), out)
	})

	t.Run("histogram", func(t *testing.T) {
		dp := m.SetEmptyHistogram().DataPoints().AppendEmpty()
		dp.SetCount(12)
		dp.SetMax(44)
		dp.SetMin(3)
		dp.SetSum(120)
		_, ok := copyMetric(dest, m, "copied.test.metric", 1)
		require.False(t, ok)
	})
}

func TestHasAny(t *testing.T) {
	// p returns a numberic data point having the given attributes.
	p := func(m map[string]any) pmetric.NumberDataPoint {
		v := pmetric.NewNumberDataPoint()
		if err := v.Attributes().FromRaw(m); err != nil {
			t.Fatalf("Error generating data point: %v", err)
		}
		return v
	}
	for i, tt := range []struct {
		attrs map[string]any
		tags  []kv
		out   bool
	}{
		{
			attrs: map[string]any{
				"fruit": "apple",
				"human": "Ann",
			},
			tags: []kv{{"human", "Ann"}},
			out:  true,
		},
		{
			attrs: map[string]any{
				"fruit": "apple",
				"human": "Ann",
			},
			tags: []kv{{"human", "ann"}},
			out:  false,
		},
		{
			attrs: map[string]any{
				"fruit":   "apple",
				"human":   "Ann",
				"company": "Paul",
			},
			tags: []kv{{"human", "ann"}, {"company", "Paul"}},
			out:  true,
		},
		{
			attrs: map[string]any{
				"fruit":   "apple",
				"human":   "Ann",
				"company": "Paul",
			},
			tags: []kv{{"fruit", "apple"}, {"company", "Paul"}},
			out:  true,
		},
		{
			attrs: map[string]any{
				"fruit":   "apple",
				"human":   "Ann",
				"company": "Paul",
			},
			tags: nil,
			out:  true,
		},
		{
			attrs: map[string]any{
				"fruit":   "apple",
				"human":   "Ann",
				"company": "Paul",
				"number":  4,
			},
			tags: []kv{{"number", "4"}},
			out:  false,
		},
		{
			attrs: nil,
			tags:  []kv{{"number", "4"}},
			out:   false,
		},
		{
			attrs: nil,
			tags:  nil,
			out:   true,
		},
	} {
		require.Equal(t, hasAny(p(tt.attrs), tt.tags...), tt.out, fmt.Sprint(i))
	}
}

// sameExceptName validates that metrics a and b are the same by disregarding
// their names.
func sameExceptName(t *testing.T, a, b pmetric.Metric) {
	aa, bb := pmetric.NewMetric(), pmetric.NewMetric()
	a.CopyTo(aa)
	b.CopyTo(bb)
	aa.SetName("🙂")
	bb.SetName("🙂")
	require.Equal(t, aa, bb)
}
