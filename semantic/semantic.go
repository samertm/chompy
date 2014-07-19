package semantic

import (
	"fmt"
	"log"

	"github.com/samertm/chompy/parse"
)

// Check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func Check(n parse.Node) string {
	t, ok := n.(*parse.Tree)
	if !ok {
		log.Fatal("Needed a tree.")
	}
	if t.Valid() != true {
		errs := collectErrors(t)
		for _, s := range errs {
			fmt.Println(s)
		}
		log.Fatal("Tree is not valid ):")
	}
	errs := createStables(t)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return "yay"
}

// allChildren will iterate through every single child for a node
// and thrown them down "kids". Does a depth-first search of the
// tree.
// Closes kids when done.
// NOTE on the api, should it be (node, kids) or (kids, node)?
func allChildren(node parse.Node, kids chan<- parse.Node) {
	// Preamble: set up first channel in stack.
	// chans is a stack. New channels are pushed and popped on
	// the right.
	defer close(kids)
	chans := make([]chan parse.Node, 0, 1)
	chans = append(chans, make(chan parse.Node))
	go node.Children(chans[0])
	for len(chans) != 0 {
		next, ok := <-chans[len(chans)-1]
		if !ok {
			// The channel has closed, pop it off the
			// stack.
			chans = chans[:len(chans)-1]
			continue
		}
		// next is a node. Create a new channel to recieve
		// its children and push it into the stack.
		chans = append(chans, make(chan parse.Node))
		go next.Children(chans[len(chans)-1])
		// Pass next to kids.
		kids <- next
	}
}

// Collects all the errors starting at node from Erro nodes.
func collectErrors(node parse.Node) []string {
	// Preamble.
	errors := make([]string, 0, 5)
	nodes := make(chan parse.Node)
	go allChildren(node, nodes)
	// Collect errors.
	for n := range nodes {
		switch n.(type) {
		case *parse.Erro:
			e := n.(*parse.Erro)
			errors = append(errors, e.Desc)
		default:
			if n.Valid() == false {
				fmt.Println("INVALID:", n)
			}
		}
	}
	return errors
}

func createStables(t *parse.Tree) []string {
	// Let's initialize errors so we can report any we see
	errs := make([]string, 0)
	// First, let's create an Stable to hold the information
	// about the root tree's children.
	//rootStable := stable.New(nil)
	// We're going to iterate through t's children and add them
	// to rootStable.
	kids := make(chan parse.Node)
	go t.Children(kids)
	for kid := range kids {
		switch kid.(type) {
		case *parse.Pkg:
			break
		case *parse.Impts:
			break
		case *parse.Funcdecl:
			fmt.Println("FOUND", kid)
		case *parse.Consts:
			fmt.Println("FOUND", kid)
		case *parse.Vars:
			fmt.Println("FOUND", kid)
		default:
			errs = append(errs, kid.String())
		}
	}
	return errs
}
