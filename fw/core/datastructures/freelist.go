package datastructures

// Based on https://stackoverflow.com/a/48330314
//
// A FreeList data structure which is basically an array that lets you remove
// elements from anywhere in constant-time (leaving holes behind which get
// reclaimed upon subsequent insertions in constant-time).
//

type freeListElement[T any] struct {
	element     T
	nextFreeIdx int
}

// FreeList is a data structure that allows for constant-time removals and
// insertions.
type FreeList[T any] struct {
	data         []freeListElement[T]
	firstFreeIdx int
}

// NewFreeList creates a new FreeList with the given preallocated capacity.
func NewFreeList[T any](prealloc int) *FreeList[T] {
	return &FreeList[T]{
		data:         make([]freeListElement[T], prealloc),
		firstFreeIdx: -1,
	}
}

// Clear removes all elements from the FreeList.
func (fl *FreeList[T]) Clear() {
	fl.data = fl.data[:0]
	fl.firstFreeIdx = -1
}

// Insert adds an element to the FreeList and returns its index.
func (fl *FreeList[T]) Insert(element T) int {
	if fl.firstFreeIdx == -1 {
		// -1 means there are no holes in the array
		elem := freeListElement[T]{element: element}
		fl.data = append(fl.data, elem)
		return len(fl.data) - 1
	} else {
		// reuse a hole in this case
		idx := fl.firstFreeIdx
		fl.firstFreeIdx = fl.data[idx].nextFreeIdx
		fl.data[idx] = freeListElement[T]{element: element}
		return idx
	}
}

// Remove removes the element at the given index from the FreeList.
func (fl *FreeList[T]) Remove(idx int) {
	fl.data[idx].nextFreeIdx = fl.firstFreeIdx
	fl.firstFreeIdx = idx
}

// Get returns the element at the given index.
func (fl *FreeList[T]) Get(idx int) T {
	return fl.data[idx].element
}

// Set sets the element at the given index.
func (fl *FreeList[T]) Set(idx int, element T) {
	fl.data[idx].element = element
}
