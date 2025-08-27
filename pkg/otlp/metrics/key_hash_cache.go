// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"encoding/ascii85"
	"encoding/binary"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/twmb/murmur3"
)

// KeyHashCache is a wrapper around the go-cache library.
// It uses a hash function to compute the key from the string
// to be memory efficient.
// Should be used when the key length is large.
type KeyHashCache struct {
	cache *gocache.Cache
}

func newKeyHashCache(cache *gocache.Cache) KeyHashCache {
	return KeyHashCache{
		cache: cache,
	}
}

type KeyHashCacheKey string

func (m *KeyHashCache) Get(s KeyHashCacheKey) (interface{}, bool) {
	return m.cache.Get(string(s))
}

func (m *KeyHashCache) Set(s KeyHashCacheKey, v interface{}, expiration time.Duration) {
	m.cache.Set(string(s), v, expiration)
}

func (m *KeyHashCache) ComputeKey(s string) KeyHashCacheKey {
	h1, h2 := murmur3.StringSum128(s)
	var bytes [16]byte
	binary.LittleEndian.PutUint64(bytes[0:], h1)
	binary.LittleEndian.PutUint64(bytes[8:], h2)

	var buf [64]byte
	n := ascii85.Encode(buf[:], bytes[:])
	return KeyHashCacheKey(buf[:n])
}
