package semantic

import "github.com/samertm/chompy/parse"

func treeWalks(t *parse.Tree) sErrors {
	walks := []func(*parse.Tree) sErrors {
		checkPackage,
		checkImports,
		checkMain,
	}
	for _, fn := range walks {
		s := fn(t)
		if len(s) != 0 {
			return s
		}
	}
	return nil
}

func checkPackage(t *parse.Tree) sErrors {
	if len(t.Kids) == 0 {
		return nil
	}
	_, ok := t.Kids[0].(*parse.Pkg)
	if !ok {
		return sErrors{"First statement must be package statement"}
	}
	return nil
}

func checkImports(t *parse.Tree) sErrors {
	if len(t.Kids) < 2 {
		return nil
	}
	lastImport := 1
	for i := 1; i < len(t.Kids); i++ {
		_, imptsOk := t.Kids[i].(*parse.Impts)
		_, imptOk := t.Kids[i].(*parse.Impt)
		if imptsOk || imptOk {
			if i > lastImport + 1 {
				return sErrors{"Cannot have imports after other statements"}
			}
			lastImport = i
		}
	}
	return nil
}

func checkMain(t *parse.Tree) sErrors {
	for _, kid := range t.Kids {
		switch f := kid.(type) {
		case *parse.Funcdecl:
			if f.Name.Name == "main" {
				if f.Func.Sig.Params == nil &&
					f.Func.Sig.Result == nil {
					return nil
				}
			}
		}
	}
	return sErrors{"Did not find main"}
}
