package gadget

import (
	"fmt"
	"github.com/iamgoroot/dbietool/bun_render"
	"github.com/iamgoroot/dbietool/models"
	"github.com/iamgoroot/dbietool/render"
	"github.com/iamgoroot/dbietool/tr"

	"github.com/iancoleman/strcase"
	"log"
	"strings"
)

type TypeHandler[Model any, Result interface {
	Merge(Result) Result
}] interface {
	OnEmbeddedInterface(Model) Result
	OnInterfaceMethod(method models.Method) Result
}

type GoTypeHandler struct {
}

func (h GoTypeHandler) OnInterfaceMethod(method models.Method) *render.Result {
	name := method.MethodName

	var result tr.RendererResult
	switch {
	case strings.HasPrefix(name, "SelectBy"):
		fieldName := strings.TrimPrefix(name, "SelectBy")
		selectFieldType := ""
		callMethod := "SelectOne"
		pageParamName := ""
		op := "Eq"
		for _, operator := range operators {
			if strings.HasSuffix(name, operator) {
				op = operator
			}
		}
		fieldName = strings.TrimSuffix(fieldName, op)
		for _, param := range method.In {
			switch {
			case strings.ToLower(param.Name) == strings.ToLower(fieldName):
				selectFieldType = fmt.Sprint(param.TypePrefix, param.Type)
			case param.Type == "Page":
				pageParamName = "page"
			default:
				selectFieldType = fmt.Sprint(param.TypePrefix, param.Type)
			}
		}
		for _, param := range method.Out {
			switch {
			case strings.ToLower(param.Type) == strings.ToLower(method.ModelName):
				method.ModelPkg = fmt.Sprint(param.TypePrefix, method.ModelPkg)
				switch param.TypePrefix {
				case "[]":
					callMethod = "Select"
				}
			case strings.HasPrefix(param.Type, "dbie.Paginated"):
				method.ModelName = param.Type
				method.ModelPkg = ""
				callMethod = "SelectPage"
			}
		}
		tableName := strings.TrimPrefix(method.Entity.ModelName, "dbie.Paginated[")
		tableName = strings.TrimSuffix(tableName, "]")

		nameTokens := strings.Split(tableName, ".")
		switch len(nameTokens) {
		case 1:
			tableName = nameTokens[0]
		case 2:
			tableName = nameTokens[1]
		}
		tableName = strcase.ToSnake(tableName)
		result = bun_render.SelectByField.With(models.SelectorMethod{
			Method:          method,
			CallMethod:      callMethod,
			PageParamName:   pageParamName,
			SelectTableName: tableName,
			SelectFieldName: strings.ToLower(fieldName),
			SelectFieldType: strings.ToLower(selectFieldType),
			Op:              op,
		})
	}
	return (&render.Result{
		Renderers: []tr.RendererResult{result},
	}).ImportByName(`"github.com/iamgoroot/dbie"`)
}

var _ TypeHandler[models.Entity, *render.Result] = GoTypeHandler{}

func (h GoTypeHandler) OnEmbeddedInterface(e models.Entity) *render.Result {
	log.Println("handles interface with entity", e)
	result := render.Result{
		Renderers: []tr.RendererResult{
			bun_render.Struct.With(e),
			bun_render.Constr.With(e),
		},
	}
	return (&result).ImportByName(e.ModelPkg)
}
