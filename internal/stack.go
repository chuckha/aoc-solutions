package internal

type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	if s.Empty() {
		panic("pop on empty stack")
	}
	tmp := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return tmp
}

func (s *Stack[T]) Empty() bool {
	return len(s.items) == 0
}

func (s *Stack[T]) PopAll() []T {
	out := s.items
	s.items = make([]T, 0)
	return out
}
