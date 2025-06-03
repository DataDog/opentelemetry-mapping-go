// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2022-present Datadog, Inc.

package util

import (
	"unicode"
	"unicode/utf8"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

// GetOTelAttrVal returns the matched value as a string in the input map with the given keys.
// If there are multiple keys present, the first matched one is returned.
// If normalize is true, normalize the return value with NormalizeTagValue.
func GetOTelAttrVal(attrs pcommon.Map, normalize bool, keys ...string) string {
	val := ""
	for _, key := range keys {
		attrval, exists := attrs.Get(key)
		if exists {
			val = attrval.AsString()
			break
		}
	}

	if normalize {
		val = NormalizeTagValue(val)
	}

	return val
}

// NormalizeTagValue applies some normalization to ensure the tag value matches the backend requirements.
// It should be used for cases where we have just the tag_value as the input (instead of tag_key:tag_value).
//
//nolint:revive
func NormalizeTagValue(v string) string {
	return normalize(v, false)
}

// GetOTelAttrFromEitherMap returns the matched value as a string in either attribute map with the given keys.
// If there are multiple keys present, the first matched one is returned.
// If the key is present in both maps, map1 takes precedence.
// If normalize is true, normalize the return value with NormalizeTagValue.
func GetOTelAttrFromEitherMap(map1 pcommon.Map, map2 pcommon.Map, normalize bool, keys ...string) string {
	if val := GetOTelAttrVal(map1, normalize, keys...); val != "" {
		return val
	}
	return GetOTelAttrVal(map2, normalize, keys...)
}

var isAlphaLookup = [256]bool{}
var isAlphaNumLookup = [256]bool{}
var isValidASCIIStartCharLookup = [256]bool{}
var isValidASCIITagCharLookup = [256]bool{}

func init() {
	for i := 0; i < 256; i++ {
		isAlphaLookup[i] = isAlpha(byte(i))
		isAlphaNumLookup[i] = isAlphaNum(byte(i))
		isValidASCIIStartCharLookup[i] = isValidASCIIStartChar(byte(i))
		isValidASCIITagCharLookup[i] = isValidASCIITagChar(byte(i))
	}
}

func normalize(v string, removeDigitStartChar bool) string {
	// Fast path: Check if the tag is valid and only contains ASCII characters,
	// if yes return it as-is right away. For most use-cases this reduces CPU usage.
	if isNormalizedASCIITag(v, removeDigitStartChar) {
		return v
	}
	// the algorithm works by creating a set of cuts marking start and end offsets in v
	// that have to be replaced with underscore (_)
	if len(v) == 0 {
		return ""
	}
	var (
		trim  int      // start character (if trimming)
		cuts  [][2]int // sections to discard: (start, end) pairs
		chars int      // number of characters processed
	)
	var (
		i    int  // current byte
		r    rune // current rune
		jump int  // tracks how many bytes the for range advances on its next iteration
	)
	tag := []byte(v)
	for i, r = range v {
		jump = utf8.RuneLen(r) // next i will be i+jump
		if r == utf8.RuneError {
			// On invalid UTF-8, the for range advances only 1 byte (see: https://golang.org/ref/spec#For_range (point 2)).
			// However, utf8.RuneError is equivalent to unicode.ReplacementChar so we should rely on utf8.DecodeRune to tell
			// us whether this is an actual error or just a unicode.ReplacementChar that was present in the string.
			_, width := utf8.DecodeRune(tag[i:])
			jump = width
		}
		// fast path; all letters (and colons) are ok
		switch {
		case r >= 'a' && r <= 'z' || r == ':':
			chars++
			goto end
		case r >= 'A' && r <= 'Z':
			// lower-case
			tag[i] += 'a' - 'A'
			chars++
			goto end
		}
		if unicode.IsUpper(r) {
			// lowercase this character
			if low := unicode.ToLower(r); utf8.RuneLen(r) == utf8.RuneLen(low) {
				// but only if the width of the lowercased character is the same;
				// there are some rare edge-cases where this is not the case, such
				// as \u017F (ſ)
				utf8.EncodeRune(tag[i:], low)
				r = low
			}
		}
		switch {
		case unicode.IsLetter(r):
			chars++
		// If it's not a unicode letter, and it's the first char, and digits are allowed for the start char,
		// we should goto end because the remaining cases are not valid for a start char.
		case removeDigitStartChar && chars == 0:
			trim = i + jump
			goto end
		case unicode.IsDigit(r) || r == '.' || r == '/' || r == '-':
			chars++
		default:
			// illegal character
			chars++
			if n := len(cuts); n > 0 && cuts[n-1][1] >= i {
				// merge intersecting cuts
				cuts[n-1][1] += jump
			} else {
				// start a new cut
				cuts = append(cuts, [2]int{i, i + jump})
			}
		}
	end:
		if i+jump >= 2*maxTagLength {
			// bail early if the tag contains a lot of non-letter/digit characters.
			// If a tag is test🍣🍣[...]🍣, then it's unlikely to be a properly formatted tag
			break
		}
		if chars >= maxTagLength {
			// we've reached the maximum
			break
		}
	}

	tag = tag[trim : i+jump] // trim start and end
	if len(cuts) == 0 {
		// tag was ok, return it as it is
		return string(tag)
	}
	delta := trim // cut offsets delta
	for _, cut := range cuts {
		// start and end of cut, including delta from previous cuts:
		start, end := cut[0]-delta, cut[1]-delta

		if end >= len(tag) {
			// this cut includes the end of the string; discard it
			// completely and finish the loop.
			tag = tag[:start]
			break
		}
		// replace the beginning of the cut with '_'
		tag[start] = '_'
		if end-start == 1 {
			// nothing to discard
			continue
		}
		// discard remaining characters in the cut
		copy(tag[start+1:], tag[end:])

		// shorten the slice
		tag = tag[:len(tag)-(end-start)+1]

		// count the new delta for future cuts
		delta += cut[1] - cut[0] - 1
	}
	return string(tag)
}

// This code is borrowed from dd-go metric normalization

// fast isAlpha for ascii
func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// fast isAlphaNumeric for ascii
func isAlphaNum(b byte) bool {
	return isAlpha(b) || (b >= '0' && b <= '9')
}

func isValidNormalizedMetricName(name string) bool {
	if name == "" {
		return false
	}
	if !isAlphaLookup[name[0]] {
		return false
	}
	for j := 1; j < len(name); j++ {
		b := name[j]
		if !(isAlphaNumLookup[b] || (b == '.' && !(name[j-1] == '_')) || (b == '_' && !(name[j-1] == '_'))) {
			return false
		}
	}
	return true
}

const maxTagLength = 200

func isNormalizedASCIITag(tag string, checkValidStartChar bool) bool {
	if len(tag) == 0 {
		return true
	}
	if len(tag) > maxTagLength {
		return false
	}
	i := 0
	if checkValidStartChar {
		if !isValidASCIIStartCharLookup[tag[0]] {
			return false
		}
		i++
	}
	for ; i < len(tag); i++ {
		b := tag[i]
		// TODO: Attempt to optimize this check using SIMD/vectorization.
		if isValidASCIITagCharLookup[b] {
			continue
		}
		if b == '_' {
			// an underscore is only okay if followed by a valid non-underscore character
			i++
			if i == len(tag) || !isValidASCIITagCharLookup[tag[i]] {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func isValidASCIIStartChar(c byte) bool {
	return ('a' <= c && c <= 'z') || c == ':'
}

func isValidASCIITagChar(c byte) bool {
	return isValidASCIIStartChar(c) || ('0' <= c && c <= '9') || c == '.' || c == '/' || c == '-'
}
