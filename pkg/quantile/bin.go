// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import (
	"strings"
)

func MaxOf[T uint16 | uint32]() T {
	return ^T(0)
}

func maxBinWidth[T uint16 | uint32]() T {
	return MaxOf[T]()
}

type bin[T uint16 | uint32] struct {
	k Key
	n T
}

// incrSafe performs `b.n += by` safely handling overflows. When an overflow
// occurs, we set b.n to it's max, and return the leftover amount to increment.
func (b *bin[T]) incrSafe(by int) int {
	next := by + int(b.n)

	if next > int(maxBinWidth[T]()) {
		b.n = maxBinWidth[T]()
		return next - int(maxBinWidth[T]())
	}

	b.n = T(next)
	return 0
}

// appendSafe appends 1 or more bins with the given key safely handing overflow by
// inserting multiple buckets when needed.
//
//	(1) n <= maxBinWidth :  1 bin
//	(2) n > maxBinWidth  : >1 bin
func appendSafe[T uint16 | uint32](bins []bin[T], k Key, n int) []bin[T] {
	if n <= int(maxBinWidth[T]()) {
		return append(bins, bin[T]{k: k, n: T(n)})
	}

	// on overflow, insert multiple bins with the same key.
	// put full bins at end

	// TODO|PROD: Add validation func that sorts by key and then n (smaller bin first).
	r := T(n % int(maxBinWidth[T]()))
	if r != 0 {
		bins = append(bins, bin[T]{k: k, n: r})
	}

	for i := 0; i < n/int(maxBinWidth[T]()); i++ {
		bins = append(bins, bin[T]{k: k, n: maxBinWidth[T]()})
	}

	return bins
}

type binList[T uint16 | uint32] []bin[T]

func (bins binList[T]) nSum() uint64 {
	s := uint64(0)
	for _, b := range bins {
		s += uint64(b.n)
	}
	return s
}

func (bins binList[T]) Cap() int {
	return cap(bins)
}

func (bins binList[T]) Len() int {
	return len(bins)
}

func (bins binList[T]) ensureLen(newLen int) binList[T] {
	for cap(bins) < newLen {
		bins = append(bins[:cap(bins)], bin[T]{})
	}

	return bins[:newLen]
}

func (bins binList[T]) String() string {
	var w strings.Builder
	printBins(&w, bins, defaultBinPerLine)
	return w.String()
}
