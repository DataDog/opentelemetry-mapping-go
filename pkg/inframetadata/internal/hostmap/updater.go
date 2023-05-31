package hostmap

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/multierr"
)

// updater handles updates to a hostmetadata payload,
// recording changes and (non-fatal) errors.
type updater struct {
	// changed records if there have been any changes.
	changed bool
	// attrs stores attributes to use for update of the payload.
	attrs pcommon.Map
	// err records any non-fatal errors.
	err error
}

// newUpdater creates a new Updater.
func newUpdater(res pcommon.Resource) *updater {
	return &updater{
		attrs: res.Attributes(),
	}
}

// FromValue sets field from value and records changes.
func (m *updater) FromValue(field *string, value string) {
	if *field != value {
		m.changed = true
	}

	*field = value
}

// FromFunc sets field from a function and records changes.
// The function must have signature func(pcommon.Map) (value string, ok bool, err error), with return values as follows:'
//   - value is the value that the field should be set to. If different from the existing value, a change will be recorded.
//   - ok records if the value was found. If not found, the update will be silently skipped.
//   - err records any errors found while retrieving the value. If non-nil, the update will be skipped and the error will be logged.
func (m *updater) FromFunc(field *string, fn func(pcommon.Map) (value string, ok bool, err error)) {
	value, ok, err := fn(m.attrs)
	if !ok {
		// Field not available, don't update but don't fail either
		return
	}
	if err != nil {
		m.err = multierr.Append(m.err, err)
		return
	}
	m.FromValue(field, value)
}

// FromAttributeToMap sets map field to value retrieved from a given attribute.
// If the attribute is missing in the attribute map, the update will be silently skipped.
// If the attribute is present but has an incorrect type, the update will be skipped and an error will be logged.
// If the attribute is present, the attribute will be updated. If different from the existing value, a change will be recorded.
func (m *updater) FromAttributeToMap(sourceMap map[string]string, field string, name string) {
	val, ok := m.attrs.Get(name)
	if !ok {
		// Field not available, don't update but don't fail either
		return
	}

	if val.Type() != pcommon.ValueTypeStr {
		m.err = multierr.Append(m.err,
			fmt.Errorf("%q has type %q, expected type \"Str\" instead", name, val.Type()),
		)
		return
	}

	strVal := val.Str()

	if prev, ok := sourceMap[field]; ok && prev != strVal {
		m.changed = true
	}
	sourceMap[field] = name
}

// HasChanged reports if any changes have occcurred during updates.
func (m *updater) HasChanged() bool {
	return m.changed
}

// Error reports any non-fatal errors that have occurred during updates.
func (m *updater) Error() error {
	return m.err
}
