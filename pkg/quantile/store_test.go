// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package quantile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildStore creates a store with the bins defined by a simple dsl:
//
//	<key>:<n> <key>:<n> ...
//
// For example, `0:3 1:1 2:1 2:1 3:max`
// TODO: move to main_test.go
func buildStore[T uint16 | uint32](t *testing.T, dsl string) *sparseStore[T] {
	s := &sparseStore[T]{}

	eachParsedToken[T](t, dsl, bitSize[T](), func(k Key, n uint64) {
		if n > uint64(maxBinWidth[T]()) {
			t.Fatal("n > max", n, maxBinWidth[T]())
		}

		s.count += int(n)
		s.bins = append(s.bins, bin[T]{k: k, n: T(n)})
	})

	return s
}

func TestStore(t *testing.T) {
	t.Run("uint16", func(t *testing.T) {
		testStore[uint16](t)
	})
	t.Run("uint32", func(t *testing.T) {
		testStore[uint32](t)
	})
}

func testStore[T uint16 | uint32](t *testing.T) {
	t.Run("merge", func(t *testing.T) {
		type mt struct {
			s, o, exp string
			binLimit  int
		}

		for _, tt := range []mt{
			{s: "1:1", o: "", exp: "1:1"},
			{s: "", o: "1:1", exp: "1:1"},
			{s: "1:3", o: "1:2", exp: "1:5"},
			{s: "1:max-1", o: "1:max-2", exp: "1:max-3 1:max"},
			{s: "1:1 2:1 3:1", o: "5:1 6:1 10:1", exp: "1:1 2:1 3:1 5:1 6:1 10:1"},

			// binLimit
			{
				s:        "0:1 1:1 2:1 3:1 4:1 5:1 6:1 7:1 8:1 9:1 10:1",
				o:        "0:1 1:1 2:1 3:1 4:1 5:1 6:1 7:1 8:1 9:1",
				exp:      "8:18 9:2 10:1",
				binLimit: 3,
			},
		} {

			t.Run("", func(t *testing.T) {
				var (
					c   = Default()
					s   = buildStore[T](t, tt.s)
					o   = buildStore[T](t, tt.o)
					exp = buildStore[T](t, tt.exp)
				)

				if tt.binLimit != 0 {
					c.binLimit = tt.binLimit
				}

				// TODO|TEST: check that o is not mutated.
				s.merge(c, o)

				if exp.count != s.count {
					t.Errorf("s.count=%d, want %d", s.count, exp.count)
				}

				if nsum := s.bins.nSum(); exp.count != int(nsum) {
					t.Errorf("nSum=%d, want %d", nsum, exp.count)
				}

				require.Equal(t, exp.bins.String(), s.bins.String())
				// don't compare binPool (lazy initialized sync.Pools)
				require.EqualValues(t, exp.count, s.count)
				require.EqualValues(t, exp.bins, s.bins)
			})
		}

	})

	t.Run("trimLeft", func(t *testing.T) {
		for _, tt := range []struct {
			s, e string
			b    int
		}{
			{},
			{s: "1:1", e: "1:1"},
			{s: "1:1", e: "1:1", b: 1},
			{
				// TODO: if the trimmed size is the same as before trimming adds error
				// with no benefit.
				s: "1:max 2:max 3:max",
				e: "2:max 2:max 3:max",
				b: 2,
			},
			{
				s: "1:max 1:max 1:1 2:max 3:1 4:1",
				e: "1:max 1:max 2:1 2:max 3:1 4:1",
				b: 3,
			},
			{
				s: "1:max-1 2:max-1 3:1",
				e: "1:max-1 2:max-1 3:1",
				b: 3,
			},
		} {
			t.Run("", func(t *testing.T) {
				var (
					c   = Default()
					s   = buildStore[T](t, tt.s)
					exp = buildStore[T](t, tt.e)
				)

				if tt.b != 0 {
					c.binLimit = tt.b
				}
				s.bins = s.trimLeft(s.bins, tt.b)

				if exp.count != s.count {
					t.Errorf("s.count=%d, want %d", s.count, exp.count)
				}

				if nsum := s.bins.nSum(); exp.count != int(nsum) {
					t.Errorf("nSum=%d, want %d", nsum, exp.count)
				}

				require.Equal(t, exp.bins.String(), s.bins.String())
				// don't compare binPool (lazy initialized sync.Pools)
				require.EqualValues(t, exp.count, s.count)
				require.EqualValues(t, exp.bins, s.bins)
			})
		}
	})

	t.Run("insert", func(t *testing.T) {
		type insertTest struct {
			s    *sparseStore[T]
			keys []Key
			exp  string
		}

		c := func(startState string, expected string, keys ...Key) insertTest {
			return insertTest{
				s:    buildStore[T](t, startState),
				keys: keys,
				exp:  expected,
			}
		}

		for _, tt := range []insertTest{
			c("",
				"0:3 1:1 2:1 5:1 9:1",
				0, 0, 0, 1, 2, 5, 9,
			),
			c("0:2", "-3:1 -2:1 -1:1 0:2", -1, -2, -3),
			c("0:2", "0:4", 0, 0),
			c("0:max", "0:1 0:max", 0),
			c("0:max 0:max", "0:1 0:max 0:max", 0),
			c("0:1 0:max 0:max", "0:3 0:max 0:max", 0, 0),
			c("1:1 3:1 4:1 5:1 6:1 7:1", "1:1 2:1 3:2 4:1 5:1 6:1 7:1", 2, 3),
			c("1:1 3:1", "1:1 2:3 3:1", 2, 2, 2),
			c("0:max-3", "0:2 0:max", make([]Key, 5)...),
		} {
			// TODO|TEST: that we never exceed binLimit.
			t.Run("", func(t *testing.T) {
				s := tt.s
				s.InsertKeys(Default(), tt.keys)

				exp := buildStore[T](t, tt.exp)
				if exp.count != s.count {
					t.Errorf("s.count=%d, want %d", s.count, exp.count)
				}

				if nsum := s.bins.nSum(); exp.count != int(nsum) {
					t.Errorf("nSum=%d, want %d", nsum, exp.count)
				}

				require.Equal(t, exp.bins.String(), s.bins.String())
				require.Equal(t, exp.bins, s.bins)
			})
		}
	})
}

func TestCols(t *testing.T) {
	t.Run("uint16", func(t *testing.T) {
		testCols[uint16](t)
	})
	t.Run("uint32", func(t *testing.T) {
		testCols[uint32](t)
	})
}

func testCols[T uint16 | uint32](t *testing.T) {
	for _, tt := range []struct {
		store string
		k     []int32
		n     []uint32
	}{
		{
			store: "",
		},
		{
			store: "0:1 1:1 2:2 3:1 4:1 5:1 8:1 9:1 10:max",
			k:     []int32{0, 1, 2, 3, 4, 5, 8, 9, 10},
			n:     []uint32{1, 1, 2, 1, 1, 1, 1, 1, uint32(maxBinWidth[T]())},
		},
		{
			store: "0:1 0:max",
			k:     []int32{0, 0},
			n:     []uint32{1, uint32(maxBinWidth[T]())},
		},
	} {
		st := buildStore[T](t, tt.store)
		k, n := st.Cols()
		assert.Equal(t, k, tt.k, "keys don't match")
		assert.Equal(t, n, tt.n, "values don't match")
	}
}
