package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func FirstNElementsOfSlice[Element any](slice []Element, n int) []Element {
	if len(slice) < n {
		return slice
	}
	return slice[:n]
}

func AddToQueryIfSet(q url.Values, key string, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func MapToCommaString[t any](slice []t) string {
	var result []string
	for _, v := range slice {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return strings.Join(result, ",")
}

type Set[t comparable] struct {
	m map[t]struct{}
}

func (s *Set[t]) Add(item t) {
	if s.m == nil {
		s.m = make(map[t]struct{})
	}
	s.m[item] = struct{}{}
}

func (s *Set[t]) Remove(item t) {
	delete(s.m, item)
}

func (s *Set[t]) Has(item t) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[t]) Len() int {
	return len(s.m)
}

func (s *Set[t]) ToSlice() []t {
	var items []t
	for k := range s.m {
		items = append(items, k)
	}
	return items
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := new(Set[T])
	for _, v := range items {
		s.Add(v)
	}
	return s
}
