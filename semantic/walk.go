package semantic

import (
	"reflect"
	
	"github.com/samertm/chompy/parse"
)

type walkFn func(parse.Node) bool

// Walks through all children
func walkAll(node parse.Node, kids chan<- parse.Node) {
	walkAllHooks(node, kids, nil)
}

// walkAllHooks will iterate through every single child for a node
// and thrown them down "kids", and does dispatch based on the type
// of the node from "hooks". Does a depth-first search of the tree.
// Closes kids when done.
//
// Will walk as normal if kids is
//
// hooks is a map from strings (in the form "*parse.TYPE" to match
// the node types, which are all pointers and mostly from the package
// "parse") to functions. The function returns a bool indicating
// whether we should walk over its children.
//
// As a special case, if the key "all" in hooks is set, it will be
// run for all node types, and the rest of the keys will be ignored.
//
// NOTE on the api, should it be (node, kids) or (kids, node)?
func walkAllHooks(node parse.Node, kids chan<- parse.Node, hooks map[string]walkFn) {
	// Preamble: set up first channel in stack.
	// chans is a stack. New channels are pushed and popped on
	// the right.
	if kids != nil {
		defer close(kids)
	}
	chans := make([]chan parse.Node, 0, 1)
	chans = append(chans, make(chan parse.Node))
	go node.Children(chans[0])
	allFn, allFnSet := hooks["all"]
	for len(chans) != 0 {
		next, ok := <-chans[len(chans)-1]
		if !ok {
			// The channel has closed, pop it off the
			// stack.
			chans = chans[:len(chans)-1]
			continue
		}
		if allFnSet {
			val := allFn(next)
			if !val {
				goto NEXT
			}
		} else if hooks != nil {
			// We need to look inside hooks and dispatch
			// based on the type of next. I assume that
			// reflect is expensive, so we don't take
			// this path if hooks is nil.
			// This needs to be in an else block (more
			// strictly, it needs to be in another block
			// in general) because goto cannot jump over
			// variable declarations.
			typ := reflect.TypeOf(next).String()
			fn, ok := hooks[typ]
			if ok {
				val := fn(next)
				if !val {
					goto NEXT
				}
			}
		}
		// next is a node. Create a new channel to recieve
		// its children and push it into the stack.
		chans = append(chans, make(chan parse.Node))
		go next.Children(chans[len(chans)-1])
	NEXT:
		// Pass next to kids.
		if kids != nil {
			kids <- next
		}
	}
}
