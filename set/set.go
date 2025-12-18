package set

import (
	"iter"
	"maps"
)

type StringSet struct {
	data map[string]any
}

func Empty() *StringSet {
	return &StringSet{data: make(map[string]any)}
}

func New(items []string) *StringSet {
	s := StringSet{data: make(map[string]any, len(items))}
	for _, v := range items {
		s.data[v] = nil
	}

	return &s
}

func (s *StringSet) Add(key string) {
	s.data[key] = nil
}

func (s *StringSet) Remove(key string) {
	delete(s.data, key)
}

func (s *StringSet) Items() iter.Seq[string] {
	return maps.Keys(s.data)
}

func (s *StringSet) Contains(key string) bool {
	_, ok := s.data[key]
	return ok
}
