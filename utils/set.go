package utils

type IdFieldInterface[T comparable] interface {
	Id() T
}

type Set[T comparable, S IdFieldInterface[T]] struct {
	data []S
}

func NewSet[T comparable, S IdFieldInterface[T]](data []S) *Set[T, S] {
	return &Set[T, S]{data: data}
}

func (s *Set[T, S]) ToMap() map[T]S {
	m := make(map[T]S)
	for _, v := range s.data {
		m[v.Id()] = v
	}
	return m
}

func (s *Set[T, S]) Find(id T) S {
	for _, v := range s.data {
		if v.Id() == id {
			return v
		}
	}
	return *new(S)
}

func (s *Set[T, S]) Length() int {
	return len(s.data)
}

func (s *Set[T, S]) Get(index int) S {
	return s.data[index]
}

func Select[U any, T comparable, S IdFieldInterface[T]](data *Set[T, S], f func(S) U) []U {
	result := make([]U, data.Length())
	for i := 0; i < data.Length(); i++ {
		result[i] = f(data.Get(i))
	}
	return result
}
