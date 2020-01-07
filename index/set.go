package index

type Set struct {
	data map[string]bool
}

func (s *Set) Add(entry string) bool {
	if _, exists := s.data[entry]; exists {
		return false
	}
	s.data[entry] = true
	return true
}

func (s *Set) Exists(entry string) bool {
	_, exists := s.data[entry]
	return exists
}

func (s *Set) Remove(entry string) bool {
	if _, exists := s.data[entry]; !exists {
		return false
	}
	delete(s.data, entry)
	return true
}

func (s *Set) Entries() (keys []string) {
	for k, _ := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func NewSet() *Set {
	return &Set{
		data: map[string]bool{},
	}
}
