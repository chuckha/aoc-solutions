package internal

// An Item is something we manage in a priority queue.
type HeapItem[T any] struct {
	Value    T   // The value of the item; arbitrary.
	Priority int // The priority of the item in the queue.
}

// A PriorityQueue implements heap.Interface and holds Items.
type FPQ[T any] []*HeapItem[T]

func (pq FPQ[T]) Len() int { return len(pq) }

func (pq FPQ[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq FPQ[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *FPQ[T]) Push(x any) {
	item := x.(*HeapItem[T])
	*pq = append(*pq, item)
}

func (pq *FPQ[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
