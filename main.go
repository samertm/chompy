package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic"
)

var JFLKDSJKL string
var _ = fmt.Print // debugging

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected filename")
		return
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	source, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	compile(source)
}

func compile(src []byte) {
	_, tokens := lex.Lex("bro", string(src))
	tree, err := parse.Start(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	code, err := semantic.Gen(tree)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(code))
}
