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

type SetV2 map[string]struct{}

func (s SetV2) Remove(key string) SetV2 {
	out := s.Copy()
	delete(out, key)
	return out
}

func (s SetV2) Add(key string) SetV2 {
	out := s.Copy()
	out[key] = struct{}{}
	return out
}

func (s SetV2) Copy() SetV2 {
	out := make(SetV2)
	for k, v := range s {
		out[k] = v
	}
	return out
}

type SetV3 map[string]int

func (s SetV3) Remove(key string) SetV3 {
	out := s.Copy()
	delete(out, key)
	return out
}

func (s SetV3) Add(key string, val int) SetV3 {
	out := s.Copy()
	out[key] = val
	return out
}

func (s SetV3) Copy() SetV3 {
	out := make(SetV3)
	for k, v := range s {
		out[k] = v
	}
	return out
}
