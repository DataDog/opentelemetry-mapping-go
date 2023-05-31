package hostmap

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/multierr"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/gohai"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/payload"
)

// HostMap stores a map from hostnames to host metadata payloads.
// Host metadata payloads can be updated with the information from a given resource.
// At any time, the internal map can be extracted and the map will be cleared out.
type HostMap struct {
	// mutex for the host map.
	mutex sync.Mutex
	// hosts map
	hosts map[string]payload.HostMetadata
}

// New creates a new HostMap
func New() (*HostMap, error) {
	return &HostMap{
		mutex: sync.Mutex{},
		hosts: make(map[string]payload.HostMetadata),
	}, nil
}

// Update the information about a given host by providing a resource.
// The function reports:
//   - Whether the information about the `host` has changed
//   - Any non-fatal errors that may have occurred during the update
//
// Partial modifications will still be applied even with non-fatal errors.
func (m *HostMap) Update(host string, res pcommon.Resource) (changed bool, err error) {
	hm := payload.HostMetadata{
		Flavor:  "otelcol-contrib",
		Meta:    &payload.Meta{},
		Tags:    &payload.HostTags{},
		Payload: gohai.NewEmpty(),
	}
	updater := newUpdater(res)
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var existingHost bool
	if oldHM, ok := m.hosts[host]; ok {
		existingHost = true
		hm = oldHM
	}

	updater.FromValue(&hm.InternalHostname, host)

	// Meta section
	updater.FromFunc(&hm.Meta.InstanceID,
		func(m pcommon.Map) (string, bool, error) {
			return "", false, fmt.Errorf("not implemented: InstanceID")
		},
	)
	updater.FromFunc(&hm.Meta.EC2Hostname,
		func(m pcommon.Map) (string, bool, error) {
			return "", false, fmt.Errorf("not implemented: EC2Hostname")
		},
	)
	updater.FromValue(&hm.Meta.Hostname, host)

	// Gohai - Platform
	updater.FromValue(hm.Gohai.Gohai.Platform.(map[string]*string)["hostname"], host)
	for field, attribute := range platformAttributesMap {
		updater.FromAttributeToMap(hm.Gohai.Gohai.Platform.(map[string]string), field, attribute)
	}

	m.hosts[host] = hm
	changed = existingHost && updater.HasChanged()
	err = multierr.Append(err, updater.Error())
	return
}

// Extract all the host metadata payloads and clear them from the HostMap.
func (m *HostMap) Extract() map[string]payload.HostMetadata {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	hosts := m.hosts
	m.hosts = make(map[string]payload.HostMetadata)
	return hosts
}
