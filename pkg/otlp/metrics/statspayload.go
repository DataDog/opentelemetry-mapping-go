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
	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/trace"
	"github.com/gogo/protobuf/jsonpb"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// keyStatsPayload is the key for the stats payload in the attributes map.
// This is used as Metric name and Attribute key.
const keyStatsPayload = "dd.internal.stats.payload"

var marshaler = &jsonpb.Marshaler{}

// UnsetHostnamePlaceholder is the string used as a hostname when the hostname can not be extracted from span attributes
// by the processor. Upon decoding the metrics, the Translator will use its configured fallback SourceProvider to replace
// it with the correct hostname.
//
// This isn't the most ideal approach to the problem, but provides the better user experience by avoiding the need to
// duplicate the "exporter::datadog::hostname" configuration field as "processor::datadog::hostname". The hostname can
// also not be left empty in case of failure to obtain it, because empty has special meaning. An empty hostname means
// that we are in a Lambda environment. Thus, we must use a placeholder.
const UnsetHostnamePlaceholder = "__unset__"

// keyAPMStats specifies the key name of the resource attribute which identifies resource metrics
// as being an APM Stats Payload. The presence of the key results in them being treated and consumed
// differently by the Translator.
const keyAPMStats = "_dd.apm_stats"

// StatsToMetrics converts a StatsPayload to a pdata.Metrics
func (t *Translator) StatsToMetrics(sp *pb.StatsPayload) (pmetric.Metrics, error) {
	payload, err := marshaler.MarshalToString(sp)
	if err != nil {
		t.logger.Error("Failed to marshal stats payload", zap.Error(err))
		return pmetric.NewMetrics(), err
	}
	mmx := pmetric.NewMetrics()
	rmx := mmx.ResourceMetrics().AppendEmpty()
	rmx.Resource().Attributes().PutBool(keyAPMStats, true)

	smx := rmx.ScopeMetrics().AppendEmpty()
	mslice := smx.Metrics()
	mx := mslice.AppendEmpty()
	mx.SetName(keyStatsPayload)
	sum := mx.SetEmptySum()
	sum.SetIsMonotonic(false)
	dp := sum.DataPoints().AppendEmpty()
	dp.Attributes().PutStr(keyStatsPayload, payload)
	return mmx, nil
}
