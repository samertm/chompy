package parse

import (
	"github.com/samertm/chompy/lex"

	"fmt"
)

type Node interface {
	Eval() string
}

type grammarFn func(*parser) Node

type tree struct {
	kids []Node
}

func (t *tree) Eval() (s string) {
	for _, k := range t.kids {
		s += k.Eval()
	}
	return
}

type pkg struct {
	name string
}

func (p *pkg) Eval() string {
	return fmt.Sprintln("in package ", p.name)
}

type impts struct {
	imports []Node
}

func (i *impts) Eval() (s string) {
	s += fmt.Sprintln("start imports")
	for _, im := range i.imports {
		s += im.Eval()
	}
	s += fmt.Sprintln("end imports")
	return
}

type impt struct {
	pkgName  string
	imptName string
}

func (i *impt) Eval() string {
	return fmt.Sprintln("import: pkgName: " + i.pkgName + " imptName: " + i.imptName)
}

type erro struct {
	desc string
}

func (e *erro) Eval() string {
	return fmt.Sprintln("error: ", e.desc)
}

func Start(toks chan lex.Token) Node {
	p := &parser{
		toks:    toks,
		oldToks: make([]*lex.Token, 0),
		nodes:   make(chan Node),
	}
	t := sourceFile(p)
	return t
}

// should the states return their list?... probably but not rn
// every nonterminal function assumes that it is in the correct starting state,
// except for sourceFile
func sourceFile(p *parser) *tree {
	tr := &tree{kids: make([]Node, 0)}
	if !p.accept(topPackageClause) {
		tr.kids = append(tr.kids, &erro{"PackageClause not found"})
		return tr
	}
	pkg := packageClause(p)
	tr.kids = append(tr.kids, pkg)
	if err := p.expect(tokSemicolon); err != nil {
		tr.kids = append(tr.kids, err)
		return tr
	}
	p.next()
	for p.accept(topImportDecl) {
		impts := importDecl(p)
		tr.kids = append(tr.kids, impts)
		if err := p.expect(tokSemicolon); err != nil {
			tr.kids = append(tr.kids, err)
		}
		p.next()
	}
	for p.accept(topTopLevelDecl...) {
		// things
	}
	close(p.nodes)
	return tr
}

func packageClause(p *parser) Node {
	p.next() // eat "package"
	if err := p.expect(topPackageName); err != nil {
		return err
	}
	return packageName(p)
}

func packageName(p *parser) Node {
	t := p.next()
	// should I sanity-check t?
	return &pkg{name: t.Val}
}

func importDecl(p *parser) Node {
	p.next() // eat "import"
	i := &impts{imports: make([]Node, 0)}
	if p.accept(tokOpenParen) {
		p.next() // eat "("
		for p.accept(topImportSpec...) {
			i.imports = append(i.imports, importSpec(p))
			if err := p.expect(tokSemicolon); err != nil {
				return err
			}
			p.next() // eat ";"
		}
		if err := p.expect(tokCloseParen); err != nil {
			return err
		}
		p.next() // eat ")"
		return i
	}
	// a single importSpec
	if !p.accept(topImportSpec...) {
		return &erro{"expected importSpec"}
	}
	i.imports = append(i.imports, importSpec(p))
	return i
}

func importSpec(p *parser) Node {
	i := &impt{}
	if p.accept(tokDot) {
		p.next() // eat dot
		i.pkgName = "."
	}
	if p.accept(topPackageName) {
		t := p.next() // t is the package name
		if i.pkgName == "." {
			// a dot was already processed
			return &erro{"expected tokString"}
		}
		i.pkgName = t.Val
	}
	if !p.accept(topImportPath) {
		return &erro{"expected tokString"}
	}
	// process importPath here.
	t := p.next()
	i.imptName = t.Val
	return i
}
