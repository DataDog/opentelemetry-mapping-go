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

var (
	metricAttrs = map[string]string{
		"kafka.producer.response-rate":               "producer-metrics",
		"kafka.producer.request-latency-avg":         "producer-metrics",
		"kafka.producer.outgoing-byte-rate":          "producer-metrics",
		"kafka.producer.io-wait-time-ns-avg":         "producer-metrics",
		"kafka.producer.byte-rate":                   "producer-topic-metrics",
		"kafka.consumer.total.bytes-consumed-rate":   "consumer-fetch-manager-metrics",
		"kafka.consumer.total.records-consumed-rate": "consumer-fetch-manager-metrics",
		"kafka.network.io":                           "BrokerTopicMetrics",
		"kafka.purgatory.size":                       "DelayedOperationPurgatory",
		"kafka.partition.under_replicated":           "ReplicaManager",
		"kafka.isr.operation.count":                  "ReplicaManager",
		"kafka.leader.election.rate":                 "ControllerStats",
		"kafka.partition.offline":                    "KafkaController",
		"kafka.request.time.avg":                     "RequestMetrics",
		"jvm.gc.collections.count":                   "GarbageCollector",
		"jvm.gc.collections.elapsed":                 "GarbageCollector",
	}
)

// remapMetrics extracts any Datadog specific metrics from m and appends them to all.
func remapMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	remapSystemMetrics(all, m)
	remapContainerMetrics(all, m)
	remapKafkaMetrics(all, m)
	remapJvmMetrics(all, m)
}

func remapKafkaMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	name := m.Name()
	if !strings.HasPrefix(name, "kafka.") {
		// not a kafka metric
		return
	}
	switch name {
	// no change necessary, as - gets converted to _.
	case "kafka.producer.request-rate":
	case "kafka.producer.response-rate":
	case "kafka.producer.request-latency-avg":

	case "kafka.producer.outgoing-byte-rate":
		copyMetric(all, m, "kafka.producer.bytes_out", 1)
	case "kafka.producer.io-wait-time-ns-avg":
		copyMetric(all, m, "kafka.producer.io_wait", 1)
	case "kafka.producer.byte-rate":
		// may need to add tags client, topic.
		copyMetric(all, m, "kafka.producer.bytes_out", 1)
	case "kafka.consumer.total.bytes-consumed-rate":
		copyMetric(all, m, "kafka.consumer.bytes_in", 1)
	case "kafka.consumer.total.records-consumed-rate":
		copyMetric(all, m, "kafka.consumer.messages_in", 1)
	case "kafka.network.io":
		copyMetric(all, m, "kafka.net.bytes_out.rate", 1, kv{"state", "out"})
		copyMetric(all, m, "kafka.net.bytes_in.rate", 1, kv{"state", "in"})
	case "kafka.purgatory.size":
		copyMetric(all, m, "kafka.request.producer_request_purgatory.size", 1, kv{"type", "produce"})
		copyMetric(all, m, "kafka.request.fetch_request_purgatory.size", 1, kv{"type", "fetch"})
	case "kafka.partition.under_replicated":
		copyMetric(all, m, "kafka.replication.under_replicated_partitions", 1)
	case "kafka.isr.operation.count":
		copyMetric(all, m, "kafka.replication.isr_shrinks.rate", 1, kv{"operation", "shrink"})
		copyMetric(all, m, "kafka.replication.isr_expands.rate", 1, kv{"operation", "extract"})
	case "kafka.leader.election.rate":
		copyMetric(all, m, "kafka.replication.leader_elections.rate", 1)
	case "kafka.partition.offline":
		copyMetric(all, m, "kafka.replication.offline_partitions_count", 1)
	case "kafka.request.time.avg":
		copyMetric(all, m, "kafka.request.produce.time.avg", 1, kv{"type", "produce"})
		copyMetric(all, m, "kafka.request.fetch_consumer.time.avg", 1, kv{"type", "fetchconsumer"})
		copyMetric(all, m, "kafka.request.fetch_follower.time.avg", 1, kv{"type", "fetchfollower"})
	}
}

func remapJvmMetrics(all pmetric.MetricSlice, m pmetric.Metric) {
	name := m.Name()
	if !strings.HasPrefix(name, "jvm.") {
		// not a jvm metric
		return
	}
	switch name {
	case "jvm.gc.collections.count":
		// Young Gen Collectors
		copyMetric(all, m, "jvm.gc.minor_collection_count", 1, kv{"name", "Copy"})
		copyMetric(all, m, "jvm.gc.minor_collection_count", 1, kv{"name", "PS Scavenge"})
		copyMetric(all, m, "jvm.gc.minor_collection_count", 1, kv{"name", "ParNew"})
		copyMetric(all, m, "jvm.gc.minor_collection_count", 1, kv{"name", "G1 Young Generation"})
		// Old Gen Collectors
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "MarkSweepCompact"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "PS MarkSweep"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "ConcurrentMarkSweep"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "G1 Mixed Generation"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "G1 Old Generation"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "Shenandoah Cycles"})
		copyMetric(all, m, "jvm.gc.major_collection_count", 1, kv{"name", "ZGC"})

	case "jvm.gc.collections.elapsed":
		// Young Gen Collectors
		copyMetric(all, m, "jvm.gc.minor_collection_time", 1, kv{"name", "Copy"})
		copyMetric(all, m, "jvm.gc.minor_collection_time", 1, kv{"name", "PS Scavenge"})
		copyMetric(all, m, "jvm.gc.minor_collection_time", 1, kv{"name", "ParNew"})
		copyMetric(all, m, "jvm.gc.minor_collection_time", 1, kv{"name", "G1 Young Generation"})
		// Old Gen Collectors
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "MarkSweepCompact"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "PS MarkSweep"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "ConcurrentMarkSweep"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "G1 Mixed Generation"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "G1 Old Generation"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "Shenandoah Cycles"})
		copyMetric(all, m, "jvm.gc.major_collection_time", 1, kv{"name", "ZGC"})
	}
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
		copyMetric(all, m, "system.load.1", 1)
	case "system.cpu.load_average.5m":
		copyMetric(all, m, "system.load.5", 1)
	case "system.cpu.load_average.15m":
		copyMetric(all, m, "system.load.15", 1)
	case "system.cpu.utilization":
		copyMetric(all, m, "system.cpu.idle", divPercentage, kv{"state", "idle"})
		copyMetric(all, m, "system.cpu.user", divPercentage, kv{"state", "user"})
		copyMetric(all, m, "system.cpu.system", divPercentage, kv{"state", "system"})
		copyMetric(all, m, "system.cpu.iowait", divPercentage, kv{"state", "wait"})
		copyMetric(all, m, "system.cpu.stolen", divPercentage, kv{"state", "steal"})
	case "system.memory.usage":
		copyMetric(all, m, "system.mem.total", divMebibytes)
		copyMetric(all, m, "system.mem.usable", divMebibytes,
			kv{"state", "free"},
			kv{"state", "cached"},
			kv{"state", "buffered"},
		)
	case "system.network.io":
		copyMetric(all, m, "system.net.bytes_rcvd", 1, kv{"direction", "receive"})
		copyMetric(all, m, "system.net.bytes_sent", 1, kv{"direction", "transmit"})
	case "system.paging.usage":
		copyMetric(all, m, "system.swap.free", divMebibytes, kv{"state", "free"})
		copyMetric(all, m, "system.swap.used", divMebibytes, kv{"state", "used"})
	case "system.filesystem.utilization":
		copyMetric(all, m, "system.disk.in_use", 1)
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
		if addm, ok := copyMetric(all, m, "container.cpu.usage", 1); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.usage.usermode":
		if addm, ok := copyMetric(all, m, "container.cpu.user", 1); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.usage.system":
		if addm, ok := copyMetric(all, m, "container.cpu.system", 1); ok {
			addm.SetUnit("nanocore")
		}
	case "container.cpu.throttling_data.throttled_time":
		copyMetric(all, m, "container.cpu.throttled", 1)
	case "container.cpu.throttling_data.throttled_periods":
		copyMetric(all, m, "container.cpu.throttled.periods", 1)
	case "container.memory.usage.total":
		copyMetric(all, m, "container.memory.usage", 1)
	case "container.memory.active_anon":
		copyMetric(all, m, "container.memory.kernel", 1)
	case "container.memory.hierarchical_memory_limit":
		copyMetric(all, m, "container.memory.limit", 1)
	case "container.memory.usage.limit":
		copyMetric(all, m, "container.memory.soft_limit", 1)
	case "container.memory.total_cache":
		copyMetric(all, m, "container.memory.cache", 1)
	case "container.memory.total_swap":
		copyMetric(all, m, "container.memory.swap", 1)
	case "container.blockio.io_service_bytes_recursive":
		copyMetric(all, m, "container.io.write", 1, kv{"operation", "write"})
		copyMetric(all, m, "container.io.read", 1, kv{"operation", "read"})
	case "container.blockio.io_serviced_recursive":
		copyMetric(all, m, "container.io.write.operations", 1, kv{"operation", "write"})
		copyMetric(all, m, "container.io.read.operations", 1, kv{"operation", "read"})
	case "container.network.io.usage.tx_bytes":
		copyMetric(all, m, "container.net.sent", 1)
	case "container.network.io.usage.tx_packets":
		copyMetric(all, m, "container.net.sent.packets", 1)
	case "container.network.io.usage.rx_bytes":
		copyMetric(all, m, "container.net.rcvd", 1)
	case "container.network.io.usage.rx_packets":
		copyMetric(all, m, "container.net.rcvd.packets", 1)
	}
}

// kv represents a key/value pair.
type kv struct{ K, V string }

// copyMetric copies metric m to dest. The new metric's name will be newname, and all of its datapoints will
// be divided by div. If filter is provided, only the data points that have *either* of the specified string
// attributes will be copied over. If the filtering results in no datapoints, no new metric is added to dest.
//
// copyMetric returns the new metric and reports whether it was added to dest.
//
// Please note that copyMetric is restricted to the metric types Sum and Gauge.
func copyMetric(dest pmetric.MetricSlice, m pmetric.Metric, newname string, div float64, filter ...kv) (pmetric.Metric, bool) {
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
		// Rather than having metricAttrs map we can extend copyMetric to take in an attribute map for attributes
		// that need to be added to the new metric. To avoid cluttering this PR and chaging func signature, using
		// map.
		if t := metricAttrs[m.Name()]; t != "" {
			dp.Attributes().PutStr("type", t)
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
