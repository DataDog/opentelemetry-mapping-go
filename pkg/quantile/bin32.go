// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import (
	"math"
)

const (
	maxBin32Width = math.MaxUint32
)

type bin32 struct {
	k Key
	n uint32
}

// incrSafe performs `b.n += by` safely handling overflows. When an overflow
// occurs, we set b.n to it's max, and return the leftover amount to increment.
func (b *bin32) incrSafe(by int) int {
	next := by + int(b.n)

	if next > maxBin32Width {
		b.n = maxBin32Width
		return next - maxBin32Width
	}

	b.n = uint32(next)
	return 0
}

// appendSafe appends 1 or more bins with the given key safely handing overflow by
// inserting multiple buckets when needed.
//
//	(1) n <= maxBin32Width :  1 bin
//	(2) n > maxBin32Width  : >1 bin
func appendSafe32(bins []bin32, k Key, n int) []bin32 {
	if n <= maxBin32Width {
		return append(bins, bin32{k: k, n: uint32(n)})
	}

	// on overflow, insert multiple bins with the same key.
	// put full bins at end

	// TODO|PROD: Add validation func that sorts by key and then n (smaller bin first).
	r := uint32(n % maxBin32Width)
	if r != 0 {
		bins = append(bins, bin32{k: k, n: r})
	}

	for i := 0; i < n/maxBin32Width; i++ {
		bins = append(bins, bin32{k: k, n: maxBin32Width})
	}

	return bins
}

type binList32 []bin32

func (bins binList32) nSum() uint64 {
	var s uint64
	s = 0
	for _, b := range bins {
		s += uint64(b.n)
	}
	return s
}

func (bins binList32) Cap() int {
	return cap(bins)
}

func (bins binList32) Len() int {
	return len(bins)
}

func (bins binList32) ensureLen(newLen int) binList32 {
	for cap(bins) < newLen {
		bins = append(bins[:cap(bins)], bin32{})
	}

	return bins[:newLen]
}
