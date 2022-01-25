package internal

// An Item is something we manage in a priority queue.
type FastItem[T any] struct {
	Value    T   // The value of the item; arbitrary.
	Priority int // The priority of the item in the queue.
	Found    int
}

// A PriorityQueue implements heap.Interface and holds Items.
type FastPQ[T any] []*FastItem[T]

func (pq FastPQ[T]) Len() int { return len(pq) }

func (pq FastPQ[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	if pq[i].Priority < pq[j].Priority {
		return true
	}
	if pq[i].Priority > pq[j].Priority {
		return false
	}
	return pq[i].Found < pq[j].Found
}

func (pq FastPQ[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *FastPQ[T]) Push(x interface{}) {
	item := x.(*FastItem[T])
	*pq = append(*pq, item)
}

func (pq *FastPQ[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
