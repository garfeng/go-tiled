/*
Copyright (c) 2017 Lauris Bukšis-Haberkorns <lauris@nix.lv>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package utils contains some generic internal utilities
package utils

import (
	"sort"
)

// SortAnySlice sorts a slice with given less method
func SortAnySlice[T any](data []T, lessMethod func(a, b T) bool) []T {
	s := &sortable[T]{
		data:       data,
		lessMethod: lessMethod,
	}
	sort.Sort(s)
	return s.data
}

type sortable[T any] struct {
	data       []T
	lessMethod func(a, b T) bool
}

// Swap implements from sort.Interface
func (s *sortable[T]) Swap(i, j int) {
	tmp := (s.data)[i]
	(s.data)[i] = (s.data)[j]
	(s.data)[j] = tmp
}

// Less implements from sort.Interface
func (s *sortable[T]) Less(i, j int) bool {
	return s.lessMethod(s.data[i], s.data[j])
}

// Len implements from sort.Interface
func (s *sortable[T]) Len() int {
	return len(s.data)
}
