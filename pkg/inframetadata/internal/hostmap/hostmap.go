// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package hostmap

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"
	"go.uber.org/multierr"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/gohai"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/payload"
)

// HostMap maps from hostnames to host metadata payloads.
type HostMap struct {
	// mu is the mutex for the host map and updater.
	mu sync.Mutex
	// hosts map
	hosts map[string]payload.HostMetadata
}

// New creates a new HostMap.
func New() (*HostMap, error) {
	return &HostMap{
		mu:    sync.Mutex{},
		hosts: make(map[string]payload.HostMetadata),
	}, nil
}

func getStrField(m pcommon.Map, key string) (string, bool, error) {
	val, ok := m.Get(key)
	if !ok {
		// Field not available, don't update but don't fail either
		return "", false, nil
	}

	if val.Type() != pcommon.ValueTypeStr {
		return "", false, fmt.Errorf("%q has type %q, expected type \"Str\" instead", key, val.Type())
	}

	return val.Str(), true, nil
}

func isAWS(m pcommon.Map) (bool, error) {
	cloudProvider, ok, err := getStrField(m, conventions.AttributeCloudProvider)
	if err != nil {
		return false, err
	} else if !ok {
		// no cloud provider field
		return false, nil
	}
	return cloudProvider == conventions.AttributeCloudProviderAWS, nil
}

func getInstanceID(m pcommon.Map) (string, bool, error) {
	if onAWS, err := isAWS(m); err != nil || !onAWS {
		return "", onAWS, err
	}
	return getStrField(m, conventions.AttributeHostID)
}

func getEC2Hostname(m pcommon.Map) (string, bool, error) {
	if onAWS, err := isAWS(m); err != nil || !onAWS {
		return "", onAWS, err
	}
	return getStrField(m, conventions.AttributeHostName)
}

// Update the information about a given host by providing a resource.
// The function reports:
//   - Whether the information about the `host` has changed
//   - Any non-fatal errors that may have occurred during the update
//
// Partial modifications will still be applied even with non-fatal errors.
func (m *HostMap) Update(host string, res pcommon.Resource) (changed bool, err error) {
	md := payload.HostMetadata{
		Flavor:  "otelcol-contrib",
		Meta:    &payload.Meta{},
		Tags:    &payload.HostTags{},
		Payload: gohai.NewEmpty(),
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	var found bool
	if old, ok := m.hosts[host]; ok {
		found = true
		md = old
	}

	md.InternalHostname = host
	md.Meta.Hostname = host

	// InstanceID field
	if instanceId, ok, fieldErr := getInstanceID(res.Attributes()); fieldErr != nil {
		err = multierr.Append(err, fieldErr)
	} else if ok {
		old := md.Meta.InstanceID
		changed = changed || old != instanceId
		md.Meta.InstanceID = instanceId
	}

	// EC2Hostname field
	if EC2Hostname, ok, fieldErr := getEC2Hostname(res.Attributes()); fieldErr != nil {
		err = multierr.Append(err, fieldErr)
	} else if ok {
		old := md.Meta.EC2Hostname
		changed = changed || old != EC2Hostname
		md.Meta.EC2Hostname = EC2Hostname
	}

	// Gohai - Platform
	md.Gohai.Gohai.Platform.(map[string]string)["hostname"] = host
	for field, attribute := range platformAttributesMap {
		strVal, ok, fieldErr := getStrField(res.Attributes(), attribute)
		if fieldErr != nil {
			err = multierr.Append(err, fieldErr)
		} else if ok {
			old := md.Gohai.Gohai.Platform.(map[string]string)[field]
			changed = changed || old != strVal
			md.Gohai.Gohai.Platform.(map[string]string)[field] = strVal
		}
	}

	m.hosts[host] = md
	changed = changed && found
	return
}

// Flush all the host metadata payloads and clear them from the HostMap.
func (m *HostMap) Flush() map[string]payload.HostMetadata {
	m.mu.Lock()
	defer m.mu.Unlock()
	hosts := m.hosts
	m.hosts = make(map[string]payload.HostMetadata)
	return hosts
}
