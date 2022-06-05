package inspect

import (
	"bytes"
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"go/ast"
)

func (inspector *Inspector) handleIntefaceMethod(entity inspected.Entity, method *ast.Field) *inspected.Result {
	signature := inspected.Method{
		Entity: entity,
	}
	switch m := method.Type.(type) {

	case *ast.FuncType:
		signature.MethodName = method.Names[0].Name
		for _, r := range m.Params.List {
			signature.In = append(signature.In, processInParam(r.Type))
		}
		for _, r := range m.Results.List {
			signature.Out = append(signature.Out, processParam(r.Type))
		}

		fmt.Println("===>", signature)
	default:
		fmt.Println("type", m)
	}
	return inspector.OnInterfaceMethod(signature)

}
func processInParam(r ast.Expr) (param inspected.Param) {
	switch t := r.(type) {
	case *ast.Ellipsis:
		sliceType := processInParam(t.Elt)
		sliceType.TypePrefix = "..."
		return sliceType
	case *ast.SliceExpr:
		fmt.Println(t.X, t.Low)
	case *ast.SelectorExpr:
		param.Type = t.Sel.Name
		if ident, ok := t.X.(*ast.Ident); ok {
			param.Name = ident.Name
			//f.Imports[param.Type]
		}
	case *ast.StarExpr:
		var param bytes.Buffer //TODO:
		param.WriteString("*")
		if ident, ok := t.X.(*ast.Ident); ok {
			param.WriteString(ident.Name)
		}
		fmt.Println("	star", param.String())
	case *ast.Ident:
		//param.Name = t
		param.Type = t.Name
		//fmt.Println("	func type parameter: [", i, "]", p.Names[0], t.MethodName)
	default:
		fmt.Println("	defaulted Param", param)
	}
	return
}
func processParam(r ast.Expr) (param inspected.Param) {
	switch t := r.(type) {
	case *ast.StarExpr:
	case *ast.Ident:
	//fmt.Println("	return parameter type [", i, "] =>>> ", t.MethodName)
	case *ast.ArrayType:
		param := processParam(t.Elt)
		param.TypePrefix = "[]"
		return param
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return inspected.Param{
				Name: ident.Name,
				Type: t.Sel.Name,
			}
		}
		return inspected.Param{
			Type: t.Sel.Name,
		}
		//fmt.Println("	return parameter type with import [", i, "] ==>>>", t.X, t.Sel.MethodName)
	case *ast.IndexExpr:
		if val, ok := t.X.(*ast.SelectorExpr); ok {
			if ident, ok := val.X.(*ast.Ident); ok {
				//if ident.Name == "dbie" &&  == "Paginated" {
				if ind, ok := t.Index.(*ast.SelectorExpr); ok {
					identInd := ind.X.(*ast.Ident)
					tp := fmt.Sprintf("%s.%s[%s.%s]", ident.String(), val.Sel.Name, identInd.Name, ind.Sel.Name)
					return inspected.Param{
						Name:       val.Sel.Name,
						TypePrefix: "",
						Type:       tp,
					}
				}
				//}
			}
		}
	default:
		fmt.Println("	defaulted return Param", r)
	}
	return inspected.Param{}
}
