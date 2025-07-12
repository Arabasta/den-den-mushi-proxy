package ds

import (
	"sync"
)

type CircularArray[T any] struct {
	data []T

	start    int
	size     int
	capacity int

	mu sync.RWMutex
}

func NewCircularArray[T any](capacity int) *CircularArray[T] {
	return &CircularArray[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

func (c *CircularArray[T]) Add(item T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.size >= c.capacity {
		// overwrite oldest
		c.data[c.start] = item
		c.start = (c.start + 1) % c.capacity
	} else {
		c.data[(c.start+c.size)%c.capacity] = item
		c.size++
	}
}

// GetAll returns all items in order from oldest to newest
func (c *CircularArray[T]) GetAll() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]T, c.size)
	for i := 0; i < c.size; i++ {
		out[i] = c.data[(c.start+i)%c.capacity]
	}

	return out
}

func (c *CircularArray[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.size
}
