package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/samertm/chompy/lex"
	"github.com/samertm/chompy/parse"
	"github.com/samertm/chompy/semantic"
)

var _ = fmt.Print // debugging

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected filename")
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	source, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	compile(source)
}

func compile(src []byte) {
	_, tokens := lex.Lex("bro", string(src))
	tree := parse.Start(tokens)
	fmt.Print(string(semantic.Gen(tree)))
}
