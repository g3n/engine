// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import "testing"

// Test simple matrix operations
func Test(t *testing.T) {

	m := NewMatrix()

	// m.Get(1, 1) // panics with "runtime error: index out of range" as expected

	m.Set(2,4, true)
	m.Set(3,2, true)
	if m.Get(1,1) != false {
		t.Error("Get failed")
	}
	if m.Get(2,4) != true {
		t.Error("Get failed")
	}
	if m.Get(3,2) != true {
		t.Error("Get failed")
	}

	m.Set(2,4, false)
	m.Set(3,2, false)
	if m.Get(2,4) != false {
		t.Error("Get failed")
	}
	if m.Get(3,2) != false {
		t.Error("Get failed")
	}

	// m.Get(100, 100) // panics with "runtime error: index out of range" as expected
}
