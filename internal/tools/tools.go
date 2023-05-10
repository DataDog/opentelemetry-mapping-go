// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

//go:build tools
// +build tools

package tools

// These imports are used to track test and build tool dependencies.
// This is the currently recommended approach: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/frapposelli/wwhrd"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "go.opentelemetry.io/build-tools/chloggen"
	_ "go.opentelemetry.io/build-tools/multimod"
	_ "golang.org/x/exp/cmd/apidiff"
)
