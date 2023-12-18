// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package hostmap

import (
	"fmt"
	"strings"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
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
func New() *HostMap {
	return &HostMap{
		hosts: make(map[string]payload.HostMetadata),
	}
}

// strField gets a field as string from a resource attribute map.
// It can handle fields of type "Str" and "Int".
// It returns:
// - The field's value, if available
// - Whether the field was present in the map
// - Any errors found in the process
func strField(m pcommon.Map, key string) (string, bool, error) {
	val, ok := m.Get(key)
	if !ok {
		// Field not available, don't update but don't fail either
		return "", false, nil
	}

	var value string
	switch val.Type() {
	case pcommon.ValueTypeStr:
		value = val.Str()
	case pcommon.ValueTypeInt:
		value = val.AsString()
	default:
		return "", false, fmt.Errorf("%q has type %q, expected type \"Str\" instead", key, val.Type())
	}

	return value, true, nil
}

// strSliceField gets a field as a slice from a resource attribute map.
// It can handle fields of type "Slice".
// It returns:
// - The field's value, if available
// - Whether the field was present in the map
// - Any errors found in the process
func strSliceField(m pcommon.Map, key string) ([]string, bool, error) {
	val, ok := m.Get(key)
	if !ok {
		// Field not available, don't update but don't fail either
		return nil, false, nil
	}
	if val.Type() != pcommon.ValueTypeSlice {
		return nil, false, fmt.Errorf("%q has type %q, expected type \"Slice\" instead", key, val.Type())
	}
	if val.Slice().Len() == 0 {
		return nil, false, fmt.Errorf("%q is an empty slice, expected at least one item", key)
	}

	var strSlice []string
	for i := 0; i < val.Slice().Len(); i++ {
		item := val.Slice().At(i)
		if item.Type() != pcommon.ValueTypeStr {
			return nil, false, fmt.Errorf("%s[%d] has type %q, expected type \"Str\" instead", key, i, item.Type())
		}
		strSlice = append(strSlice, item.Str())
	}
	return strSlice, true, nil
}

// isIPv4 checks if a string is an IPv4 address.
// From https://stackoverflow.com/a/48519490
func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

var macReplacer = strings.NewReplacer("-", ":")

// ieeeRAtoGolangFormat converts a MAC address from IEEE RA format to the Go format for MAC addresses.
// The Gohai payload expects MAC addresses in the Go format.
//
// Per the spec: "MAC Addresses MUST be represented in IEEE RA hexadecimal form: as hyphen-separated
// octets in uppercase hexadecimal form from most to least significant."
//
// Golang returns MAC addresses as colon-separated octets in lowercase hexadecimal form from most
// to least significant, so we need to:
// - Replace hyphens with colons
// - Convert to lowercase
//
// This is the inverse of toIEEERA from the resource detection processor system detector.
func ieeeRAtoGolangFormat(IEEERAMACaddress string) string {
	return strings.ToLower(macReplacer.Replace(IEEERAMACaddress))
}

// isAWS checks if a resource attribute map
// is coming from an AWS VM.
func isAWS(m pcommon.Map) (bool, error) {
	cloudProvider, ok, err := strField(m, conventions.AttributeCloudProvider)
	if err != nil {
		return false, err
	} else if !ok {
		// no cloud provider field
		return false, nil
	}
	return cloudProvider == conventions.AttributeCloudProviderAWS, nil
}

// instanceID gets the AWS EC2 instance ID from a resource attribute map.
// It returns:
// - The EC2 instance id if available
// - Whether the instance id was found
// - Any errors found retrieving the ID
func instanceID(m pcommon.Map) (string, bool, error) {
	if onAWS, err := isAWS(m); err != nil || !onAWS {
		return "", onAWS, err
	}
	return strField(m, conventions.AttributeHostID)
}

// ec2Hostname gets the AWS EC2 OS hostname from a resource attribute map.
// It returns:
// - The EC2 OS hostname if available
// - Whether the EC2 OS hostname was found
// - Any errors found retrieving the ID
func ec2Hostname(m pcommon.Map) (string, bool, error) {
	if onAWS, err := isAWS(m); err != nil || !onAWS {
		return "", onAWS, err
	}
	return strField(m, conventions.AttributeHostName)
}

// Set a hardcoded host metadata payload.
func (m *HostMap) Set(md payload.HostMetadata) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hosts[md.Meta.Hostname] = md
	return nil
}

// newOrFetchHostMetadata returns the host metadata payload for a given host or creates a new one.
// This method is NOT thread-safe and should be called while holding the m.mu mutex.
func (m *HostMap) newOrFetchHostMetadata(host string) (payload.HostMetadata, bool) {
	md, ok := m.hosts[host]
	if !ok {
		md = payload.HostMetadata{
			Flavor:  "otelcol-contrib",
			Meta:    &payload.Meta{},
			Tags:    &payload.HostTags{},
			Payload: gohai.NewEmpty(),
		}
	}
	return md, ok
}

// Update the information about a given host by providing a resource.
// The function reports:
//   - Whether the information about the `host` has changed
//   - The host metadata payload stored
//   - Any non-fatal errors that may have occurred during the update
//
// Non-fatal errors are local to the specific field where they happened
// and do not change the other fields. If when filling a field a non-fatal
// error is raised, the error will be reported, the field will be left
// empty and further fields will still be filled.
//
// The order in which resource attributes are read does not affect the final
// host metadata payload, even if non-fatal errors are raised during execution.
func (m *HostMap) Update(host string, res pcommon.Resource) (changed bool, md payload.HostMetadata, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	md, found := m.newOrFetchHostMetadata(host)
	md.InternalHostname = host
	md.Meta.Hostname = host

	// InstanceID field
	if iid, ok, err2 := instanceID(res.Attributes()); err2 != nil {
		err = multierr.Append(err, err2)
	} else if ok {
		old := md.Meta.InstanceID
		changed = changed || old != iid
		md.Meta.InstanceID = iid
	}

	// EC2Hostname field
	if ec2Host, ok, err2 := ec2Hostname(res.Attributes()); err2 != nil {
		err = multierr.Append(err, err2)
	} else if ok {
		old := md.Meta.EC2Hostname
		changed = changed || old != ec2Host
		md.Meta.EC2Hostname = ec2Host
	}

	// Gohai - Platform
	md.Platform()["hostname"] = host
	for field, attribute := range platformAttributesMap {
		strVal, ok, fieldErr := strField(res.Attributes(), attribute)
		if fieldErr != nil {
			err = multierr.Append(err, fieldErr)
		} else if ok {
			old := md.Platform()[field]
			changed = changed || old != strVal
			md.Platform()[field] = strVal
		}
	}

	// Gohai - CPU
	for field, attribute := range cpuAttributesMap {
		strVal, ok, fieldErr := strField(res.Attributes(), attribute)
		if fieldErr != nil {
			err = multierr.Append(err, fieldErr)
		} else if ok {
			old := md.CPU()[field]
			changed = changed || old != strVal
			md.CPU()[field] = strVal
		}
	}

	// Gohai - Network
	if macAddresses, ok, fieldErr := strSliceField(res.Attributes(), attributeHostMAC); fieldErr != nil {
		err = multierr.Append(err, fieldErr)
	} else if ok {
		old := md.Network()[fieldNetworkMACAddress]
		// Take the first MAC addresses for consistency with the Agent's implementation
		// Map from IEEE RA format to the Go format for MAC addresses.
		new := ieeeRAtoGolangFormat(macAddresses[0])
		changed = changed || old != new
		md.Network()[fieldNetworkMACAddress] = new
	}

	if ipAddresses, ok, fieldErr := strSliceField(res.Attributes(), attributeHostIP); fieldErr != nil {
		err = multierr.Append(err, fieldErr)
	} else if ok {
		oldIPv4 := md.Network()[fieldNetworkIPAddressIPv4]
		oldIPv6 := md.Network()[fieldNetworkIPAddressIPv6]

		var foundIPv4 bool
		var foundIPv6 bool
		// Take the first IPv4 and the first IPv6 addresses for consistency with the Agent's implementation
		for _, ip := range ipAddresses {
			if foundIPv4 && foundIPv6 {
				break
			}

			if !foundIPv4 && isIPv4(ip) {
				changed = changed || oldIPv4 != ip
				md.Network()[fieldNetworkIPAddressIPv4] = ip
				foundIPv4 = true
			} else if !foundIPv6 { // not IPv4, so it must be IPv6
				changed = changed || oldIPv6 != ip
				md.Network()[fieldNetworkIPAddressIPv6] = ip
				foundIPv6 = true
			}
		}
	}

	m.hosts[host] = md
	changed = changed && found
	return
}

func (m *HostMap) UpdateFromMetric(host string, metric pmetric.Metric) {
	var point pmetric.NumberDataPoint

	// Take last available point
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		lastIndex := metric.Gauge().DataPoints().Len() - 1
		point = metric.Gauge().DataPoints().At(lastIndex)
	case pmetric.MetricTypeSum:
		lastIndex := metric.Sum().DataPoints().Len() - 1
		point = metric.Sum().DataPoints().At(lastIndex)
	default:
		// unsupported type
		return
	}

	// Take value from point
	var value float64
	switch point.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		value = float64(point.IntValue())
	case pmetric.NumberDataPointValueTypeDouble:
		value = point.DoubleValue()
	default:
		// unsupported type
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	md, _ := m.newOrFetchHostMetadata(host)

	// Gohai - CPU
	data, ok := cpuMetricsMap[metric.Name()]
	if ok {
		if data.ConversionFactor != 0 {
			value = value * data.ConversionFactor
		}
		md.CPU()[data.FieldName] = fmt.Sprintf("%g", value)
	}
}

// Flush all the host metadata payloads and clear them from the HostMap.
func (m *HostMap) Flush() map[string]payload.HostMetadata {
	m.mu.Lock()
	defer m.mu.Unlock()
	hosts := m.hosts
	m.hosts = make(map[string]payload.HostMetadata)
	return hosts
}
