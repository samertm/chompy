package semantic

import "log"

type stable struct {
	// change interface{} to something else?
	scopes []map[string]interface{}
}

func NewStable() *stable {
	s := &stable{scopes: make([]map[string]interface{}, 0, 1)}
	s.scopes = append(s.scopes, make(map[string]interface{}))
	return s
}

func (s *stable) Push() {
	s.scopes = append(scopes, make([]map[string]interface{}))
}

func (s *stable) Pop() {
	if len(s.scopes) == 0 {
		log.Fatal("Popped scope that did not exist")
	}
	s.scopes = s.scopes[:len(s.scopes)-1]
}

// returns true if the spot already exists
func (s *stable) Insert(key string, i interface{}) bool {
	if len(s.scopes) == 0 {
		log.Fatal("Popped scope that did not exist")
		return false
	}
	// insert into top scope
	_, exists := s.scopes[len(s.scopes)-1][key]
	s.scopes[len(s.scopes)-1][key] = i
	return exists
}

// follows item, ok result pattern of maps
func (s *stable) Lookup(key string) (interface{}, bool) {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		if val, ok := s.scopes[i][key]; ok {
			return val, true
		}
	}
	return nil, false
}
