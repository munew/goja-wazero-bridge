package wbridge

import "github.com/dop251/goja"

type externrefStorage struct {
	data  map[uintptr]goja.Value
	index uintptr
}

func newExternrefStorage() *externrefStorage {
	return &externrefStorage{
		data: make(map[uintptr]goja.Value),
		index: 1,
	}
}

func (s *externrefStorage) Get(key uint64) (goja.Value, bool) {
	v, ok := s.data[uintptr(key)]
	return v, ok
}

func (s *externrefStorage) Set(value goja.Value) uintptr {
	v := s.index
	s.data[v] = value
	s.index++
	return v
}
