package inspect

import (
	"fmt"
	"github.com/iamgoroot/dbietool/inspect/gadget"
	"github.com/iamgoroot/dbietool/models"
	"github.com/iamgoroot/dbietool/render"
	"go/ast"
	"sort"
)

func (f *SingleFile) handleOneType(result *render.Result, handler gadget.TypeHandler[models.Entity, *render.Result], gen *ast.GenDecl) {
	//var run []func()
	for _, spec := range gen.Specs {
		if ts, ok := spec.(*ast.TypeSpec); ok && contains(f.WithTypes, ts.Name.Name) {
			switch typeItem := ts.Type.(type) {
			case *ast.InterfaceType:
				f.onInteface(handler, ts, typeItem)
				sort.Slice(typeItem.Methods.List, func(i, j int) bool {
					return len(typeItem.Methods.List[i].Names) < len(typeItem.Methods.List[j].Names)
				})
				entity := &models.Entity{Name: ts.Name.Name}
				for _, method := range typeItem.Methods.List {
					if len(method.Names) == 0 {
						res := f.handleEmdeddedInterfacesAndGetModelName(handler, entity, method)
						result.Merge(res)
						continue
					}
					res := f.handleIntefaceMethod(handler, *entity, method)
					result.Merge(res)
				}
			case *ast.IndexExpr:
				fmt.Println("handle alias")
			}
		}
	}
}

func (f *SingleFile) onInteface(handler gadget.TypeHandler[models.Entity, *render.Result], ts *ast.TypeSpec, item *ast.InterfaceType) {
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
	//e := models.Entity{
	//	MethodName: ts.MethodName.String(),
	//	//ModelName: ts.MethodName.String(),
	//}

}
