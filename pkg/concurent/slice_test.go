package concurent

import (
	"reflect"
	"sync"
	"testing"
)

const countCases = 1000

func TestNewSlice(t *testing.T) {
	arr := NewSlice[int]()
	if arr == nil {
		t.Error()
	}
	arr.slice = append(arr.slice, 1)
	if reflect.TypeOf(arr.slice[0]).Name() != "int" {
		t.Error()
	}
}

func TestSlice_Append(t *testing.T) {
	arr := makeSlice(countCases)
	if len(arr.slice) != countCases {
		t.Error()
	}
}

func TestSlice_Remove(t *testing.T) {
	arr := makeSlice(countCases)
	const cases = countCases / 2

	wg := sync.WaitGroup{}
	wg.Add(cases)
	for i := 0; i < cases; i++ {
		i := i
		go func() {
			arr.Remove(i)
			wg.Done()
		}()
	}
	wg.Wait()

	if len(arr.slice) != cases {
		t.Error()
	}
}

func TestSlice_Get(t *testing.T) {
	arr := NewSlice[int]()
	for i := 0; i < countCases; i++ {
		arr.Append(i)
	}
	for i := 0; i < countCases; i++ {
		if arr.Get(i) != i {
			t.Fatal()
		}
	}
}

func makeSlice(size int) *Slice[int] {
	arr := NewSlice[int]()
	wg := sync.WaitGroup{}
	wg.Add(size)
	for i := 0; i < size; i++ {
		i := i
		go func() {
			arr.Append(i)
			wg.Done()
		}()
	}
	wg.Wait()
	return arr
}
