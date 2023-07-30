// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cmp provides functions related to comparing ordered values.
package cmp

import "golang.org/x/exp/constraints"

// Min returns the minimal value in x and ys.
// For floating-point numbers, Min propagates NaNs (any NaN value in x or ys
// forces the output to be NaN).
func Min[T constraints.Ordered](x T, ys ...T) T {
	for i := range ys {
		x = min(x, ys[i])
	}
	return x
}

// MinFunc returns the minimal value in x and ys, using cmp to compare elements.
// If there is more than one minimal element
// according to the cmp function, MinFunc returns the first one.
func MinFunc[T any](cmp func(x, y T) int, x T, ys ...T) T {
	for i := range ys {
		if cmp(ys[i], x) < 0 {
			x = ys[i]
		}
	}
	return x
}

// Max returns the maximal value in x and ys.
// For floating-point, Max propagates NaNs (any NaN value in x or ys
// forces the output to be NaN).
func Max[T constraints.Ordered](x T, ys ...T) T {
	for i := range ys {
		x = max(x, ys[i])
	}
	return x
}

// MaxFunc returns the maximal value in x and ys, using cmp to compare elements.
// If there is more than one maximal element
// according to the cmp function, MaxFunc returns the first one.
func MaxFunc[T any](cmp func(x, y T) int, x T, ys ...T) T {
	for i := range ys {
		if cmp(ys[i], x) > 0 {
			x = ys[i]
		}
	}
	return x
}

// min is a version of the predeclared function from the Go 1.21 release.
func min[T constraints.Ordered](x, y T) T {
	if x < y || isNaN(x) {
		return x
	}
	return y
}

// max is a version of the predeclared function from the Go 1.21 release.
func max[T constraints.Ordered](x, y T) T {
	if x > y || isNaN(x) {
		return x
	}
	return y
}

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}
