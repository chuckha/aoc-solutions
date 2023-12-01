package internal

type Queue[T any] struct {
	data []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data: []T{},
	}
}

func (q *Queue[T]) Enqueue(items ...T) {
	q.data = append(q.data, items...)
}

func (q *Queue[T]) Empty() bool {
	return len(q.data) == 0
}

func (q *Queue[T]) Dequeue() T {
	out := q.data[0]
	q.data = q.data[1:]
	return out
}

func (q *Queue[T]) Internal() []T {
	return q.data
}
