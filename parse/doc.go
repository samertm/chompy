// Structure of the parser:
//
//  - parse.go: contains the parser type and accompanying methods.
//  - nodes.go: contains all of the node types.
//  - toplevelsets.go: top level set is a vocab word I made for what is
//    usually referred to as a "start set". Every identifier for a top
//    level set starts with "top", and the rest of the identifier maps
//    directly to a nonterminal in grammar.txt. Each top level set is
//    either a single token or a list of tokens. These are used with
//    "func (*parser) accept" and "func (*parser) expect". It also
//    contains definitions for important tokens. The defined tokens start
//    with "tok".
//  - grammar.go: contains the recursive-descent part of the parser.
//    Every parsing function represents an element of the grammar from
//    grammar.txt and maps directly to that file, except that the
//    identifier starts with a lowercase letter.
package parse
