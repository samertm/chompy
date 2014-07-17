package semantic

import "log"

// Let's create a type to hold information about the variables in
// our program.
type NodeInfo struct {
	// What variables do we need? Probably a pointer to a type
	T Type
	// What else? We don't need the identifier name because
	// that's stored in the symbol table. There may be other
	// things but I'm not sure what they are.
}

// Okay, let's create our Type type. Type will hold all the
// information we need to generate code for a specific type.
interface Type  {
	Equal(Type) bool
}

type Func struct {
	Args []Type
}

func (f *Func) Equal(t Type) bool {
	fn, ok := t.(*Func)
	if !ok {
		return false
	}
	if len(f.Args) != len(fn.Args) [
		return false
	}
	for i := 0; i < len(f.Args); i++ {
		if f.Args[i].Equal(fn.Args[i]) == false {
			return false
		}
	}
	return true
}

func Basic struct {
	Name string
	// this is a pointer type if true
	Pointer bool
}

func (b *Basic) Equal(t Type) bool {
	ba, ok := t.(*Basic)
	if !ok {
		return false
	}
	return b.Name == ba.Name && b.Pointer == ba.Pointer
}

type stable struct {
	scope []map[string]NodeInfo
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
