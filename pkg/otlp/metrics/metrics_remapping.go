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
	"strings"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

const (
	// divMebibytes specifies the number of bytes in a mebibyte.
	divMebibytes = 1024 * 1024
	// divPercentage specifies the division necessary for converting fractions to percentages.
	divPercentage = 0.01
)

// remapMetrics extracts any Datadog specific metrics from m and appends them to all.
func remapMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	remapSystemMetrics(all, m)
	remapContainerMetrics(all, m)
	remapKafkaMetrics(all, m)
	remapJvmMetrics(all, m)
}

// remapSystemMetrics extracts system metrics from m and appends them to all.
func remapSystemMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	name := m.Name()
	if !strings.HasPrefix(name, "process.") && !strings.HasPrefix(name, "system.") {
		// not a system metric
		return
	}
	switch name {
	case "system.cpu.load_average.1m":
		copyMetricWithAttr(all, m, "system.load.1", 1, attributesMapping{})
	case "system.cpu.load_average.5m":
		copyMetricWithAttr(all, m, "system.load.5", 1, attributesMapping{})
	case "system.cpu.load_average.15m":
		copyMetricWithAttr(all, m, "system.load.15", 1, attributesMapping{})
	case "system.cpu.utilization":
		copyMetricWithAttr(all, m, "system.cpu.idle", divPercentage, attributesMapping{}, kv{"state", "idle"})
		copyMetricWithAttr(all, m, "system.cpu.user", divPercentage, attributesMapping{}, kv{"state", "user"})
		copyMetricWithAttr(all, m, "system.cpu.system", divPercentage, attributesMapping{}, kv{"state", "system"})
		copyMetricWithAttr(all, m, "system.cpu.iowait", divPercentage, attributesMapping{}, kv{"state", "wait"})
		copyMetricWithAttr(all, m, "system.cpu.stolen", divPercentage, attributesMapping{}, kv{"state", "steal"})
	case "system.memory.usage":
		copyMetricWithAttr(all, m, "system.mem.total", divMebibytes, attributesMapping{})
		copyMetricWithAttr(all, m, "system.mem.usable", divMebibytes, attributesMapping{},
			kv{"state", "free"},
			kv{"state", "cached"},
			kv{"state", "buffered"},
		)
	case "system.network.io":
		copyMetricWithAttr(all, m, "system.net.bytes_rcvd", 1, attributesMapping{}, kv{"direction", "receive"})
		copyMetricWithAttr(all, m, "system.net.bytes_sent", 1, attributesMapping{}, kv{"direction", "transmit"})
	case "system.paging.usage":
		copyMetricWithAttr(all, m, "system.swap.free", divMebibytes, attributesMapping{}, kv{"state", "free"})
		copyMetricWithAttr(all, m, "system.swap.used", divMebibytes, attributesMapping{}, kv{"state", "used"})
	case "system.filesystem.utilization":
		copyMetricWithAttr(all, m, "system.disk.in_use", 1, attributesMapping{})
	}
	// process.* and system.* metrics need to be prepended with the otel.* namespace
	m.SetName("otel." + m.Name())
}

// remapContainerMetrics extracts system metrics from m and appends them to all.
func remapContainerMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	name := m.Name()
	if !strings.HasPrefix(name, "container.") {
		// not a container metric
		return
	}
	switch name {
	case "container.cpu.usage.total":
		if addm, ok := copyMetricWithAttr(all, m, "container.cpu.usage", 1, attributesMapping{}); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.usage.usermode":
		if addm, ok := copyMetricWithAttr(all, m, "container.cpu.user", 1, attributesMapping{}); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.usage.system":
		if addm, ok := copyMetricWithAttr(all, m, "container.cpu.system", 1, attributesMapping{}); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.throttling_data.throttled_time":
		copyMetricWithAttr(all, m, "container.cpu.throttled", 1, attributesMapping{})
	case "container.cpu.throttling_data.throttled_periods":
		copyMetricWithAttr(all, m, "container.cpu.throttled.periods", 1, attributesMapping{})
	case "container.memory.usage.total":
		copyMetricWithAttr(all, m, "container.memory.usage", 1, attributesMapping{})
	case "container.memory.active_anon":
		copyMetricWithAttr(all, m, "container.memory.kernel", 1, attributesMapping{})
	case "container.memory.hierarchical_memory_limit":
		copyMetricWithAttr(all, m, "container.memory.limit", 1, attributesMapping{})
	case "container.memory.usage.limit":
		copyMetricWithAttr(all, m, "container.memory.soft_limit", 1, attributesMapping{})
	case "container.memory.total_cache":
		copyMetricWithAttr(all, m, "container.memory.cache", 1, attributesMapping{})
	case "container.memory.total_swap":
		copyMetricWithAttr(all, m, "container.memory.swap", 1, attributesMapping{})
	case "container.blockio.io_service_bytes_recursive":
		copyMetricWithAttr(all, m, "container.io.write", 1, attributesMapping{}, kv{"operation", "write"})
		copyMetricWithAttr(all, m, "container.io.read", 1, attributesMapping{}, kv{"operation", "read"})
	case "container.blockio.io_serviced_recursive":
		copyMetricWithAttr(all, m, "container.io.write.operations", 1, attributesMapping{}, kv{"operation", "write"})
		copyMetricWithAttr(all, m, "container.io.read.operations", 1, attributesMapping{}, kv{"operation", "read"})
	case "container.network.io.usage.tx_bytes":
		copyMetricWithAttr(all, m, "container.net.sent", 1, attributesMapping{})
	case "container.network.io.usage.tx_packets":
		copyMetricWithAttr(all, m, "container.net.sent.packets", 1, attributesMapping{})
	case "container.network.io.usage.rx_bytes":
		copyMetricWithAttr(all, m, "container.net.rcvd", 1, attributesMapping{})
	case "container.network.io.usage.rx_packets":
		copyMetricWithAttr(all, m, "container.net.rcvd.packets", 1, attributesMapping{})
	}
}

type (
	// kv represents a key/value pair.
	kv struct{ K, V string }

	// attributesMapping contains to mapping of attributes from OTel to DD.
	attributesMapping struct {
		// fixed represents attributes that need to be mapped where the value is
		// already known based on the OTel metric name.
		fixed map[string]string
		// dynamic represents attributes that need to be mapped where the value needs
		// to be dynamically pulled from a data point attribute. Typically when the OTel
		// metric and DD metric have different conventions (e.g. group vs consumer_group).
		dynamic map[string]string
	}
)

// copyMetric copies metric m to dest. The new metric's name will be newname, and all of its datapoints will
// be divided by div. If filter is provided, only the data points that have *either* of the specified string
// attributes will be copied over. If the filtering results in no datapoints, no new metric is added to dest.
// It will add any attributes specified in attributesMapping, by either pulling the value from the datapoint
// for dynamic attributes, or setting the given attribute for fixed attributes.
//
// copyMetric returns the new metric and reports whether it was added to dest.
//
// Please note that copyMetric is restricted to the metric types Sum and Gauge.
func copyMetricWithAttr(dest pmetric.MetricSlice, m pmetric.Metric, newname string, div float64, attributesMapping attributesMapping, filter ...kv) (pmetric.Metric, bool) {
	newm := pmetric.NewMetric()
	m.CopyTo(newm)
	newm.SetName(newname)
	var dps pmetric.NumberDataPointSlice
	switch newm.Type() {
	case pmetric.MetricTypeGauge:
		dps = newm.Gauge().DataPoints()
	case pmetric.MetricTypeSum:
		dps = newm.Sum().DataPoints()
	default:
		// invalid metric type
		return newm, false
	}
	dps.RemoveIf(func(dp pmetric.NumberDataPoint) bool {
		if !hasAny(dp, filter...) {
			return true
		}
		switch dp.ValueType() {
		case pmetric.NumberDataPointValueTypeInt:
			if div >= 1 {
				// avoid division by zero
				dp.SetIntValue(dp.IntValue() / int64(div))
			}
		case pmetric.NumberDataPointValueTypeDouble:
			if div != 0 {
				dp.SetDoubleValue(dp.DoubleValue() / div)
			}
		}
		// attributes mapping
		if fixed := attributesMapping.fixed; fixed != nil {
			for k, v := range fixed {
				dp.Attributes().PutStr(k, v)
			}
		}
		if dynamic := attributesMapping.dynamic; dynamic != nil {
			for old, new := range dynamic {
				if v, ok := dp.Attributes().Get(old); ok {
					dp.Attributes().PutStr(new, v.AsString())
				}
			}
		}
		return false
	})
	if dps.Len() > 0 {
		// if we have datapoints, copy it
		addm := dest.AppendEmpty()
		newm.CopyTo(addm)
		return addm, true
	}
	return newm, false
}

// hasAny reports whether point has any of the given string tags.
// If no tags are provided it returns true.
func hasAny(point pmetric.NumberDataPoint, tags ...kv) bool {
	if len(tags) == 0 {
		return true
	}
	attr := point.Attributes()
	for _, tag := range tags {
		v, ok := attr.Get(tag.K)
		if !ok {
			continue
		}
		if v.Str() == tag.V {
			return true
		}
	}
	return false
}
