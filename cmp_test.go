// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmp

import (
	"fmt"
	"math"
	"testing"
)

type S struct {
	a int
	b string
}

func (s S) cmp(s2 S) int {
	return s.a - s2.a
}

func TestMinMax(t *testing.T) {
	t.Parallel()

	intCmp := func(x, y int) int { return x - y }

	tests := []struct {
		data    []int
		wantMin int
		wantMax int
	}{
		{[]int{7}, 7, 7},
		{[]int{1, 2}, 1, 2},
		{[]int{2, 1}, 1, 2},
		{[]int{1, 2, 3}, 1, 3},
		{[]int{3, 2, 1}, 1, 3},
		{[]int{2, 1, 3}, 1, 3},
		{[]int{2, 2, 3}, 2, 3},
		{[]int{3, 2, 3}, 2, 3},
		{[]int{0, 2, -9}, -9, 2},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			if gotMin := Min(tt.data[0], tt.data[1:]...); tt.wantMin != gotMin {
				t.Errorf("Min got %v, want %v", gotMin, tt.wantMin)
			}
			if gotMinFunc := MinFunc(intCmp, tt.data[0], tt.data[1:]...); tt.wantMin != gotMinFunc {
				t.Errorf("MinFunc got %v, want %v", gotMinFunc, tt.wantMin)
			}

			if gotMax := Max(tt.data[0], tt.data[1:]...); tt.wantMax != gotMax {
				t.Errorf("Max got %v, want %v", gotMax, tt.wantMax)
			}
			if gotMaxFunc := MaxFunc(intCmp, tt.data[0], tt.data[1:]...); tt.wantMax != gotMaxFunc {
				t.Errorf("MaxFunc got %v, want %v", gotMaxFunc, tt.wantMax)
			}
		})
	}

	svals := []S{
		{1, "a"},
		{2, "a"},
		{1, "b"},
		{2, "b"},
	}

	gotMin := MinFunc(S.cmp, svals[0], svals[1:]...)
	wantMin := S{1, "a"}
	if wantMin != gotMin {
		t.Errorf("MinFunc(%v, %v...) = %v, want %v", svals[0], svals[1:], gotMin, wantMin)
	}

	gotMax := MaxFunc(S.cmp, svals[0], svals[1:]...)
	wantMax := S{2, "a"}
	if wantMax != gotMax {
		t.Errorf("MaxFunc(%v, %v...) = %v, want %v", svals[0], svals[1:], gotMax, wantMax)
	}
}

func TestMinMaxNaNs(t *testing.T) {
	t.Parallel()

	fs := []float64{1.0, 999.9, 3.14, -400.4, -5.14}
	if fmin := Min(fs[0], fs[1:]...); fmin != -400.4 {
		t.Errorf("got min %v, want -400.4", fmin)
	}
	if fmax := Max(fs[0], fs[1:]...); fmax != 999.9 {
		t.Errorf("got max %v, want 999.9", fmax)
	}

	// No matter which element of fs is replaced with a NaN, both Min and Max
	// should propagate the NaN to their output.
	for i := 0; i < len(fs); i++ {
		testfs := append(fs[:0:0], fs...)
		testfs[i] = math.NaN()

		fmin := Min(testfs[0], testfs[1:]...)
		if !math.IsNaN(fmin) {
			t.Errorf("got min %v, want NaN", fmin)
		}

		fmax := Max(testfs[0], testfs[1:]...)
		if !math.IsNaN(fmax) {
			t.Errorf("got max %v, want NaN", fmax)
		}
	}
}
