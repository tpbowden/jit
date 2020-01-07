package index

import "sort"

type SortedSet struct {
	entries *Set
	order   []string
}

func (s *SortedSet) Add(entry string) bool {
	if s.entries.Add(entry) {
		index := sort.SearchStrings(s.order, entry)
		s.order = append(s.order, "")
		copy(s.order[index+1:], s.order[index:])
		s.order[index] = entry
		return true
	}
	return false
}

func (s *SortedSet) Remove(entry string) bool {
	if s.entries.Remove(entry) {
		index := sort.SearchStrings(s.order, entry)
		copy(s.order[index:], s.order[index+1:])
		s.order[len(s.order)-1] = ""
		s.order = s.order[:len(s.order)-1]
		return true
	}
	return false
}

func (s *SortedSet) Entries() []string {
	return s.order
}

func NewSortedSet() *SortedSet {
	return &SortedSet{
		entries: NewSet(),
		order:   []string{},
	}
}
