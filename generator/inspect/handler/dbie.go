package handler

import (
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
	buntemplates "github.com/iamgoroot/dbietool/template/core"
	"github.com/iancoleman/strcase"
	"log"
	"strings"
)

var operators = []string{"Eq", "Neq", "Gt", "Gte", "Lt", "Lte", "Like", "Ilike", "Nlike", "Nilike", "In", "Nin", "Is", "Not"}

type TypeHandler[Model any, Result interface {
	Merge(Result) Result
}] interface {
	Once() template.RendererResult
	OnEmbeddedInterface(Model) Result
	OnInterfaceMethod(method inspected.Method) Result
}

var _ TypeHandler[inspected.Entity, *inspected.Result] = DbieToolHandler{}

type DbieToolHandler struct {
	Opts map[string]interface{}
}

func (h DbieToolHandler) Once() template.RendererResult {
	return buntemplates.Factory.With(h.Opts)
}

func (h DbieToolHandler) OnInterfaceMethod(method inspected.Method) *inspected.Result {
	name := method.MethodName
	switch {
	case
		// strings.HasPrefix(name, "UpdateBy"), TODO: implement
		// strings.HasPrefix(name, "DeleteBy"),
		strings.HasPrefix(name, "SelectBy"), strings.HasPrefix(name, "FindBy"):
		return h.generateQueryMethod(method)
	}

	return nil
}

func (h DbieToolHandler) generateQueryMethod(method inspected.Method) *inspected.Result {
	methodName := method.MethodName
	method.Opts = h.Opts
	sortings := strings.Split(methodName, "OrderBy")
	var orderBy []inspected.OrderBy
	if len(sortings) > 0 {
		methodName = sortings[0]
		sortings = sortings[1:]
	}
	fieldName := strings.TrimPrefix(methodName, "SelectBy")
	fieldName = strings.TrimPrefix(fieldName, "FindBy")
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
		case strings.EqualFold(param.Name, fieldName):
			selectFieldType = fmt.Sprint(param.TypePrefix, param.Type)
		case param.Type == "Page":
			pageParamName = "page"
		default:
			selectFieldType = fmt.Sprint(param.TypePrefix, param.Type)
		}
	}
	for _, param := range method.Out {
		switch {
		case strings.EqualFold(param.Type, method.ModelName):
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
		Snippets: []template.RendererResult{result},
	}).ImportByName(`"context"`).
		ImportByName(string(buntemplates.CoreImport.With(method.Opts).Bytes())).
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
		Field: strcase.ToSnake(field),
		Order: order,
	}
}

func (h DbieToolHandler) OnEmbeddedInterface(e inspected.Entity) *inspected.Result {
	e.Opts = h.Opts
	log.Println("handles interface with entity", e)
	result := inspected.Result{
		Snippets: []template.RendererResult{
			buntemplates.Struct.With(e),
			buntemplates.Constr.With(e),
		},
	}
	return (&result).ImportByName(e.ModelPkg)
}
