package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/samertm/chompy/lex"
)

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
	lex.Lex("bro", string(src))
}
