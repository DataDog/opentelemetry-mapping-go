// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import "github.com/DataDog/opentelemetry-mapping-go/pkg/quantile/summary"

type Sketch interface {
	Insert(c *Config, values []Key)
	InsertMany(c *Config, values []float64)
	InsertCounts(c *Config, kcs []KeyCount)
	Reset()

	CopyAsSketch() Sketch
	Cols() (k []int32, n []uint32)
	Basic() *summary.Summary

	BinsString() string
	MemSize() (used, allocated int)
	BinsLen() int
	BinsCap() int
	Quantile(c *Config, q float64) float64
}
