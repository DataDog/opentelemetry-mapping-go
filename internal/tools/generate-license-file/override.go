// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CopyrightOverride struct {
	dependencies map[string]string
}

const (
	overrideFilePath = ".copyright-overrides.yml"
)

var globalOverrides *CopyrightOverride

func init() {
	var err error
	if globalOverrides, err = NewOverrideFromFile(overrideFilePath); err != nil {
		panic(fmt.Sprintf("Failed to load overrides: %s", err))
	}
}

func (c *CopyrightOverride) CopyrightNotice(dependency string) ([]string, bool) {
	for pattern, notice := range c.dependencies {
		if ok, err := filepath.Match(pattern, dependency); err != nil {
			// Bad pattern, should never happen
			panic(err)
		} else if ok {
			return []string{notice}, true
		}
	}
	return nil, false
}

func NewOverrideFromFile(path string) (*CopyrightOverride, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %w", path, err)
	}

	overrides := &CopyrightOverride{
		dependencies: map[string]string{},
	}

	if err = yaml.Unmarshal(data, &overrides.dependencies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %q: %w", path, err)
	}

	return overrides, nil
}
