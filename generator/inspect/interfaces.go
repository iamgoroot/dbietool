package inspect

import (
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"go/ast"
)

func (inspector *Inspector) handleEmdeddedInterfacesAndGetModelName(
	entity *inspected.Entity,
	method *ast.Field,
) *inspected.Result {
	switch m := method.Type.(type) {
	case *ast.IndexExpr:
		if val, ok := m.X.(*ast.SelectorExpr); ok {
			if val.X.(*ast.Ident).Name == "dbie" && val.Sel.Name == "Repo" {
				if ind, ok := m.Index.(*ast.SelectorExpr); ok {
					ident := ind.X.(*ast.Ident)
					entity.ModelPkg, entity.ModelName, entity.ModelImport = ident.Name, ind.Sel.Name, ident.String()
				}
			}
		}
	case *ast.FuncType:
		fmt.Println("handle func", m)
	default:
		fmt.Println("type", m)
	}
	return inspector.OnEmbeddedInterface(*entity)
}
