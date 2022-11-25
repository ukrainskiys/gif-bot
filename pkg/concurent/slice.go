package concurent

import "sync"

type SliceInterface[T any] struct {
	Append func(T)
	Get    func(idx int) T
	Size   func() int
}

type Slice[T any] struct {
	slice []T
	mx    *sync.RWMutex
}

func NewSlice[T any](capacity int) *Slice[T] {
	return &Slice[T]{make([]T, capacity), &sync.RWMutex{}}
}

func (s *Slice[T]) Append(element T) {
	defer s.mx.Unlock()
	s.mx.Lock()
	s.slice = append(s.slice, element)
}

func (s *Slice[T]) Get(idx int) T {
	defer s.mx.RUnlock()
	s.mx.RLock()
	if len(s.slice) > idx && idx >= 0 {
		return s.slice[idx]
	} else {
		var null T
		return null
	}
}

func (s *Slice[T]) Size() int {
	defer s.mx.RUnlock()
	s.mx.RLock()
	return len(s.slice)
}
