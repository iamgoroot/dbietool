package models

import "fmt"

type (
	Entity struct {
		Name                             string
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
	}
)

func (e Entity) ModelIdent() string {
	if e.ModelPkg == "" {
		return e.ModelName
	}
	return fmt.Sprintf("%s.%s", e.ModelPkg, e.ModelName)
}

//func (e SelectorMethod) ModelIdent() string {
//	return e.SelectFieldType
//}

//func (m Method) OutParams() string {
//	b := strings.Builder{}
//	for _, param := range m.Out {
//
//	}
//	b.WriteString()
//}

//func (p Params) out(entityName string) string {
//	for i, param := range p {
//		if strings.HasSuffix(param.Type, entityName) {
//
//		}
//	}
//}
