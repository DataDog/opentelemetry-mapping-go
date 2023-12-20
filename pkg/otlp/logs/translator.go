// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logs

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
)

// Translator of OTLP logs to Datadog format
type Translator struct {
	set     component.TelemetrySettings
	otelTag string
}

// NewTranslator returns a new Translator
func NewTranslator(set component.TelemetrySettings, otelSource string) (*Translator, error) {
	return &Translator{
		set:     set,
		otelTag: "otel_source:" + otelSource,
	}, nil
}

// MapLogs from OTLP format to Datadog format.
func (t *Translator) MapLogs(ld plog.Logs) []datadogV2.HTTPLogItem {
	rsl := ld.ResourceLogs()
	var payloads []datadogV2.HTTPLogItem
	for i := 0; i < rsl.Len(); i++ {
		rl := rsl.At(i)
		sls := rl.ScopeLogs()
		res := rl.Resource()
		for j := 0; j < sls.Len(); j++ {
			sl := sls.At(j)
			lsl := sl.LogRecords()
			// iterate over Logs
			for k := 0; k < lsl.Len(); k++ {
				log := lsl.At(k)
				payload := Transform(log, res, t.set.Logger)
				ddtags := payload.GetDdtags()
				if ddtags != "" {
					payload.SetDdtags(ddtags + "," + t.otelTag)
				} else {
					payload.SetDdtags(t.otelTag)
				}
				payloads = append(payloads, payload)
			}
		}
	}
	return payloads
}
