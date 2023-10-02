package metricscommon

// RuntimeMetricPrefixLanguageMap defines the runtime metric prefixes and which languages they map to
var RuntimeMetricPrefixLanguageMap = map[string]string{
	"process.runtime.go":     "go",
	"process.runtime.dotnet": "dotnet",
	"process.runtime.jvm":    "jvm",
}

// RuntimeMetricMapping defines the fields needed to map OTel runtime metrics to their equivalent
// Datadog runtime metrics
type RuntimeMetricMapping struct {
	MappedName string                   // the Datadog runtime metric name
	Attributes []RuntimeMetricAttribute // the attribute(s) this metric originates from
}

// RuntimeMetricAttribute defines the structure for an attribute in regard to mapping runtime metrics.
// The presence of a RuntimeMetricAttribute means that a metric must be mapped from a data point
// with the given attribute(s).
type RuntimeMetricAttribute struct {
	Key    string   // the attribute name
	Values []string // the attribute value, or multiple Values if there is more than one value for the same mapping
}

// RuntimeMetricMappingList defines the structure for a list of runtime metric mappings where the Key
// represents the OTel metric name and the RuntimeMetricMapping contains the Datadog metric name
type RuntimeMetricMappingList map[string][]RuntimeMetricMapping

var goRuntimeMetricsMappings = RuntimeMetricMappingList{
	"process.runtime.go.goroutines":        {{MappedName: "runtime.go.num_goroutine"}},
	"process.runtime.go.cgo.calls":         {{MappedName: "runtime.go.num_cgo_call"}},
	"process.runtime.go.lookups":           {{MappedName: "runtime.go.mem_stats.lookups"}},
	"process.runtime.go.mem.heap_alloc":    {{MappedName: "runtime.go.mem_stats.heap_alloc"}},
	"process.runtime.go.mem.heap_sys":      {{MappedName: "runtime.go.mem_stats.heap_sys"}},
	"process.runtime.go.mem.heap_idle":     {{MappedName: "runtime.go.mem_stats.heap_idle"}},
	"process.runtime.go.mem.heap_inuse":    {{MappedName: "runtime.go.mem_stats.heap_inuse"}},
	"process.runtime.go.mem.heap_released": {{MappedName: "runtime.go.mem_stats.heap_released"}},
	"process.runtime.go.mem.heap_objects":  {{MappedName: "runtime.go.mem_stats.heap_objects"}},
	"process.runtime.go.gc.pause_total_ns": {{MappedName: "runtime.go.mem_stats.pause_total_ns"}},
	"process.runtime.go.gc.count":          {{MappedName: "runtime.go.mem_stats.num_gc"}},
}

var dotnetRuntimeMetricsMappings = RuntimeMetricMappingList{
	"process.runtime.dotnet.monitor.lock_contention.count": {{MappedName: "runtime.dotnet.threads.contention_count"}},
	"process.runtime.dotnet.exceptions.count":              {{MappedName: "runtime.dotnet.exceptions.count"}},
	"process.runtime.dotnet.gc.heap.size": {{
		MappedName: "runtime.dotnet.gc.size.gen0",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen0"},
		}},
	}, {
		MappedName: "runtime.dotnet.gc.size.gen1",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen1"},
		}},
	}, {
		MappedName: "runtime.dotnet.gc.size.gen2",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen2"},
		}},
	}, {
		MappedName: "runtime.dotnet.gc.size.loh",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"loh"},
		}},
	}},
	"process.runtime.dotnet.gc.collections.count": {{
		MappedName: "runtime.dotnet.gc.count.gen0",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen0"},
		}},
	}, {
		MappedName: "runtime.dotnet.gc.count.gen1",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen1"},
		}},
	}, {
		MappedName: "runtime.dotnet.gc.count.gen2",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "generation",
			Values: []string{"gen2"},
		}},
	}},
}

var javaRuntimeMetricsMappings = RuntimeMetricMappingList{
	"process.runtime.jvm.threads.count":          {{MappedName: "jvm.thread_count"}},
	"process.runtime.jvm.classes.current_loaded": {{MappedName: "jvm.loaded_classes"}},
	"process.runtime.jvm.system.cpu.utilization": {{MappedName: "jvm.cpu_load.system"}},
	"process.runtime.jvm.cpu.utilization":        {{MappedName: "jvm.cpu_load.process"}},
	"process.runtime.jvm.gc.duration":            {{MappedName: "jvm.gc.parnew.time"}},
	"process.runtime.jvm.memory.usage": {{
		MappedName: "jvm.heap_memory",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.non_heap_memory",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"non_heap"},
		}},
	}, {
		MappedName: "jvm.gc.old_gen_size",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"G1 Old Gen", "Tenured Gen", "PS Old Gen"},
		}, {
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.gc.eden_size",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"G1 Eden Space", "Eden Space", "Par Eden Space", "PS Eden Space"},
		}, {
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.gc.survivor_size",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"G1 Survivor Space", "Survivor Space", "Par Survivor Space", "PS Survivor Space"},
		}, {
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.gc.metaspace_size",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"Metaspace"},
		}, {
			Key:    "type",
			Values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.memory.committed": {{
		MappedName: "jvm.heap_memory_committed",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.non_heap_memory_committed",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.memory.init": {{
		MappedName: "jvm.heap_memory_init",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.non_heap_memory_init",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.memory.limit": {{
		MappedName: "jvm.heap_memory_max",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"heap"},
		}},
	}, {
		MappedName: "jvm.non_heap_memory_max",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "type",
			Values: []string{"non_heap"},
		}},
	}},
	"process.runtime.jvm.buffer.usage": {{
		MappedName: "jvm.buffer_pool.direct.used",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"direct"},
		}},
	}, {
		MappedName: "jvm.buffer_pool.mapped.used",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"mapped"},
		}},
	}},
	"process.runtime.jvm.buffer.count": {{
		MappedName: "jvm.buffer_pool.direct.count",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"direct"},
		}},
	}, {
		MappedName: "jvm.buffer_pool.mapped.count",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"mapped"},
		}},
	}},
	"process.runtime.jvm.buffer.limit": {{
		MappedName: "jvm.buffer_pool.direct.limit",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"direct"},
		}},
	}, {
		MappedName: "jvm.buffer_pool.mapped.limit",
		Attributes: []RuntimeMetricAttribute{{
			Key:    "pool",
			Values: []string{"mapped"},
		}},
	}},
}

func getRuntimeMetricsMappings() RuntimeMetricMappingList {
	res := RuntimeMetricMappingList{}
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

// RuntimeMetricsMappings defines the mappings from OTel runtime metric names to their
// equivalent Datadog runtime metric names
var RuntimeMetricsMappings = getRuntimeMetricsMappings()
