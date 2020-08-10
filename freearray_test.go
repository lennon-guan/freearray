package freearray

import "testing"

func assertInt32(t *testing.T, result int32, expected int32, caseName string) {
	if result != expected {
		t.Errorf("--> case [%s], expected %d, got %d", caseName, expected, result)
	}
}

func TestAlloc(t *testing.T) {
	l := New(4)
	l.printState()
	assertInt32(t, l.Alloc(1), 0, "alloc 0")
	l.printState()
	assertInt32(t, l.Alloc(1), 1, "alloc 1")
	l.printState()
	assertInt32(t, l.Alloc(1), 2, "alloc 2")
	l.printState()
	assertInt32(t, l.Alloc(1), 3, "alloc 3")
	l.printState()
	assertInt32(t, l.Alloc(1), -1, "alloc 4")
	l.printState()
	l.Release(2)
	l.printState()
	l.Release(0)
	l.printState()
	l.Release(1)
	l.printState()
	l.Release(1)
	l.printState()
	l.Alloc(1)
	l.printState()
}
