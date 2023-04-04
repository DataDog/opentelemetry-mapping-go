package metrics

// runtimeMetricMapping defines the fields needed to map OTel runtime metrics to their equivalent
// Datadog runtime metrics
type runtimeMetricMapping struct {
	mappedName string                   // the Datadog runtime metric name
	attributes []runtimeMetricAttribute // the attribute(s) this metric originates from
}

// runtimeMetricAttribute defines the structure for an attribute in regard to mapping runtime metrics.
// The presence of a runtimeMetricAttribute means that a metric must be mapped from a data point
// with the given attribute(s).
type runtimeMetricAttribute struct {
	key    string   // the attribute name
	values []string // the attribute value, or multiple values if there is more than one value for the same mapping
}

var goRuntimeMetricsMappings = map[string][]runtimeMetricMapping{
	"process.runtime.go.goroutines":        {{mappedName: "runtime.go.num_goroutine"}},
	"process.runtime.go.cgo.calls":         {{mappedName: "runtime.go.num_cgo_call"}},
	"process.runtime.go.lookups":           {{mappedName: "runtime.go.mem_stats.lookups"}},
	"process.runtime.go.mem.heap_alloc":    {{mappedName: "runtime.go.mem_stats.heap_alloc"}},
	"process.runtime.go.mem.heap_sys":      {{mappedName: "runtime.go.mem_stats.heap_sys"}},
	"process.runtime.go.mem.heap_idle":     {{mappedName: "runtime.go.mem_stats.heap_idle"}},
	"process.runtime.go.mem.heap_inuse":    {{mappedName: "runtime.go.mem_stats.heap_inuse"}},
	"process.runtime.go.mem.heap_released": {{mappedName: "runtime.go.mem_stats.heap_released"}},
	"process.runtime.go.mem.heap_objects":  {{mappedName: "runtime.go.mem_stats.heap_objects"}},
	"process.runtime.go.gc.pause_total_ns": {{mappedName: "runtime.go.mem_stats.pause_total_ns"}},
	"process.runtime.go.gc.count":          {{mappedName: "runtime.go.mem_stats.num_gc"}},
}

var dotnetRuntimeMetricsMappings = map[string][]runtimeMetricMapping{
	"process.runtime.dotnet.thread_pool.threads.count":     {{mappedName: "runtime.dotnet.threads.count"}},
	"process.runtime.dotnet.monitor.lock_contention.count": {{mappedName: "runtime.dotnet.threads.contention_count"}},
	"process.runtime.dotnet.exceptions.count":              {{mappedName: "runtime.dotnet.exceptions.count"}},
	"process.runtime.dotnet.gc.heap.size": {{
		mappedName: "runtime.dotnet.gc.size.gen0",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen0"},
		}},
	}, {
		mappedName: "runtime.dotnet.gc.size.gen1",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen1"},
		}},
	}, {
		mappedName: "runtime.dotnet.gc.size.gen2",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen2"},
		}},
	}, {
		mappedName: "runtime.dotnet.gc.size.loh",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"loh"},
		}},
	}},
	"process.runtime.dotnet.gc.collections.count": {{
		mappedName: "runtime.dotnet.gc.count.gen0",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen0"},
		}},
	}, {
		mappedName: "runtime.dotnet.gc.count.gen1",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen1"},
		}},
	}, {
		mappedName: "runtime.dotnet.gc.count.gen2",
		attributes: []runtimeMetricAttribute{{
			key:    "generation",
			values: []string{"gen2"},
		}},
	}},
}

var javaRuntimeMetricsMappings = map[string][]runtimeMetricMapping{
	"process.runtime.jvm.threads.count": {{mappedName: "jvm.thread_count"}},
	"process.runtime.jvm.gc.duration":   {{mappedName: "jvm.gc.parnew.time"}},
	"process.runtime.jvm.memory.usage": {{
		mappedName: "jvm.heap_memory",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"heap"},
		}},
	}, {
		mappedName: "jvm.non_heap_memory",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"non_heap"},
		}},
	}, {
		mappedName: "jvm.gc.old_gen_size",
		attributes: []runtimeMetricAttribute{{
			key:    "pool",
			values: []string{"g1_old_gen", "ps_old_gen", "tenured_gen"},
		}, {
			key:    "type",
			values: []string{"heap"},
		}},
	}},
	"process.runtime.jvm.memory.committed": {{
		mappedName: "jvm.heap_memory_committed",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"heap"},
		}},
	}, {
		mappedName: "jvm.non_heap_memory_committed",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.memory.init": {{
		mappedName: "jvm.heap_memory_init",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"heap"},
		}},
	}, {
		mappedName: "jvm.non_heap_memory_init",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.memory.limit": {{
		mappedName: "jvm.heap_memory_max",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"heap"},
		}},
	}, {
		mappedName: "jvm.non_heap_memory_max",
		attributes: []runtimeMetricAttribute{{
			key:    "type",
			values: []string{"non_heap"},
		}},
	}},
}

func getRuntimeMetricsMappings() map[string][]runtimeMetricMapping {
	res := map[string][]runtimeMetricMapping{}
	for k, v := range goRuntimeMetricsMappings {
		res[k] = v
	}
	for k, v := range dotnetRuntimeMetricsMappings {
		res[k] = v
	}
	for k, v := range javaRuntimeMetricsMappings {
		res[k] = v
	}
	return res
}

// runtimeMetricsMappings defines the mappings from OTel runtime metric names to their
// equivalent Datadog runtime metric names
var runtimeMetricsMappings = getRuntimeMetricsMappings()
