package concurent

import "sync"

type SliceInterface[T any] struct {
	Append func(T)
	Remove func(idx int)
	Get    func(idx int) T
	Size   func() int
	Array  func() []T
}

type Slice[T any] struct {
	slice []T
	mx    sync.RWMutex
}

func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{[]T{}, sync.RWMutex{}}
}

func (s *Slice[T]) Append(element T) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.slice = append(s.slice, element)
}

func (s *Slice[T]) Get(idx int) T {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if len(s.slice) > idx && idx >= 0 {
		return s.slice[idx]
	} else {
		var null T
		return null
	}
}

func (s *Slice[T]) Remove(idx int) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if len(s.slice) > idx && idx >= 0 {
		s.slice = append(s.slice[:idx], s.slice[idx+1:]...)
	}
}

func (s *Slice[T]) Size() int {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return len(s.slice)
}

func (s *Slice[T]) Array() []T {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.slice
}
