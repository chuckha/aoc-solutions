package internal

type Set[T any] map[string]T

func (s Set[T]) Insert(key string, value T) {
	s[key] = value
}

func (s Set[T]) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set[T]) Intersect(s2 Set[T]) Set[T] {
	inter := make(map[string]T)
	for k, v := range s {
		if s2.Has(k) {
			inter[k] = v
		}
	}
	return inter
}

func (s Set[T]) Union(s2 Set[T]) Set[T] {
	un := make(map[string]T)
	for k, v := range s {
		un[k] = v
	}
	for k, v := range s2 {
		un[k] = v
	}
	return un
}
