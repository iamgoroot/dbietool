package bun_render

import (
	"github.com/iamgoroot/dbietool/models"
	"github.com/iamgoroot/dbietool/tr"
)

var SelectByField = tr.Tr[models.SelectorMethod](tr.AtFunc, "\nconst fieldName{{ .MethodName }} = `\"{{ .SelectTableName }}\".\"{{ .SelectFieldName }}\"`\n"+`
func (g generated{{ .Name }}) {{ .MethodName }}(
{{- if .PageParamName }}{{ .PageParamName }} dbie.Page, {{ end -}}
{{- .SelectFieldName }} {{ .SelectFieldType -}}
) ({{ .ModelIdent }}, error) {
	return g.{{ .CallMethod }}({{ if .PageParamName }}{{ .PageParamName }}, {{ end -}}fieldName{{ .MethodName }}, dbie.{{ .Op }}, {{ .SelectFieldName }})
}
`)

//var SelectManyByField = tr.Tr[models.Entity](tr.AtConstructor, `
//func {{ .MethodName }}({{ range $key, $val := .In }} {{ $key }} {{ $val }}{{ end }}) {{ range $key, $val := .Out }} {{ $key }} {{ $val }}{{ end }}{
//}
//`)
//{{- if eq .Op "In"}}...{{ end -}}
//{{- if eq .Op "Nin"}}...{{ end -}}
