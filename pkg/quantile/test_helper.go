// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build test
// +build test

package quantile

import (
	"math"
)

func almostEqual(a, b, e float64) bool {
	return math.Abs((a-b)/a) <= e
}

// SketchesApproxEqual checks whether two SketchSeries are equal
func SketchesApproxEqual(exp, act SketchReader, e float64) bool {

	if !almostEqual(exp.Summary().Sum, act.Summary().Sum, e) {
		return false
	}

	if !almostEqual(exp.Summary().Avg, act.Summary().Avg, e) {
		return false
	}

	if !almostEqual(exp.Summary().Max, act.Summary().Max, e) {
		return false
	}

	if !almostEqual(exp.Summary().Min, act.Summary().Min, e) {
		return false
	}

	if exp.Summary().Cnt != exp.Summary().Cnt {
		return false
	}

	if exp.Count() != act.Count() {
		return false
	}

	exp_keys, exp_vals := exp.Cols()
	act_keys, act_vals := act.Cols()

	if len(exp_keys) != len(act_keys) {
		return false
	}

	for i := range exp_keys {
		if math.Abs(float64(act_keys[i]-exp_keys[i])) > 1 {
			return false
		}

		if act_vals[i] != exp_vals[i] {
			return false
		}
	}

	return true
}

// SketchesEqual checks whether two SketchSeries are equal
func SketchesEqual(exp, act SketchReader) bool {
	if exp.Summary().Sum != act.Summary().Sum {
		return false
	}

	if exp.Summary().Avg != act.Summary().Avg {
		return false
	}

	if exp.Summary().Max != act.Summary().Max {
		return false
	}

	if exp.Summary().Min != act.Summary().Min {
		return false
	}

	if exp.Summary().Cnt != exp.Summary().Cnt {
		return false
	}

	if exp.Count() != act.Count() {
		return false
	}

	exp_keys, exp_vals := exp.Cols()
	act_keys, act_vals := act.Cols()

	if len(exp_keys) != len(act_keys) {
		return false
	}

	for i := range exp_keys {
		if act_keys[i] != exp_keys[i] {
			return false
		}

		if act_vals[i] != exp_vals[i] {
			return false
		}
	}

	return true
}

type tHelper interface {
	Helper()
}
