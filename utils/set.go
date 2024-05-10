package utils

type Set[S any] struct {
	data []S
}

func NewSet[S any](data []S) *Set[S] {
	return &Set[S]{data: data}
}

func SetToMap[T comparable, S any](set *Set[S], generateKey func(S) T) map[T]S {
	m := make(map[T]S)
	for _, v := range set.data {
		key := generateKey(v)
		m[key] = v
	}
	return m
}

func (s *Set[S]) Find(search func(S) bool) S {
	for _, v := range s.data {
		if search(v) {
			return v
		}
	}
	return *new(S)
}

func (s *Set[S]) Length() int {
	return len(s.data)
}

func (s *Set[S]) Get(index int) S {
	return s.data[index]
}

func SetSelect[S any, T any](data *Set[S], f func(S) T) *Set[T] {
	var result []T
	for _, v := range data.data {
		result = append(result, f(v))
	}
	return NewSet(result)
}

func (s *Set[S]) Filter(f func(S) bool) *Set[S] {
	var result []S
	for _, v := range s.data {
		if f(v) {
			result = append(result, v)
		}
	}
	return NewSet(result)
}

func (s *Set[S]) Foreach(f func(int, S)) {
	for i, v := range s.data {
		f(i, v)
	}
}

func (s *Set[S]) ToArray() []S {
	return s.data
}
