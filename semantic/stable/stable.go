package stable

// Let's create a type to hold information about the variables in
// our program.
type NodeInfo struct {
	// What variables do we need? Probably a pointer to a type
	T      *Basic // TODO: change to Type [Issue: https://github.com/samertm/chompy/issues/15]
	Offset int
	// What else? We don't need the identifier name because
	// that's stored in the symbol table. There may be other
	// things but I'm not sure what they are.
	// We need to store the thing above it
	up *NodeInfo
}

// Okay, let's create our Type type. Type will hold all the
// information we need to generate code for a specific type.
type Type interface {
	Equal(Type) bool
	String() string
}

// type Func struct {
// 	// Might not need this for anonymous functions
// 	Name *Basic
// 	Args []Type
// 	Result []Type
// }

// func (f *Func) Equal(t Type) bool {
// 	fn, ok := t.(*Func)
// 	if !ok {
// 		return false
// 	}
// 	// Check the names and results for equality.
// 	// NOTE I'm not sure if I need this, or if this will work
// 	// with closures.
// 	return f.Name.Equal(fn) && typesEqual(f.Result, fn.Result) &&
// 		typesEqual(f.Args, fn.Args)
// }

func typesEqual(types0, types1 []Type) bool {
	if len(types0) != len(types1) {
		return false
	}
	for i := 0; i < len(types0); i++ {
		if types0[i].Equal(types1[i]) == false {
			return false
		}
	}
	return true
}

// func (f *Func) String() string {
// 	s := "func: " + f.Name.String() + "\n"
// 	s += "args: "
// 	for _, a := range f.Args {
// 		s += a.String() + "\n"
// 	}
// 	return s
// }

// Represents all types that are not functions
type Basic struct {
	Pkg  string
	Name string
	Size int
	// this is a pointer type if true
	Pointer bool
}

func (b *Basic) Equal(t Type) bool {
	ba, ok := t.(*Basic)
	if !ok {
		return false
	}
	return b.Pkg == ba.Pkg && b.Name == ba.Name && b.Pointer == ba.Pointer
}

func (b *Basic) String() string {
	s := "pkg: " + b.Pkg + " name: " + b.Name
	if b.Pointer {
		s += " *"
	}
	return s
}

// type Struct struct {
// 	Name   *Basic
// 	Fields []Type
// }

// func (s *Struct) Equal(t Type) bool {
// 	ss, ok := t.(*Struct)
// 	if !ok {
// 		return false
// 	}
// 	if s.Name != s.Name ||
// 		len(s.Fields) != len(ss.Fields) {
// 		return false
// 	}
// 	for i := 0; i < len(s.Fields); i++ {
// 		if s.Fields[i].Equal(ss.Fields[i]) == false {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (s *Struct) String() string {
// 	str := "struct: " + s.Name.String() + "\n"
// 	str += "fields: "
// 	for _, f := range s.Fields {
// 		str += f.String()
// 	}
// 	return str
// }

// TODO: add offset to symbol table [Issue: https://github.com/samertm/chompy/issues/16]
type Stable struct {
	table  map[string]*NodeInfo
	latest *NodeInfo
	up     *Stable
}

// Creates a new Stable. It is legal to pass nil as the old Stable
func New(old *Stable) *Stable {
	return &Stable{
		table: make(map[string]*NodeInfo),
		up:    old,
	}
}

func (s *Stable) Insert(name string, value *NodeInfo) {
	value.up = s.latest
	s.latest = value
	s.table[name] = value
}

func (s *Stable) Get(name string) (*NodeInfo, bool) {
	for tab := s; tab != nil; tab = tab.up {
		n, ok := tab.table[name]
		if ok {
			return n, true
		}
	}
	return nil, false
}

func (s *Stable) IterScope() chan *NodeInfo {
	ch := make(chan *NodeInfo)
	go func() {
		for _, ni := range s.table {
			ch <- ni
		}
	}()
	return ch
}
