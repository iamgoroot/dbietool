package inspect

import (
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"go/ast"
	"sort"
)

func (inspector *Inspector) handleOneType(result *inspected.Result, gen *ast.GenDecl) {
	for _, spec := range gen.Specs {
		if ts, ok := spec.(*ast.TypeSpec); ok {
			switch typeItem := ts.Type.(type) {
			case *ast.InterfaceType:
				inspector.onInteface(ts, typeItem)
				sort.Slice(typeItem.Methods.List, func(i, j int) bool {
					return len(typeItem.Methods.List[i].Names) < len(typeItem.Methods.List[j].Names)
				})
				entity := &inspected.Entity{Name: ts.Name.Name}
				for _, method := range typeItem.Methods.List {
					if len(method.Names) == 0 {
						res := inspector.handleEmdeddedInterfacesAndGetModelName(entity, method)
						result.Merge(res)
						continue
					}
					res := inspector.handleIntefaceMethod(*entity, method)
					result.Merge(res)
				}
			case *ast.IndexExpr:
				fmt.Println("handle alias")
			}
		}
	}
}

func (inspector *Inspector) onInteface(ts *ast.TypeSpec, item *ast.InterfaceType) {
	fmt.Println(ts.Name, ts.Type)

	//switch intf := item.Interface.(type) {
	//case *ast.FieldList:
	//
	//}
	//for _, field := range ts.TypeParams.List {
	fmt.Println(ts.Name, ts.Type)
	switch ft := ts.Type.(type) {
	case *ast.IndexExpr:
		fmt.Println("*ast.IndexExpr", ft.Index)
	case *ast.Ident:
		fmt.Println("*ast.Ident", ft.Name)
	case *ast.SelectorExpr:
		fmt.Println("*ast.SelectorExpr", ft.X, ft.Sel.Name)
	//case *ast.FieldList:
	//	fmt.Println("*ast.SelectorExpr", ft.X, ft.Sel.MethodName)
	default:
		fmt.Println("defaulted", ft)
	}

	////}
	//e := inspected.Entity{
	//	MethodName: ts.MethodName.String(),
	//	//ModelName: ts.MethodName.String(),
	//}

}
