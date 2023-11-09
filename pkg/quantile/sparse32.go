// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import (
	"math"
	"strings"
	"unsafe"

	"github.com/DataDog/opentelemetry-mapping-go/pkg/quantile/summary"
)

var _ memSized = (*Sketch32)(nil)

// A Sketch for tracking quantiles
// The serialized JSON of Sketch contains the summary only
// Bins are not included.
type Sketch32 struct {
	sparseStore32

	BasicSummary summary.Summary `json:"summary"`
}

func (s Sketch32) String() string {
	var b strings.Builder
	// todo
	// printSketch(&b, s, Default())
	return b.String()
}

func (s Sketch32) BinsLen() int {
	return len(s.bins)
}

func (s Sketch32) BinsCap() int {
	return cap(s.bins)
}

func (s Sketch32) Count() uint64 {
	return s.sparseStore32.count
}

func (s Sketch32) Basic() *summary.Summary {
	return &s.BasicSummary
}

// MemSize returns memory use in bytes:
//
//	used: uses len(bins)
//	allocated: uses cap(bins)
func (s Sketch32) MemSize() (used, allocated int) {
	const (
		basicSize = int(unsafe.Sizeof(summary.Summary{}))
	)

	used, allocated = s.sparseStore32.MemSize()
	used += basicSize
	allocated += basicSize
	return
}

func (s Sketch32) InsertKeys(c *Config, keys []Key) {
	s.sparseStore32.insert(c, keys)
}

// InsertMany values into the sketch.
func (s Sketch32) InsertMany(c *Config, values []float64) {
	keys := getKeyList()

	for _, v := range values {
		s.BasicSummary.Insert(v)
		keys = append(keys, c.key(v))
	}

	s.insert(c, keys)
	putKeyList(keys)
}

// Reset sketch to its empty state.
func (s Sketch32) Reset() {
	s.BasicSummary.Reset()
	s.count = 0
	s.bins = s.bins[:0] // TODO: just release to a size tiered pool.
}

// Insert a single value into the sketch.
// NOTE: InsertMany is much more efficient.
func (s Sketch32) InsertVals(c *Config, vals ...float64) {
	// TODO: remove this
	s.InsertMany(c, vals)
}

// Merge o into s, without mutating o.
func (s *Sketch32) merge(c *Config, o *Sketch32) {
	s.BasicSummary.Merge(o.BasicSummary)
	s.sparseStore32.merge(c, &o.sparseStore32)
}

// Quantile returns v such that s.count*q items are <= v.
//
// Special cases are:
//
//		Quantile(c, q <= 0)  = min
//	 Quantile(c, q >= 1)  = max
func (s Sketch32) Quantile(c *Config, q float64) float64 {
	switch {
	case s.count == 0:
		return 0
	case q <= 0:
		return s.BasicSummary.Min
	case q >= 1:
		return s.BasicSummary.Max
	}

	var (
		n     float64
		rWant = rank(s.count, q)
	)

	for i, b := range s.bins {
		n += float64(b.n)
		if n <= rWant {
			continue
		}

		weight := (n - rWant) / float64(b.n)

		vLow := c.f64(b.k)
		vHigh := vLow * c.gamma.v

		switch i {
		case s.bins.Len():
			vHigh = s.BasicSummary.Max
		case 0:
			vLow = s.BasicSummary.Min
		}

		// TODO|PROD: Interpolate between bucket boundaries, correctly handling min, max,
		// negative numbers.
		// with a gamma of 1.02, interpolating to the center gives us a 1% abs
		// error bound.
		return (vLow*weight + vHigh*(1-weight))
		// return vLow
	}

	// this can happen if count is greater than sum of bins
	return s.BasicSummary.Max
}

// CopyTo makes a deep copy of this sketch into dst.
func (s *Sketch32) CopyTo(dst *Sketch32) {
	// TODO: pool slices here?
	dst.bins = dst.bins.ensureLen(s.bins.Len())
	copy(dst.bins, s.bins)
	dst.count = s.count
	dst.BasicSummary = s.BasicSummary
}

// Copy returns a deep copy
func (s *Sketch32) Copy() *Sketch32 {
	dst := &Sketch32{}
	s.CopyTo(dst)
	return dst
}

func (s Sketch32) CopyAsSketch() Sketch {
	return s.Copy()
}

// Equals returns true if s and o are equivalent.
func (s *Sketch32) Equals(o *Sketch32) bool {
	if s.BasicSummary != o.BasicSummary {
		return false
	}

	if s.count != o.count {
		return false
	}

	if len(s.bins) != len(o.bins) {
		return false
	}

	for i := range s.bins {
		if o.bins[i] != s.bins[i] {
			return false
		}
	}

	return true
}

// ApproxEquals checks if s and o are equivalent, with e error allowed for Sum and Average
func (s *Sketch32) ApproxEquals(o *Sketch32, e float64) bool {
	if math.Abs(s.BasicSummary.Sum-o.BasicSummary.Sum) > e {
		return false
	}

	if math.Abs(s.BasicSummary.Avg-o.BasicSummary.Avg) > e {
		return false
	}

	if s.BasicSummary.Min != o.BasicSummary.Min {
		return false
	}

	if s.BasicSummary.Max != o.BasicSummary.Max {
		return false
	}

	if s.BasicSummary.Cnt != o.BasicSummary.Cnt {
		return false
	}

	if s.count != o.count {
		return false
	}

	if len(s.bins) != len(o.bins) {
		return false
	}

	for i := range s.bins {
		if o.bins[i] != s.bins[i] {
			return false
		}
	}

	return true
}
