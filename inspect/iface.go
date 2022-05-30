package inspect

import (
	"fmt"
	"github.com/iamgoroot/dbietool/inspect/gadget"
	"github.com/iamgoroot/dbietool/parse"
	"github.com/iamgoroot/dbietool/render"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

type SingleFile struct {
	//Result
	parse.File
	WithTypes []string
}

func (f *SingleFile) Parse(node ast.Node, render *render.Result) {
	render.ImportLookup = map[string]string{}
	file, ok := node.(*ast.File)
	if !ok {
		return
	}
	render.Pkg = file.Name.Name
	for _, dec := range file.Decls {
		currentTypeHandler := gadget.GoTypeHandler{}
		if gen, ok := dec.(*ast.GenDecl); ok {
			switch gen.Tok {
			case token.TYPE:
				f.handleOneType(render, currentTypeHandler, gen)
			case token.IMPORT:
				f.handleImports(render, gen)
				fmt.Println("handle import", gen.Specs)
			}
		}
	}
}

func (f *SingleFile) handleImports(r *render.Result, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		switch sp := spec.(type) {
		case *ast.ImportSpec:
			name := ""
			if sp.Name != nil {
				name = sp.Name.Name
			}
			if name == "" {
				name = filepath.Base(strings.Trim(sp.Path.Value, `"`))
			}
			if sp.Path.Value != "" {
				r.ImportLookup[name] = sp.Path.Value
			}
		default:
			fmt.Println("import defaulted", sp)
		}
	}
}
