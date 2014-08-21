

import (
        "fmt"
        "github.com/samertm/chompy/lex"
        "github.com/samertm/chompy/parse"
)

var _ = fmt.Print // debugging

func main() {
        tree := parse.Start(tokens)
        fmt.Print(tree)
}
