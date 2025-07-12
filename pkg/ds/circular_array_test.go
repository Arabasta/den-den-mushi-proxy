package ds

import (
	"sync"
	"testing"
)

func TestCircularArrayBasic(t *testing.T) {
	capacity := 3
	c := NewCircularArray[int](capacity)

	c.Add(1)
	c.Add(2)
	c.Add(3)

	got := c.GetAll()
	want := []int{1, 2, 3}

	if !equalSlices(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestCircularArrayOverwrite(t *testing.T) {
	capacity := 3
	c := NewCircularArray[int](capacity)

	// Add more than capacity
	c.Add(1)
	c.Add(2)
	c.Add(3)
	c.Add(4) // overwrites 1

	got := c.GetAll()
	want := []int{2, 3, 4}

	if !equalSlices(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestCircularArrayWrapAround(t *testing.T) {
	capacity := 3
	c := NewCircularArray[int](capacity)

	for i := 1; i <= 10; i++ {
		c.Add(i)
	}

	got := c.GetAll()
	want := []int{8, 9, 10}

	if !equalSlices(got, want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestCircularArrayLen(t *testing.T) {
	capacity := 5
	c := NewCircularArray[int](capacity)

	if c.Len() != 0 {
		t.Errorf("Expected length 0, got %d", c.Len())
	}

	c.Add(42)
	if c.Len() != 1 {
		t.Errorf("Expected length 1, got %d", c.Len())
	}

	for i := 0; i < 10; i++ {
		c.Add(i)
	}

	if c.Len() != capacity {
		t.Errorf("Expected length %d, got %d", capacity, c.Len())
	}
}

func TestCircularArrayConcurrent(t *testing.T) {
	capacity := 50
	c := NewCircularArray[int](capacity)

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			c.Add(val)
		}(i)
	}

	wg.Wait()

	if c.Len() != capacity {
		t.Errorf("Expected length %d after concurrent writes, got %d", capacity, c.Len())
	}
}

func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
