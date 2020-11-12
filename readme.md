# Gheap — a generic heap in Golang

## Initialization
```
heap := gheap.NewHeap(-1) // unbounded size

// add anything you want to the heap
type itemToHeapify struct {
    // any properties you want
}

// so long as it conforms to the orderable interface, which looks like this:
type Orderable interface {
	// Order dictates the internal ordering of the items in the heap. Heap is
	// min-ordered, so lowest-order items have the highest priority
	Order() int
}

// make your itemToHeapify conform to the interface by adding a method named
// `Order`, like this:
func (i itemToHeapify) Order int {
    // Returns the relative order of this item as an int
}
```

## Basic Usage
```
// 1. make sure you're conforming to the protocol
type itemToHeapify struct {
    key int
    val []byte // use whatever you want
}

func (i itemToHeapify) Order int {
    return i.key
}

// 2. create a heap
heap := gheap.NewHeap(-1) // unbounded size

// 3. Push and Pop at will
a, ok := heap.Push(itemToHeapify{ key: 1, val: []byte{} })
b, ok := heap.Pop()
```

- If the heap is at capacity (as specified on initialization) then pushing items will cause the lowest-ordered item to be ejected
- If the heap is empty, then popping will return (nil, false)
