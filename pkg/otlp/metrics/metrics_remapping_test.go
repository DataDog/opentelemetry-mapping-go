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

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/stretchr/testify/assert"
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

		// kafka
		{
			in:  metric("kafka.producer.request-rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.producer.request_rate", point{f: 1, attrs: map[string]any{"type": "producer-metrics"}})},
		},
		{
			in:  metric("kafka.producer.response-rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.producer.response_rate", point{f: 1, attrs: map[string]any{"type": "producer-metrics"}})},
		},
		{
			in:  metric("kafka.producer.request-latency-avg", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.producer.request_latency_avg", point{f: 1, attrs: map[string]any{"type": "producer-metrics"}})},
		},
		{
			in:  metric("kafka.producer.outgoing-byte-rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.producer.bytes_out", point{f: 1, attrs: map[string]any{"type": "producer-metrics"}})},
		},
		{
			in:  metric("kafka.producer.io-wait-time-ns-avg", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.producer.io_wait", point{f: 1, attrs: map[string]any{"type": "producer-metrics"}})},
		},
		{
			in: metric("kafka.producer.byte-rate", point{f: 1, attrs: map[string]any{"client-id": "client123"}}),
			out: []pmetric.Metric{metric("kafka.producer.bytes_out", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
				"client":    "client123",
				"type":      "producer-topic-metrics",
			}})},
		},
		{
			in:  metric("kafka.consumer.total.bytes-consumed-rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.consumer.bytes_in", point{f: 1, attrs: map[string]any{"type": "consumer-fetch-manager-metrics"}})},
		},
		{
			in:  metric("kafka.consumer.total.records-consumed-rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.consumer.messages_in", point{f: 1, attrs: map[string]any{"type": "consumer-fetch-manager-metrics"}})},
		},
		{
			in: metric("kafka.network.io",
				point{f: 1, attrs: map[string]any{
					"state": "out",
				}},
				point{f: 2, attrs: map[string]any{
					"state": "in",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.net.bytes_out.rate", point{f: 1, attrs: map[string]any{
					"type":  "BrokerTopicMetrics",
					"name":  "BytesOutPerSec",
					"state": "out",
				}}),
				metric("kafka.net.bytes_in.rate", point{f: 2, attrs: map[string]any{
					"type":  "BrokerTopicMetrics",
					"name":  "BytesInPerSec",
					"state": "in",
				}}),
			},
		},
		{
			in: metric("kafka.purgatory.size",
				point{f: 1, attrs: map[string]any{
					"type": "produce",
				}},
				point{f: 2, attrs: map[string]any{
					"type": "fetch",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.request.producer_request_purgatory.size", point{f: 1, attrs: map[string]any{
					"type":             "DelayedOperationPurgatory",
					"name":             "PurgatorySize",
					"delayedOperation": "Produce",
				}}),
				metric("kafka.request.fetch_request_purgatory.size", point{f: 2, attrs: map[string]any{
					"type":             "DelayedOperationPurgatory",
					"name":             "PurgatorySize",
					"delayedOperation": "Fetch",
				}}),
			},
		},
		{
			in: metric("kafka.partition.under_replicated", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.under_replicated_partitions", point{f: 1, attrs: map[string]any{
				"type": "ReplicaManager",
				"name": "UnderReplicatedPartitions",
			}})},
		},
		{
			in: metric("kafka.isr.operation.count",
				point{f: 1, attrs: map[string]any{
					"operation": "shrink",
				}},
				point{f: 2, attrs: map[string]any{
					"operation": "expand",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.replication.isr_shrinks.rate", point{f: 1, attrs: map[string]any{
					"type":      "ReplicaManager",
					"name":      "IsrShrinksPerSec",
					"operation": "shrink",
				}}),
				metric("kafka.replication.isr_expands.rate", point{f: 2, attrs: map[string]any{
					"type":      "ReplicaManager",
					"name":      "IsrExpandsPerSec",
					"operation": "expand",
				}}),
			},
		},
		{
			in: metric("kafka.leader.election.rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.leader_elections.rate", point{f: 1, attrs: map[string]any{
				"type": "ControllerStats",
				"name": "LeaderElectionRateAndTimeMs",
			}})},
		},
		{
			in: metric("kafka.partition.offline", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.offline_partitions_count", point{f: 1, attrs: map[string]any{
				"type": "KafkaController",
				"name": "OfflinePartitionsCount",
			}})},
		},
		{
			in: metric("kafka.request.time.avg",
				point{f: 1, attrs: map[string]any{
					"type": "produce",
				}},
				point{f: 2, attrs: map[string]any{
					"type": "fetchconsumer",
				}},
				point{f: 3, attrs: map[string]any{
					"type": "fetchfollower",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.request.produce.time.avg", point{f: 1, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "Produce",
				}}),
				metric("kafka.request.fetch_consumer.time.avg", point{f: 2, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "FetchConsumer",
				}}),
				metric("kafka.request.fetch_follower.time.avg", point{f: 3, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "FetchFollower",
				}}),
			},
		},
		{
			in: metric("kafka.message.count", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.messages_in.rate", point{f: 1, attrs: map[string]any{
				"type": "BrokerTopicMetrics",
				"name": "MessagesInPerSec",
			}})},
		},
		{
			in: metric("kafka.request.failed",
				point{f: 1, attrs: map[string]any{
					"type": "produce",
				}},
				point{f: 2, attrs: map[string]any{
					"type": "fetch",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.request.produce.failed.rate", point{f: 1, attrs: map[string]any{
					"type": "BrokerTopicMetrics",
					"name": "FailedProduceRequestsPerSec",
				}}),
				metric("kafka.request.fetch.failed.rate", point{f: 2, attrs: map[string]any{
					"type": "BrokerTopicMetrics",
					"name": "FailedFetchRequestsPerSec",
				}}),
			},
		},
		{
			in: metric("kafka.request.time.99p",
				point{f: 1, attrs: map[string]any{
					"type": "produce",
				}},
				point{f: 2, attrs: map[string]any{
					"type": "fetchconsumer",
				}},
				point{f: 3, attrs: map[string]any{
					"type": "fetchfollower",
				}},
			),
			out: []pmetric.Metric{
				metric("kafka.request.produce.time.99percentile", point{f: 1, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "Produce",
				}}),
				metric("kafka.request.fetch_consumer.time.99percentile", point{f: 2, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "FetchConsumer",
				}}),
				metric("kafka.request.fetch_follower.time.99percentile", point{f: 3, attrs: map[string]any{
					"type":    "RequestMetrics",
					"name":    "TotalTimeMs",
					"request": "FetchFollower",
				}}),
			},
		},
		{
			in: metric("kafka.partition.count", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.partition_count", point{f: 1, attrs: map[string]any{
				"type": "ReplicaManager",
				"name": "PartitionCount",
			}})},
		},
		{
			in: metric("kafka.max.lag", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.max_lag", point{f: 1, attrs: map[string]any{
				"type":     "ReplicaFetcherManager",
				"name":     "MaxLag",
				"clientId": "replica",
			}})},
		},
		{
			in: metric("kafka.controller.active.count", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.active_controller_count", point{f: 1, attrs: map[string]any{
				"type": "KafkaController",
				"name": "ActiveControllerCount",
			}})},
		},
		{
			in: metric("kafka.unclean.election.rate", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.replication.unclean_leader_elections.rate", point{f: 1, attrs: map[string]any{
				"type": "ControllerStats",
				"name": "UncleanLeaderElectionsPerSec",
			}})},
		},
		{
			in: metric("kafka.request.queue", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.request.channel.queue.size", point{f: 1, attrs: map[string]any{
				"type": "RequestChannel",
				"name": "RequestQueueSize",
			}})},
		},
		{
			in: metric("kafka.logs.flush.time.count", point{f: 1}),
			out: []pmetric.Metric{metric("kafka.log.flush_rate.rate", point{f: 1, attrs: map[string]any{
				"type": "LogFlushStats",
				"name": "LogFlushRateAndTimeMs",
			}})},
		},

		{
			in: metric("kafka.consumer.bytes-consumed-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.consumer.bytes_consumed", point{f: 1, attrs: map[string]any{
				"type":      "consumer-fetch-manager-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.consumer.records-consumed-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.consumer.records_consumed", point{f: 1, attrs: map[string]any{
				"type":      "consumer-fetch-manager-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.consumer.fetch-size-avg", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.consumer.fetch_size_avg", point{f: 1, attrs: map[string]any{
				"type":      "consumer-fetch-manager-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.producer.compression-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.producer.compression_rate", point{f: 1, attrs: map[string]any{
				"type":      "producer-topic-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.producer.record-error-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.producer.record_error_rate", point{f: 1, attrs: map[string]any{
				"type":      "producer-topic-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.producer.record-retry-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.producer.record_retry_rate", point{f: 1, attrs: map[string]any{
				"type":      "producer-topic-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},
		{
			in: metric("kafka.producer.record-send-rate", point{f: 1, attrs: map[string]any{
				"client-id": "client123",
			}}),
			out: []pmetric.Metric{metric("kafka.producer.record_send_rate", point{f: 1, attrs: map[string]any{
				"type":      "producer-topic-metrics",
				"client-id": "client123",
				"client":    "client123",
			}})},
		},

		// kafka metrics receiver
		{
			in: metric("kafka.partition.current_offset", point{f: 1, attrs: map[string]any{
				"group": "group123",
			}}),
			out: []pmetric.Metric{metric("kafka.broker_offset", point{f: 1, attrs: map[string]any{
				"group":          "group123",
				"consumer_group": "group123",
			}})},
		},
		{
			in: metric("kafka.consumer_group.lag", point{f: 1, attrs: map[string]any{
				"group": "group123",
			}}),
			out: []pmetric.Metric{metric("kafka.consumer_lag", point{f: 1, attrs: map[string]any{
				"group":          "group123",
				"consumer_group": "group123",
			}})},
		},
		{
			in: metric("kafka.consumer_group.offset", point{f: 1, attrs: map[string]any{
				"group": "group123",
			}}),
			out: []pmetric.Metric{metric("kafka.consumer_offset", point{f: 1, attrs: map[string]any{
				"group":          "group123",
				"consumer_group": "group123",
			}})},
		},

		// jvm
		{
			in: metric("jvm.gc.collections.count",
				point{f: 1, attrs: map[string]any{"name": "Copy"}},
				point{f: 2, attrs: map[string]any{"name": "PS Scavenge"}},
				point{f: 3, attrs: map[string]any{"name": "ParNew"}},
				point{f: 4, attrs: map[string]any{"name": "G1 Young Generation"}},
			),
			out: []pmetric.Metric{
				metric("jvm.gc.minor_collection_count",
					point{f: 1, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "Copy",
					}},
					point{f: 2, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "PS Scavenge",
					}},
					point{f: 3, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ParNew",
					}},
					point{f: 4, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Young Generation",
					}}),
			},
		},
		{
			in: metric("jvm.gc.collections.count",
				point{f: 1, attrs: map[string]any{"name": "MarkSweepCompact"}},
				point{f: 2, attrs: map[string]any{"name": "PS MarkSweep"}},
				point{f: 3, attrs: map[string]any{"name": "ConcurrentMarkSweep"}},
				point{f: 4, attrs: map[string]any{"name": "G1 Mixed Generation"}},
				point{f: 5, attrs: map[string]any{"name": "G1 Old Generation"}},
				point{f: 6, attrs: map[string]any{"name": "Shenandoah Cycles"}},
				point{f: 7, attrs: map[string]any{"name": "ZGC"}},
			),
			out: []pmetric.Metric{
				metric("jvm.gc.major_collection_count",
					point{f: 1, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "MarkSweepCompact",
					}},
					point{f: 2, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "PS MarkSweep",
					}},
					point{f: 3, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ConcurrentMarkSweep",
					}},
					point{f: 4, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Mixed Generation",
					}},
					point{f: 5, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Old Generation",
					}},
					point{f: 6, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "Shenandoah Cycles",
					}},
					point{f: 7, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ZGC",
					}}),
			},
		},
		{
			in: metric("jvm.gc.collections.elapsed",
				point{f: 1, attrs: map[string]any{"name": "Copy"}},
				point{f: 2, attrs: map[string]any{"name": "PS Scavenge"}},
				point{f: 3, attrs: map[string]any{"name": "ParNew"}},
				point{f: 4, attrs: map[string]any{"name": "G1 Young Generation"}},
			),
			out: []pmetric.Metric{
				metric("jvm.gc.minor_collection_time",
					point{f: 1, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "Copy",
					}},
					point{f: 2, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "PS Scavenge",
					}},
					point{f: 3, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ParNew",
					}},
					point{f: 4, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Young Generation",
					}}),
			},
		},
		{
			in: metric("jvm.gc.collections.elapsed",
				point{f: 1, attrs: map[string]any{"name": "MarkSweepCompact"}},
				point{f: 2, attrs: map[string]any{"name": "PS MarkSweep"}},
				point{f: 3, attrs: map[string]any{"name": "ConcurrentMarkSweep"}},
				point{f: 4, attrs: map[string]any{"name": "G1 Mixed Generation"}},
				point{f: 5, attrs: map[string]any{"name": "G1 Old Generation"}},
				point{f: 6, attrs: map[string]any{"name": "Shenandoah Cycles"}},
				point{f: 7, attrs: map[string]any{"name": "ZGC"}},
			),
			out: []pmetric.Metric{
				metric("jvm.gc.major_collection_time",
					point{f: 1, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "MarkSweepCompact",
					}},
					point{f: 2, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "PS MarkSweep",
					}},
					point{f: 3, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ConcurrentMarkSweep",
					}},
					point{f: 4, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Mixed Generation",
					}},
					point{f: 5, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "G1 Old Generation",
					}},
					point{f: 6, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "Shenandoah Cycles",
					}},
					point{f: 7, attrs: map[string]any{
						"type": "GarbageCollector",
						"name": "ZGC",
					}}),
			},
		},
	} {
		lena := dest.Len()
		checkprefix := strings.HasPrefix(tt.in.Name(), "system.") ||
			strings.HasPrefix(tt.in.Name(), "process.") ||

			tt.in.Name() == "kafka.producer.request-rate" ||
			tt.in.Name() == "kafka.producer.response-rate" ||
			tt.in.Name() == "kafka.producer.request-latency-avg" ||

			tt.in.Name() == "kafka.consumer.fetch-size-avg" ||
			tt.in.Name() == "kafka.producer.compression-rate" ||
			tt.in.Name() == "kafka.producer.record-error-rate" ||
			tt.in.Name() == "kafka.producer.record-retry-rate" ||
			tt.in.Name() == "kafka.producer.record-send-rate"
		remapMetrics(dest, tt.in)
		if checkprefix {
			require.True(t, strings.HasPrefix(tt.in.Name(), "otel."), "system.* and process.*  and a subset of kafka metrics need to be prepended with the otel.* namespace")
		}
		require.Equal(t, dest.Len()-lena, len(tt.out), "unexpected number of metrics added")
		for i, out := range tt.out {
			assert.NoError(t, pmetrictest.CompareMetric(out, dest.At(dest.Len()-len(tt.out)+i)))
		}
	}

}

func TestCopyMetricWithAttr(t *testing.T) {
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
			out, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{})
			require.True(t, ok)
			require.Equal(t, m.Name(), "test.metric")
			require.Equal(t, out.Name(), "copied.test.metric")
			sameExceptName(t, m, out)
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("div", func(t *testing.T) {
			out, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 2, attributesMapping{})
			require.True(t, ok)
			require.Equal(t, out.Name(), "copied.test.metric")
			require.Equal(t, out.Gauge().DataPoints().At(0).DoubleValue(), 6.)
			require.Equal(t, out.Gauge().DataPoints().At(1).IntValue(), int64(12))
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("filter", func(t *testing.T) {
			out, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{}, kv{"human", "Ann"})
			require.True(t, ok)
			require.Equal(t, out.Name(), "copied.test.metric")
			require.Equal(t, out.Gauge().DataPoints().Len(), 1)
			require.Equal(t, out.Gauge().DataPoints().At(0).IntValue(), int64(24))
			require.Equal(t, out.Gauge().DataPoints().At(0).Attributes().AsRaw(), map[string]any{"human": "Ann", "age": int64(25)})
			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("attributesMapping", func(t *testing.T) {
			out, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{
				fixed:   map[string]string{"fixed.attr": "ok"},
				dynamic: map[string]string{"fruit": "remapped_fruit"},
			})
			require.True(t, ok)
			require.Equal(t, m.Name(), "test.metric")
			require.Equal(t, out.Name(), "copied.test.metric")

			aa, bb := pmetric.NewMetric(), pmetric.NewMetric()
			m.CopyTo(aa)
			out.CopyTo(bb)

			aa.SetName("common.name")
			// add attributes mappings manually.
			aa.Gauge().DataPoints().At(0).Attributes().PutStr("fixed.attr", "ok")
			aa.Gauge().DataPoints().At(0).Attributes().PutStr("remapped_fruit", "apple")
			aa.Gauge().DataPoints().At(1).Attributes().PutStr("fixed.attr", "ok")

			bb.SetName("common.name")
			require.Equal(t, aa, bb)

			require.Equal(t, dest.At(dest.Len()-1), out)
		})

		t.Run("none", func(t *testing.T) {
			_, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{}, kv{"human", "Paul"})
			require.False(t, ok)
		})
	})

	t.Run("sum", func(t *testing.T) {
		dp := m.SetEmptySum().DataPoints().AppendEmpty()
		dp.SetDoubleValue(12)
		dp.Attributes().FromRaw(map[string]any{"fruit": "apple", "count": 15})
		out, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{})
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
		_, ok := copyMetricWithAttr(dest, m, "copied.test.metric", 1, attributesMapping{})
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
