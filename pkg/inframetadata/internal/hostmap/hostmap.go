package hostmap

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/pdata/pcommon"
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
	// u is the updater
	u *updater
}

// New creates a new HostMap.
func New() (*HostMap, error) {
	return &HostMap{
		mu:    sync.Mutex{},
		hosts: make(map[string]payload.HostMetadata),
		u:     &updater{},
	}, nil
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
	m.u.Reset(res)

	var found bool
	if old, ok := m.hosts[host]; ok {
		found = true
		md = old
	}

	m.u.FromValue(&md.InternalHostname, host)

	// Meta section
	m.u.FromFunc(&md.Meta.InstanceID,
		func(m pcommon.Map) (string, bool, error) {
			return "", false, fmt.Errorf("not implemented: InstanceID")
		},
	)
	m.u.FromFunc(&md.Meta.EC2Hostname,
		func(m pcommon.Map) (string, bool, error) {
			return "", false, fmt.Errorf("not implemented: EC2Hostname")
		},
	)
	m.u.FromValue(&md.Meta.Hostname, host)

	// Gohai - Platform
	m.u.FromValue(md.Gohai.Gohai.Platform.(map[string]*string)["hostname"], host)
	for field, attribute := range platformAttributesMap {
		m.u.FromAttributeToMap(md.Gohai.Gohai.Platform.(map[string]string), field, attribute)
	}

	m.hosts[host] = md
	changed = found && m.u.HasChanged()
	err = multierr.Append(err, m.u.Error())
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
