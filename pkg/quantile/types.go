// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import "github.com/DataDog/opentelemetry-mapping-go/pkg/quantile/summary"

type SketchReader interface {
	Summary() *summary.Summary
	Cols() (k []int32, n []uint32)
	Count() int
	CopyInterface() SketchReader
}

type SketchWriter interface {
	Reset()
	InsertKeys(c *Config, keys []Key)
	InsertCounts(c *Config, kcs []KeyCount)
}

type SketchRW interface {
	SketchReader
	SketchWriter
}
