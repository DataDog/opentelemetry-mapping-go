// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package hostmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	conventions "go.opentelemetry.io/collector/semconv/v1.18.0"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/internal/testutils"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata/payload"
)

func TestUpdate(t *testing.T) {
	hostInfo := []struct {
		hostname        string
		attributes      map[string]any
		expectedChanged bool
		expectedErr     string
	}{
		{
			hostname: "host-1-hostid",
			attributes: map[string]any{
				conventions.AttributeCloudProvider: conventions.AttributeCloudProviderAWS,
				conventions.AttributeHostID:        "host-1-hostid",
				conventions.AttributeHostName:      "host-1-hostname",
				conventions.AttributeOSDescription: "Fedora Linux",
				conventions.AttributeHostArch:      conventions.AttributeHostArchAMD64,
				attributeKernelName:                "GNU/Linux",
				attributeKernelRelease:             "5.19.0-43-generic",
				attributeKernelVersion:             "#44~22.04.1-Ubuntu SMP PREEMPT_DYNAMIC Mon May 22 13:39:36 UTC 2",
			},
			expectedChanged: false,
		},
		{
			// Same as #1, but missing some attributes
			hostname: "host-1-hostid",
			attributes: map[string]any{
				conventions.AttributeCloudProvider: conventions.AttributeCloudProviderAWS,
				conventions.AttributeHostID:        "host-1-hostid",
				conventions.AttributeHostName:      "host-1-hostname",
				conventions.AttributeOSDescription: "Fedora Linux",
			},
			expectedChanged: false,
		},
		{
			// Same as #1 but wrong type and an update
			hostname: "host-1-hostid",
			attributes: map[string]any{
				conventions.AttributeCloudProvider: conventions.AttributeCloudProviderAWS,
				conventions.AttributeHostID:        "host-1-hostid",
				conventions.AttributeHostName:      "host-1-hostname",
				conventions.AttributeOSDescription: true, // wrong type
				conventions.AttributeHostArch:      conventions.AttributeHostArchAMD64,
				attributeKernelName:                "GNU/Linux",
				attributeKernelRelease:             "5.19.0-43-generic",
				attributeKernelVersion:             "#82~18.04.1-Ubuntu SMP Fri Apr 16 15:10:02 UTC 2021", // changed
			},
			expectedChanged: true,
			expectedErr:     "\"os.description\" has type \"Bool\", expected type \"Str\" instead",
		},
		{
			// Different host, partial information, on Azure
			hostname: "host-2-hostid",
			attributes: map[string]any{
				conventions.AttributeCloudProvider: conventions.AttributeCloudProviderAzure,
				conventions.AttributeHostID:        "host-2-hostid",
				conventions.AttributeHostName:      "host-2-hostname",
				conventions.AttributeHostArch:      conventions.AttributeHostArchARM64,
			},
		},
	}

	hostMap, err := New()
	require.NoError(t, err)

	for _, info := range hostInfo {
		changed, err := hostMap.Update(info.hostname, testutils.NewResourceFromMap(t, info.attributes))
		assert.Equal(t, info.expectedChanged, changed)
		if info.expectedErr != "" {
			assert.EqualError(t, err, info.expectedErr)
		} else {
			assert.NoError(t, err)
		}
	}

	hosts := hostMap.Flush()
	assert.Len(t, hosts, 2)

	if assert.Contains(t, hosts, "host-1-hostid") {
		md := hosts["host-1-hostid"]
		assert.Equal(t, md.InternalHostname, "host-1-hostid")
		assert.Equal(t, md.Flavor, "otelcol-contrib")
		assert.Equal(t, md.Meta, &payload.Meta{
			InstanceID:  "host-1-hostid",
			EC2Hostname: "host-1-hostname",
			Hostname:    "host-1-hostid",
		})
		assert.Equal(t, md.Tags, &payload.HostTags{})
		assert.Equal(t, md.Payload.Gohai.Gohai.Platform, map[string]string{
			"hostname":                    "host-1-hostid",
			fieldPlatformOS:               "Fedora Linux",
			fieldPlatformProcessor:        "amd64",
			fieldPlatformMachine:          "amd64",
			fieldPlatformHardwarePlatform: "amd64",
			fieldPlatformKernelName:       "GNU/Linux",
			fieldPlatformKernelRelease:    "5.19.0-43-generic",
			fieldPlatformKernelVersion:    "#82~18.04.1-Ubuntu SMP Fri Apr 16 15:10:02 UTC 2021",
		})
		assert.Nil(t, md.Payload.Gohai.Gohai.CPU)
		assert.Nil(t, md.Payload.Gohai.Gohai.FileSystem)
		assert.Nil(t, md.Payload.Gohai.Gohai.Memory)
		assert.Nil(t, md.Payload.Gohai.Gohai.Network)
	}

	if assert.Contains(t, hosts, "host-2-hostid") {
		md := hosts["host-2-hostid"]
		assert.Equal(t, md.InternalHostname, "host-2-hostid")
		assert.Equal(t, md.Flavor, "otelcol-contrib")
		assert.Equal(t, md.Meta, &payload.Meta{
			Hostname: "host-2-hostid",
		})
		assert.Equal(t, md.Tags, &payload.HostTags{})
		assert.Equal(t, md.Payload.Gohai.Gohai.Platform, map[string]string{
			"hostname":                    "host-2-hostid",
			fieldPlatformProcessor:        "arm64",
			fieldPlatformMachine:          "arm64",
			fieldPlatformHardwarePlatform: "arm64",
		})
		assert.Nil(t, md.Payload.Gohai.Gohai.CPU)
		assert.Nil(t, md.Payload.Gohai.Gohai.FileSystem)
		assert.Nil(t, md.Payload.Gohai.Gohai.Memory)
		assert.Nil(t, md.Payload.Gohai.Gohai.Network)
	}

	assert.Empty(t, hostMap.Flush(), "returned map must be empty after double flush")
}
