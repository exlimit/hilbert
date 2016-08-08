// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package hilbert is for mapping values to and from space-filling curves, such as Hilbert and Peano
// curves.
package hilbert

// Hilbert represents a 2D Hilbert space of order N for mapping to and from.
// Implements SpaceFilling interface.
type Hilbert struct {
	N int
}

// New returns a Hilbert space which maps integers to and from the curve.
// n must be a power of two.
func NewHilbert(n int) (*Hilbert, error) {
	if n <= 0 {
		return nil, ErrNotPositive
	}

	// Test if power of two
	if (n & (n - 1)) != 0 {
		return nil, ErrNotPowerOfTwo
	}

	return &Hilbert{
		N: n,
	}, nil
}

// GetDimensions returns the width and height of the 2D space.
func (s *Hilbert) GetDimensions() (int, int) {
	return s.N, s.N
}

// Map transforms a one dimension value, t, in the range [0, n^2-1] to coordinates on the Hilbert
// curve in the two-dimension space, where x and y are within [0,n-1].
func (s *Hilbert) Map(t int) (x, y int, err error) {
	if t < 0 || t >= s.N*s.N {
		return -1, -1, ErrOutOfRange
	}

	for i := 1; i < s.N; i = i * 2 {
		rx := t&2 == 2
		ry := t&1 == 1
		if rx {
			ry = !ry
		}

		x, y = s.rotate(i, x, y, rx, ry)

		if rx {
			x = x + i
		}
		if ry {
			y = y + i
		}

		t /= 4
	}

	return
}

// MapInverse transform coordinates on Hilbert Curve from (x,y) to t.
func (s *Hilbert) MapInverse(x, y int) (t int, err error) {
	if x < 0 || x >= s.N || y < 0 || y >= s.N {
		return -1, ErrOutOfRange
	}

	for i := s.N / 2; i > 0; i = i / 2 {
		rx := (x & i) > 0
		ry := (y & i) > 0

		a := 0
		if rx {
			a = 3
		}
		t += i * i * (a ^ b2i(ry))

		x, y = s.rotate(i, x, y, rx, ry)
	}

	return
}

// rotate rotates and flips the quadrant appropriately.
func (s *Hilbert) rotate(n, x, y int, rx, ry bool) (int, int) {
	if !ry {
		if rx {
			x = n - 1 - x
			y = n - 1 - y
		}

		x, y = y, x
	}
	return x, y
}
