package rum

import (
	"net/http"

	"go.opentelemetry.io/collector/pdata/plog"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
)

func ToLogs(payload map[string]any, req *http.Request) plog.Logs {
	results := plog.NewLogs()
	rl := results.ResourceLogs().AppendEmpty()
	rl.SetSchemaUrl(semconv.SchemaURL)
	parseRUMRequestIntoResource(rl.Resource(), payload, req.URL.Query().Get("ddforward"))

	in := rl.ScopeLogs().AppendEmpty()
	in.Scope().SetName(InstrumentationScopeName)

	newLogRecord := in.LogRecords().AppendEmpty()

	flatPayload := flattenJSON(payload)

	setOTLPAttributes(flatPayload, newLogRecord.Attributes())

	return results
}
