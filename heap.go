package gheap

import (
	"math"
)

// Heapable defines the properties than any item must have
// to be heap-ordered.
type Heapable interface {
	// Order dictates the internal ordering of
	// the items in the heap. Heap is min-ordered.
	Order() int
}

// Heap is a priority queue (min-heap)
type Heap struct {
	storage []Heapable
	maxSize int
}

// NewHeap returns a Heap of the specified size. If size < 0
// heap size is unbounded.
func NewHeap(maxSize int) *Heap {
	if maxSize < 0 {
		return &Heap{maxSize: math.MaxInt64}
	}
	return &Heap{maxSize: maxSize}
}

// Heapify returns a Heap of the specified size using the given
// source slice as its backing storage, and heap-sorts it in < O(n) time.
func Heapify(source []Heapable, maxSize int) *Heap {
	result := &Heap{storage: source, maxSize: maxSize}
	result.heapify()
	return result
}

// Push adds an item to the heap.
// The second return val, if true, indicates that the heap is at its
// maximum capacity the highest priority item was popped and returned
// to you as the first return val
func (h *Heap) Push(val Heapable) (Heapable, bool) {
	h.storage = append(h.storage, val)
	h.percolateUp(len(h.storage) - 1)
	if len(h.storage) > h.maxSize {
		return h.Pop()
	}
	return nil, false
}

// UnsafeStorage yields a shallow copy of the underlying storage of the heap.
// The behaviour following the mutation of the result is undefined
func (h *Heap) UnsafeStorage() []Heapable {
	result := make([]Heapable, 0, len(h.storage))
	copy(result, h.storage)
	return result
}

// Pop removes the highest priority item from the heap.
// The second return val, if false, indicates that the heap is empty
// and that a nil value was returned to you as the first return val
func (h *Heap) Pop() (Heapable, bool) {
	switch len(h.storage) {
	case 0:
		return nil, false
	case 1:
		return h.removeLast(), true
	default:
		result := h.storage[0]
		h.storage[0] = h.removeLast()
		h.percolateDown(0)
		return result, true
	}
}

// Peak returns the highest priority item in the heap without
// dequeuing it.
// The second return val, if false, indicates that the heap is empty
// and that a nil value was returned to you as the first return val
func (h *Heap) Peak() (Heapable, bool) {
	if len(h.storage) > 0 {
		return h.storage[0], true
	}
	return nil, false
}

// Size returns the number of items in the Heap
func (h *Heap) Size() int {
	return len(h.storage)
}

func (h *Heap) removeLast() Heapable {
	result := h.storage[len(h.storage)-1]
	h.storage = h.storage[:len(h.storage)-1]
	return result
}

func (h *Heap) percolateUp(i int) {
	parentIndex := h.parentIndex(i)
	for parentIndex >= 0 && parentIndex < i && !h.inOrder(parentIndex, i) {
		h.storage[parentIndex], h.storage[i] = h.storage[i], h.storage[parentIndex]
		i = parentIndex
		parentIndex = h.parentIndex(i)
	}
}

func (h *Heap) percolateDown(i int) {
	childIndex := h.highestPriorityChildIndex(i)
	for childIndex > -1 && !h.inOrder(i, childIndex) {
		h.storage[i], h.storage[childIndex] = h.storage[childIndex], h.storage[i]
		i = childIndex
		childIndex = h.highestPriorityChildIndex(i)
	}
}

// Returns the highest priority child index.
// If there are no children, returns -1
func (h *Heap) highestPriorityChildIndex(parentIndex int) int {
	left, right := h.leftChildIndex(parentIndex), h.rightChildIndex(parentIndex)
	switch {
	case left >= len(h.storage):
		return -1 // no children
	case right >= len(h.storage):
		return left // no right child
	// both children exist
	case h.storage[left].Order() <= h.storage[right].Order():
		return left // left child greater or equal priority
	default:
		return right // right child greater priority
	}
}

func (h *Heap) inOrder(parentIndex, childIndex int) bool {
	return h.storage[parentIndex].Order() < h.storage[childIndex].Order()
}

func (h *Heap) parentIndex(childIndex int) int {
	return (childIndex - 1) / 2
}

func (h *Heap) leftChildIndex(parentIndex int) int {
	return parentIndex*2 + 1
}

func (h *Heap) rightChildIndex(parentIndex int) int {
	return parentIndex*2 + 2
}

func (h *Heap) heapify() {
	if len(h.storage) == 0 {
		return
	}
	parentIndex := (len(h.storage) - 1) / 2 // skip the bottom row
	for parentIndex >= 0 {
		h.percolateDown(parentIndex)
		parentIndex--
	}
}
