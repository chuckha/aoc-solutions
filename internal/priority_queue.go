package internal

type Item[T any] struct {
	Data     T
	Priority int
}

type PriorityQueue[T any] struct {
	Items []Item[T]
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		Items: make([]Item[T], 0),
	}
}

func (p *PriorityQueue[T]) Add(item T, priority int) {
	p.Items = append(p.Items, Item[T]{Data: item, Priority: priority})
}

func (p *PriorityQueue[T]) Empty() bool {
	return len(p.Items) == 0
}

func (p *PriorityQueue[T]) Size() int {
	return len(p.Items)
}

func (p *PriorityQueue[T]) Pull() T {
	minPrio := p.Items[0].Priority
	min := p.Items[0].Data
	minidx := 0
	for i, item := range p.Items {
		if item.Priority < minPrio {
			minPrio = item.Priority
			min = item.Data
			minidx = i
		}
	}
	p.Items = append(p.Items[:minidx], p.Items[minidx+1:]...)
	return min
}
