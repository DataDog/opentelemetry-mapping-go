// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import "github.com/DataDog/opentelemetry-mapping-go/pkg/quantile/summary"

type Sketch interface {
	InsertCounts(c *Config, kcs []KeyCount)
	InsertKeys(c *Config, values []Key)
	InsertMany(c *Config, values []float64)
	InsertVals(c *Config, vals ...float64)
	Reset()

	Basic() *summary.Summary
	Cols() (k []int32, n []uint32)

	BinsCap() int
	BinsLen() int
	CopyAsSketch() Sketch
	Count() uint64
	MemSize() (used, allocated int)
	Quantile(c *Config, q float64) float64
	String() string
}
