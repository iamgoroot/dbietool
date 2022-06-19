package inspected

import "fmt"

type (
	Entity struct {
		Name, Pkg                        string
		ModelName, ModelPkg, ModelImport string
	}

	Method struct {
		Entity
		MethodName string
		In, Out    []Param
	}

	Param struct {
		Name, TypePrefix, Type string
		TypeParam              *Param
	}
	SelectorMethod struct {
		Method
		SelectFieldName, SelectFieldType string
		CallMethod                       string
		Op                               string
		PageParamName                    string
		SelectTableName                  string
		OrderBy                          []OrderBy
	}
	OrderBy struct {
		Field, Order string
	}
)

func (e Entity) ModelIdent() string {
	if e.ModelPkg == "" {
		return e.ModelName
	}
	return fmt.Sprintf("%s.%s", e.ModelPkg, e.ModelName)
}
