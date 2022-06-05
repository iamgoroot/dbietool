package bun_templates

import (
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
)

var Struct = template.Tr[inspected.Entity](template.AtTypes, `
type generatedUserRepo struct {
	dbie.Repo[{{ .ModelIdent }}]
}
`)

var Constr = template.Tr[inspected.Entity](template.AtConstructor, `
func New{{ .Name }}(db *bun.DB, ctx context.Context) {{ .Name }} {
	return generatedUserRepo{
		Repo: dbie.NewRepo[{{ .ModelIdent }}](
			dbie.BunCore[{{ .ModelIdent }}]{Context: ctx, DB: db},
		),
	}
}
`)

var SelectByField = template.Tr[inspected.SelectorMethod](template.AtFunc,
	"{{if .OrderBy}}"+
		"var sort{{ .MethodName }}Setting = []dbie.Sort{\n"+
		"{{ range $val := .OrderBy }}\t{Field: `{{ $val.Field }}`, Order: {{ $val.Order }}},\n{{ end }}"+
		"}{{ end }}"+`

func (g generated{{ .Name }}) {{ .MethodName }}(
{{- if .PageParamName }}{{ .PageParamName }} dbie.Page, {{ end -}}
{{- .SelectFieldName }} {{ .SelectFieldType -}}
) ({{ .ModelIdent }}, error) {
	return g.{{ .CallMethod }}({{ if .PageParamName }}{{ .PageParamName }}, {{ end -}}`+
		"`\"{{- .SelectTableName }}\".\"{{ .SelectFieldName }}\"`"+
		`, dbie.{{ .Op }}, {{ .SelectFieldName }}{{if .OrderBy}}, sort{{ .MethodName }}Setting...{{ end }})		
}
`)
