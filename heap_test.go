package gheap_test

import (
	"gheap"
	"math"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Heap", func() {
	var subject *gheap.Heap

	Describe("empty state inspection", func() {
		BeforeEach(func() {
			subject = gheap.NewHeap(-1)
		})

		Describe("peak", func() {
			It("returns nil", func() {
				_, exists := subject.Peak()
				Expect(exists).To(Equal(false))
			})
		})

		Describe("Size", func() {
			It("returns 0", func() {
				Expect(subject.Size()).To(Equal(0))
			})
		})
	})

	Context("when size is unbounded", func() {

		Describe("Push", func() {
			BeforeEach(func() {
				subject = gheap.NewHeap(-1)
			})

			Context("when the heap is empty", func() {
				item := testItem{1, []byte{}}
				BeforeEach(func() {
					Expect(subject.Size()).To(Equal(0))
					subject.Push(item)
				})

				It("increases in size", func() {
					Expect(subject.Size()).To(Equal(1))
				})

				It("places the new item at the head", func() {
					obj, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(equal(obj, item)).To(Equal(true))
				})
			})

			Context("when the heap has a lower-priority item at the head", func() {
				a, b := testItem{1, []byte{}}, testItem{2, []byte{}}
				BeforeEach(func() {
					subject.Push(a)
					subject.Push(b)
				})

				It("Increases in size", func() {
					Expect(subject.Size()).To(Equal(2))
				})

				It("does not replace the head item", func() {
					item, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(item.Order()).To(Equal(1))
				})
			})

			Context("when the heap has a higher-priority item at the head", func() {
				a, b := testItem{1, []byte{}}, testItem{2, []byte{}}
				BeforeEach(func() {
					subject.Push(b)
					subject.Push(a)
				})

				It("Increases in size", func() {
					Expect(subject.Size()).To(Equal(2))
				})

				It("does not replace the head item", func() {
					item, ok := subject.Peak()
					Expect(ok).To(Equal(true))
					Expect(item.Order()).To(Equal(1))
				})
			})

			Context("when the newest item requires just one swap", func() {
				BeforeEach(func() {
					subject.Push(testItem{4, []byte{}})
					subject.Push(testItem{5, []byte{}})
					subject.Push(testItem{8, []byte{}})
					subject.Push(testItem{6, []byte{}})
					subject.Push(testItem{9, []byte{}})
					subject.Push(testItem{9, []byte{}})
					subject.Push(testItem{7, []byte{}})
				})

				It("does not violate the heap ordering property", func() {
					assertHeapOrdering(subject)
				})
			})
		})

		Describe("Pop", func() {
			BeforeEach(func() {
				subject = gheap.NewHeap(-1)
			})

			Context("when the heap is empty", func() {
				It("returns nil", func() {
					_, exists := subject.Pop()
					Expect(exists).To(Equal(false))
				})
			})

			Context("when the heap has a single item", func() {
				item := testItem{1, []byte{}}
				BeforeEach(func() {
					subject.Push(item)
				})

				It("returns that item", func() {
					obj, ok := subject.Pop()
					Expect(ok).To(Equal(true))
					Expect(equal(obj, item)).To(Equal(true))
				})
			})

			Context("when the heap contains both higher and lower priority items", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 0, val: []byte{}})
					subject.Push(testItem{key: 5, val: []byte{}})
					subject.Push(testItem{key: 1, val: []byte{}})
					subject.Push(testItem{key: 4, val: []byte{}})
					subject.Push(testItem{key: 3, val: []byte{}})
					Expect(subject.Size()).To(Equal(5))
					subject.Push(testItem{key: 2, val: []byte{}}) // should sift to the middle
					Expect(subject.Size()).To(Equal(6))
				})

				It("sorts items by their given order", func() {
					lastVal := math.MinInt64
					for subject.Size() > 0 {
						assertHeapOrdering(subject)
						top, ok := subject.Pop()
						Expect(ok).To(Equal(true))
						Expect(top.Order() > lastVal).To(Equal(true))
						lastVal = top.Order()
					}
				})
			})
		})
	})

	Context("when a size is specified", func() {
		heapSize := 5

		Describe("Push", func() {
			BeforeEach(func() {
				subject = gheap.NewHeap(heapSize)
			})

			Context("when <= size items are inserted", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 1, val: []byte{}})
					subject.Push(testItem{key: 5, val: []byte{}})
					subject.Push(testItem{key: 2, val: []byte{}})
					subject.Push(testItem{key: 4, val: []byte{}})
					subject.Push(testItem{key: 3, val: []byte{}})
				})

				It("allows all items to exist inside", func() {
					Expect(subject.Size()).To(Equal(heapSize))
				})
			})

			Context("when additional items are inserted", func() {
				BeforeEach(func() {
					subject.Push(testItem{key: 0, val: []byte{}})
					subject.Push(testItem{key: 5, val: []byte{}})
					subject.Push(testItem{key: 1, val: []byte{}})
					subject.Push(testItem{key: 4, val: []byte{}})
					subject.Push(testItem{key: 3, val: []byte{}})
					Expect(subject.Size()).To(Equal(heapSize))
				})

				It("does not exceed maximum size", func() {
					subject.Push(testItem{key: 2, val: []byte{}})
					Expect(subject.Size()).To(Equal(heapSize))
				})

				It("retains the lower-priority items", func() {
					subject.Push(testItem{key: 2, val: []byte{}})
					sortedContents := make([]int, 0, 5)
					for subject.Size() > 0 {
						assertHeapOrdering(subject)
						item, ok := subject.Pop()
						Expect(ok).To(Equal(true))
						sortedContents = append(sortedContents, item.Order())
					}
					Expect(sortedContents).To(BeEquivalentTo([]int{1, 2, 3, 4, 5})) // zero is missing
				})

				It("ejects the highest-priority item", func() {
					item, overflowed := subject.Push(testItem{key: 2, val: []byte{}})
					Expect(overflowed).To(Equal(true))
					Expect(item.Order()).To(Equal(0))
				})
			})
		})
	})

	Describe("robustness", func() {

		heapSize := -1 // unbounded
		testSize := 200
		popPercent := 25

		Describe("heap ordering", func() {
			BeforeEach(func() {
				subject = gheap.NewHeap(heapSize)
			})

			It("never violates the heap ordering property", func() {
				for i := 0; i < testSize; i++ {
					if rand.Intn(100) > popPercent {
						item := testItem{key: rand.Int(), val: []byte{}}
						subject.Push(item)
					} else {
						subject.Pop()
					}
					assertHeapOrdering(subject)

				}
			})
		})
	})

	Describe("heapify", func() {
		Context("when the provided slice is empty", func() {
			BeforeEach(func() {
				subject = gheap.Heapify(make([]gheap.Heapable, 0), -1)
			})

			It("generates a valid (albeit empty) heap out of the given slice", func() {
				assertHeapOrdering(subject)
			})
		})

		Context("when the provided heap has items within it", func() {
			BeforeEach(func() {
				nums := []gheap.Heapable{
					testItem{key: 1, val: []byte{}},
					testItem{key: 9, val: []byte{}},
					testItem{key: 2, val: []byte{}},
					testItem{key: 8, val: []byte{}},
					testItem{key: 3, val: []byte{}},
					testItem{key: 7, val: []byte{}},
					testItem{key: 4, val: []byte{}},
					testItem{key: 6, val: []byte{}},
					testItem{key: 5, val: []byte{}},
					testItem{key: 4, val: []byte{}},
					testItem{key: 6, val: []byte{}},
					testItem{key: 3, val: []byte{}},
					testItem{key: 7, val: []byte{}},
					testItem{key: 2, val: []byte{}},
					testItem{key: 8, val: []byte{}},
					testItem{key: 1, val: []byte{}},
					testItem{key: 9, val: []byte{}},
				}
				subject = gheap.Heapify(nums, -1)
			})

			It("generates a valid heap out of the given slice", func() {
				assertHeapOrdering(subject)
			})
		})
	})
})

// testing helpers

type testItem struct {
	key int
	val []byte
}

func (t testItem) Order() int {
	return t.key
}

func equal(a gheap.Heapable, b testItem) bool {
	obj, coerced := a.(testItem)
	if !coerced {
		return false
	}
	return obj.key == b.key
}

func assertHeapOrdering(heap *gheap.Heap) {
	storage := heap.UnsafeStorage()
	for i, item := range storage {
		left, right := i*2+1, i*2+2
		if left < len(storage) {
			Expect(storage[left].Order() >= item.Order()).To(Equal(true))
		}
		if right < len(storage) {
			Expect(storage[right].Order() >= item.Order()).To(Equal(true))
		}
	}
}
