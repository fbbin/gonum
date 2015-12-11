// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"testing"
)

func TestDaxpyUnitary(t *testing.T) {
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want []float64
	}{
		{
			alpha: 0,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-3},
		},
		{
			alpha: 1,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-1},
		},
		{
			alpha: 3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{3},
		},
		{
			alpha: -3,
			xData: []float64{2},
			yData: []float64{-3},
			want:  []float64{-9},
		},
		{
			alpha: 0,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha: 1,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 1, 4, -2, 2, -10},
		},
		{
			alpha: 3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, 3, 6, 2, -4, -18},
		},
		{
			alpha: -3,
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  []float64{0, 1, -3, 0, -10, 14, 6},
		},
		{
			alpha: -5,
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  []float64{0, 1, -5, -2, -14, 20, 14, -18},
		},
	} {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		// Test z = alpha * x + y.
		prefix := fmt.Sprintf("test %v (z=a*x+y)", i)
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		z, zFront, zBack := newGuardedVector(test.xData, 1)
		DaxpyUnitary(test.alpha, x, y, z)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !allNaN(zFront) || !allNaN(zBack) {
			t.Errorf(msgGuard, prefix, "z", zFront, zBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, z, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, z)
		}

		// Test y = alpha * x + y.
		prefix = fmt.Sprintf("test %v (y=a*x+y)", i)
		x, xFront, xBack = newGuardedVector(test.xData, 1)
		y, yFront, yBack = newGuardedVector(test.yData, 1)
		DaxpyUnitary(test.alpha, x, y, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, y, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, y)
		}

		// Test x = alpha * x + y.
		prefix = fmt.Sprintf("test %v (x=a*x+y)", i)
		x, xFront, xBack = newGuardedVector(test.xData, 1)
		y, yFront, yBack = newGuardedVector(test.yData, 1)

		DaxpyUnitary(test.alpha, x, y, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDaxpyInc(t *testing.T) {
	for i, test := range []struct {
		alpha float64
		xData []float64
		yData []float64

		want    []float64
		wantRev []float64 // Result when one vector is traversed in reverse direction.
	}{
		{
			alpha:   0,
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    []float64{-3},
			wantRev: []float64{-3},
		},
		{
			alpha:   1,
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    []float64{-1},
			wantRev: []float64{-1},
		},
		{
			alpha:   3,
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    []float64{3},
			wantRev: []float64{3},
		},
		{
			alpha:   -3,
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    []float64{-9},
			wantRev: []float64{-9},
		},
		{
			alpha:   0,
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    []float64{0, 1, 0, 3, -4, 5, -6},
			wantRev: []float64{0, 1, 0, 3, -4, 5, -6},
		},
		{
			alpha:   1,
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    []float64{0, 1, 1, 4, -2, 2, -10},
			wantRev: []float64{-4, -2, 2, 4, -3, 5, -6},
		},
		{
			alpha:   3,
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    []float64{0, 1, 3, 6, 2, -4, -18},
			wantRev: []float64{-12, -8, 6, 6, -1, 5, -6},
		},
		{
			alpha:   -3,
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    []float64{0, 1, -3, 0, -10, 14, 6},
			wantRev: []float64{12, 10, -6, 0, -7, 5, -6},
		},
		{
			alpha:   -5,
			xData:   []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:    []float64{0, 1, -5, -2, -14, 20, 14, -18},
			wantRev: []float64{-25, 21, 15, -7, -9, 0, -6, 7},
		},
	} {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"
		n := len(test.xData)

		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}

				prefix := fmt.Sprintf("test %v, incX = %v, incY = %v", i, incX, incY)
				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)
				DaxpyInc(test.alpha, x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msgGuard, prefix, "y", yFront, yBack)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(y, incY) {
					t.Errorf("%v: modified y argument at non-stride position", prefix)
				}

				want := test.want
				if incX*incY < 0 {
					want = test.wantRev
				}
				if !equalStrided(want, y, incY) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, want, y)
				}
			}
		}
	}
}
