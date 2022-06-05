package inspect

import (
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

type merger[T any] interface{ Merge(T) T }
type TypeHandler[Result merger[Result]] interface {
	OnEmbeddedInterface(entity inspected.Entity) Result
	OnInterfaceMethod(method inspected.Method) Result
}

type Inspector struct {
	TypeHandler[*inspected.Result]
	ImportLookup map[string]string
}

func (inspector *Inspector) Parse(node ast.Node, render *inspected.Result) {
	render.ImportLookup = map[string]string{}
	file, ok := node.(*ast.File)
	if !ok {
		return
	}
	render.Pkg = file.Name.Name
	for _, dec := range file.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			switch gen.Tok {
			case token.TYPE:
				inspector.handleOneType(render, gen)
			case token.IMPORT:
				inspector.handleImports(render, gen)
				fmt.Println("handle import", gen.Specs)
			}
		}
	}
}

func (inspector *Inspector) handleImports(r *inspected.Result, gen *ast.GenDecl) {
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
