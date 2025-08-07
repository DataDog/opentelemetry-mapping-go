// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package rum

var OTLPAttributeToRUMPayloadKeyMapping = map[string]string{
	// _common-schema.json
	Service + "." + Name:                                    Service,
	Service + "." + Version:                                 Version,
	Session + "." + Id:                                      Session + "." + Id,
	User + "." + Id:                                         Usr + "." + Id,
	User + "." + FullName:                                   Usr + "." + Name,
	User + "." + Email:                                      Usr + "." + Email,
	User + "." + Hash:                                       Usr + "." + AnonymousId,
	User + "." + Name:                                       Account + "." + Name,
	Os + "." + Name:                                         Os + "." + Name,
	Os + "." + Version:                                      Os + "." + Version,
	Os + "." + Build:                                        Os + "." + BuildId,
	Device + "." + Model + "." + Name:                       Device + "." + Name,
	Device + "." + Model + "." + Identifier:                 Device + "." + Model,
	Device + "." + Manufacturer:                             Device + "." + Brand,
	Datadog + "." + FormatVersion:                           DD + "." + FormatVersion,
	Datadog + "." + Session + "." + Plan:                    DD + "." + Session + "." + Plan,
	Datadog + "." + Session + "." + SessionPrecondition:     DD + "." + Session + "." + SessionPrecondition,
	Datadog + "." + Configuration + "." + SessionSampleRate: DD + "." + Configuration + "." + SessionSampleRate,
	Datadog + "." + Configuration + "." + SessionReplaySampleRate: DD + "." + Configuration + "." + SessionReplaySampleRate,
	Datadog + "." + Configuration + "." + ProfilingSampleRate:     DD + "." + Configuration + "." + ProfilingSampleRate,
	Datadog + "." + BrowserSDKVersion:                             DD + "." + BrowserSDKVersion,
	Datadog + "." + SDKName:                                       DD + "." + SDKName,

	// action-schema.json
	Datadog + "." + Action + "." + Position + "." + X:      DD + "." + Action + "." + Position + "." + X,
	Datadog + "." + Action + "." + Position + "." + Y:      DD + "." + Action + "." + Position + "." + Y,
	Datadog + "." + Action + "." + Target + "." + Selector: DD + "." + Action + "." + Target + "." + Selector,
	Datadog + "." + Action + "." + Target + "." + Width:    DD + "." + Action + "." + Target + "." + Width,
	Datadog + "." + Action + "." + Target + "." + Height:   DD + "." + Action + "." + Target + "." + Height,
	Datadog + "." + Action + "." + NameSource:              DD + "." + Action + "." + NameSource,

	// error-schema.json
	Error + "." + Message: Error + "." + Message,
	Error + "." + Type:    Error + "." + Type,

	// long_task-schema.json
	Datadog + "." + Discarded: DD + "." + Discarded,
	Datadog + "." + Profiling: DD + "." + Profiling,

	// resource-schema.json
	Datadog + "." + SpanId:                        DD + "." + SpanId,
	Datadog + "." + ParentSpanId:                  DD + "." + ParentSpanId,
	Datadog + "." + TraceId:                       DD + "." + TraceId,
	Datadog + "." + RulePSR:                       DD + "." + RulePSR,
	Datadog + "." + Profiling + "." + Status:      DD + "." + Profiling + "." + Status,
	Datadog + "." + Profiling + "." + ErrorReason: DD + "." + Profiling + "." + ErrorReason,

	// _view-schema.json
	Datadog + "." + DocumentVersion:                                           DD + "." + DocumentVersion,
	Datadog + "." + PageStates:                                                DD + "." + PageStates,
	Datadog + "." + PageStates + "." + State:                                  DD + "." + PageStates + "." + State,
	Datadog + "." + PageStates + "." + StartTime:                              DD + "." + PageStates + "." + StartTime,
	Datadog + "." + ReplayStats + "." + RecordsCount:                          DD + "." + ReplayStats + "." + RecordsCount,
	Datadog + "." + ReplayStats + "." + SegmentsCount:                         DD + "." + ReplayStats + "." + SegmentsCount,
	Datadog + "." + ReplayStats + "." + SegmentsTotalRawSize:                  DD + "." + ReplayStats + "." + SegmentsTotalRawSize,
	Datadog + "." + CLS + "." + DevicePixelRatio:                              DD + "." + CLS + "." + DevicePixelRatio,
	Datadog + "." + Configuration + "." + StartSessionReplayRecordingManually: DD + "." + Configuration + "." + StartSessionReplayRecordingManually,

	// vitals-schema.json
	Datadog + "." + Vital + "." + ComputedValue: DD + "." + Vital + "." + ComputedValue,
}

var RUMPayloadKeyToOTLPAttributeMapping = map[string]string{
	// _common-schema.json
	Service:                         Service + "." + Name,
	Version:                         Service + "." + Version,
	Session + "." + Id:              Session + "." + Id,
	Usr + "." + Id:                  User + "." + Id,
	Usr + "." + Name:                User + "." + FullName,
	Usr + "." + Email:               User + "." + Email,
	Usr + "." + AnonymousId:         User + "." + Hash,
	Account + "." + Name:            User + "." + Name,
	Os + "." + Name:                 Os + "." + Name,
	Os + "." + Version:              Os + "." + Version,
	Os + "." + Build:                Os + "." + BuildId,
	Device + "." + Name:             Device + "." + Model + "." + Name,
	Device + "." + Model:            Device + "." + Model + "." + Identifier,
	Device + "." + Brand:            Device + "." + Manufacturer,
	DD + "." + FormatVersion:        Datadog + "." + FormatVersion,
	DD + "." + Session + "." + Plan: Datadog + "." + Session + "." + Plan,
	DD + "." + Session + "." + SessionPrecondition:           Datadog + "." + Session + "." + SessionPrecondition,
	DD + "." + Configuration + "." + SessionSampleRate:       Datadog + "." + Configuration + "." + SessionSampleRate,
	DD + "." + Configuration + "." + SessionReplaySampleRate: Datadog + "." + Configuration + "." + SessionReplaySampleRate,
	DD + "." + Configuration + "." + ProfilingSampleRate:     Datadog + "." + Configuration + "." + ProfilingSampleRate,
	DD + "." + BrowserSDKVersion:                             Datadog + "." + BrowserSDKVersion,
	DD + "." + SDKName:                                       Datadog + "." + SDKName,

	// action-schema.json
	DD + "." + Action + "." + Position + "." + X:      Datadog + "." + Action + "." + Position + "." + X,
	DD + "." + Action + "." + Position + "." + Y:      Datadog + "." + Action + "." + Position + "." + Y,
	DD + "." + Action + "." + Target + "." + Selector: Datadog + "." + Action + "." + Target + "." + Selector,
	DD + "." + Action + "." + Target + "." + Width:    Datadog + "." + Action + "." + Target + "." + Width,
	DD + "." + Action + "." + Target + "." + Height:   Datadog + "." + Action + "." + Target + "." + Height,
	DD + "." + Action + "." + NameSource:              Datadog + "." + Action + "." + NameSource,

	// error-schema.json
	Error + "." + Message: Error + "." + Message,
	Error + "." + Type:    Error + "." + Type,

	// long_task-schema.json
	DD + "." + Discarded: Datadog + "." + Discarded,
	DD + "." + Profiling: Datadog + "." + Profiling,

	// resource-schema.json
	DD + "." + SpanId:                        Datadog + "." + SpanId,
	DD + "." + ParentSpanId:                  Datadog + "." + ParentSpanId,
	DD + "." + TraceId:                       Datadog + "." + TraceId,
	DD + "." + RulePSR:                       Datadog + "." + RulePSR,
	DD + "." + Profiling + "." + Status:      Datadog + "." + Profiling + "." + Status,
	DD + "." + Profiling + "." + ErrorReason: Datadog + "." + Profiling + "." + ErrorReason,

	// _view-schema.json
	DD + "." + DocumentVersion:                                           Datadog + "." + DocumentVersion,
	DD + "." + PageStates:                                                Datadog + "." + PageStates,
	DD + "." + PageStates + "." + State:                                  Datadog + "." + PageStates + "." + State,
	DD + "." + PageStates + "." + StartTime:                              Datadog + "." + PageStates + "." + StartTime,
	DD + "." + ReplayStats + "." + RecordsCount:                          Datadog + "." + ReplayStats + "." + RecordsCount,
	DD + "." + ReplayStats + "." + SegmentsCount:                         Datadog + "." + ReplayStats + "." + SegmentsCount,
	DD + "." + ReplayStats + "." + SegmentsTotalRawSize:                  Datadog + "." + ReplayStats + "." + SegmentsTotalRawSize,
	DD + "." + CLS + "." + DevicePixelRatio:                              Datadog + "." + CLS + "." + DevicePixelRatio,
	DD + "." + Configuration + "." + StartSessionReplayRecordingManually: Datadog + "." + Configuration + "." + StartSessionReplayRecordingManually,

	// vitals-schema.json
	DD + "." + Vital + "." + ComputedValue: Datadog + "." + Vital + "." + ComputedValue,
}
