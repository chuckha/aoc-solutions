package internal

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data: []T{},
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.data = append(q.data, item)
}

func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue[T]) Dequeue() T {
	out := q.data[0]
	q.data = q.data[1:]
	return out
}
