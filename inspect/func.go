package inspect

import (
	"bytes"
	"fmt"
	"github.com/iamgoroot/dbietool/inspect/gadget"
	"github.com/iamgoroot/dbietool/models"
	"github.com/iamgoroot/dbietool/render"
	"go/ast"
)

func (f *SingleFile) handleIntefaceMethod(handler gadget.TypeHandler[models.Entity, *render.Result], entity models.Entity, method *ast.Field) (result *render.Result) {
	signature := models.Method{
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
	return handler.OnInterfaceMethod(signature)

}
func processInParam(r ast.Expr) (param models.Param) {
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
func processParam(r ast.Expr) (param models.Param) {
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
			return models.Param{
				Name: ident.Name,
				Type: t.Sel.Name,
			}
		}
		return models.Param{
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
					return models.Param{
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
	return models.Param{}
}
