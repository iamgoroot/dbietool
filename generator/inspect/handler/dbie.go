package handler

import (
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
	buntemplates "github.com/iamgoroot/dbietool/template/bun"
	"github.com/iancoleman/strcase"
	"log"
	"strings"
)

var operators = []string{"Eq", "Neq", "Gt", "Gte", "Lt", "Lte", "Like", "Ilike", "Nlike", "Nilike", "In", "Nin", "Is", "Not"}

type TypeHandler[Model any, Result interface {
	Merge(Result) Result
}] interface {
	OnEmbeddedInterface(Model) Result
	OnInterfaceMethod(method inspected.Method) Result
}

var _ TypeHandler[inspected.Entity, *inspected.Result] = DbieToolHandler{}

type DbieToolHandler struct{}

func (h DbieToolHandler) OnInterfaceMethod(method inspected.Method) *inspected.Result {

	name := method.MethodName
	switch {
	case
		//strings.HasPrefix(name, "UpdateBy"),
		//strings.HasPrefix(name, "DeleteBy"),
		strings.HasPrefix(name, "SelectBy"):
	default:
		return nil
	}

	return generateQueryMethod(method)
}

func generateQueryMethod(method inspected.Method) *inspected.Result {
	methodName := method.MethodName
	sortings := strings.Split(methodName, "OrderBy")
	var orderBy []inspected.OrderBy
	if len(sortings) > 0 {
		methodName = sortings[0]
		sortings = sortings[1:]
	}
	fieldName := strings.TrimPrefix(methodName, "SelectBy")
	selectFieldType := ""
	callMethod := "SelectOne"
	pageParamName := ""

	op := "Eq"
	for _, operator := range operators {
		if strings.HasSuffix(methodName, operator) {
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

	for _, sorting := range sortings {
		order := parseOrderByMethod(tableName, sorting)
		orderBy = append(orderBy, order)
	}

	result := buntemplates.SelectByField.With(inspected.SelectorMethod{
		Method:          method,
		CallMethod:      callMethod,
		PageParamName:   pageParamName,
		SelectTableName: tableName,
		SelectFieldName: strings.ToLower(fieldName),
		SelectFieldType: strings.ToLower(selectFieldType),
		Op:              op,
		OrderBy:         orderBy,
	})
	return (&inspected.Result{
		Renderers: []template.RendererResult{result},
	}).ImportByName(`"context"`).
		ImportByName(`"github.com/uptrace/bun"`).
		ImportByName(`"github.com/iamgoroot/dbie"`)
}

func parseOrderByMethod(tableName, sorting string) inspected.OrderBy {
	sortingUpper := strings.ToUpper(sorting)
	field, order := sorting, "dbie.ASC"
	switch {
	case strings.HasSuffix(sortingUpper, "ASC"): //TODO: implement other sortings. generalize
		field = sorting[:len(sorting)-3]
	case strings.HasSuffix(sortingUpper, "DESC"):
		field = sorting[:len(sorting)-4]
		order = "dbie.DESC"
	}
	return inspected.OrderBy{
		Field: fmt.Sprintf("%s.%s", tableName, strcase.ToSnake(field)),
		Order: order,
	}
}

func (h DbieToolHandler) OnEmbeddedInterface(e inspected.Entity) *inspected.Result {
	log.Println("handles interface with entity", e)
	result := inspected.Result{
		Renderers: []template.RendererResult{
			buntemplates.Struct.With(e),
			buntemplates.Constr.With(e),
		},
	}
	return (&result).ImportByName(e.ModelPkg)
}
