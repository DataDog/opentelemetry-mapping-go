// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import (
	"sync"
)

const (
	defaultBinListSize      = 2 * defaultBinLimit
	defaultKeyListSize      = 256
	defaultOverflowListSize = 16
)

type BinPool[T uint16 | uint32] struct {
	binListPool      sync.Pool
	keyListPool      sync.Pool
	overflowListPool sync.Pool
}

func NewBinPool[T uint16 | uint32]() *BinPool[T] {
	return &BinPool[T]{
		binListPool: sync.Pool{
			New: func() interface{} {
				a := make([]bin[T], 0, defaultBinListSize)
				return &a
			},
		},
		keyListPool: sync.Pool{
			New: func() interface{} {
				a := make([]Key, 0, defaultKeyListSize)
				return &a
			},
		},
		overflowListPool: sync.Pool{
			New: func() interface{} {
				a := make([]bin[T], 0, defaultOverflowListSize)
				return &a
			},
		},
	}
}

func (b *BinPool[T]) getBinList() []bin[T] {
	a := *(b.binListPool.Get().(*[]bin[T]))
	return a[:0]
}

func (b *BinPool[T]) putBinList(a []bin[T]) {
	b.binListPool.Put(&a)
}

func (b *BinPool[T]) getKeyList() []Key {
	a := *(b.keyListPool.Get().(*[]Key))
	return a[:0]
}

func (b *BinPool[T]) putKeyList(a []Key) {
	b.keyListPool.Put(&a)
}

func (b *BinPool[T]) getOverflowList() []bin[T] {
	a := *(b.overflowListPool.Get().(*[]bin[T]))
	return a[:0]
}

func (b *BinPool[T]) putOverflowList(a []bin[T]) {
	b.overflowListPool.Put(&a)
}
