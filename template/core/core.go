package bun_templates

import (
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
)

const dbObjectIdent = `
	{{- if eq . "Bun" }} *bun.DB{{ end -}}
	{{- if eq . "Gorm" }} *gorm.DB{{ end -}}
	{{- if eq . "Bee" }} orm.Ormer{{ end -}}
`

var Factory = template.Tr[inspected.Opts](template.AtConstructor, `
{{- if eq .Constr "factory" }}
{{range .Cores }}
type {{ . }} struct {
`+dbObjectIdent+`
}
{{ end -}}
{{ end -}}
`)

var Constr = template.Tr[inspected.Entity](template.AtConstructor, `
{{- if eq .Opts.Constr "func" }}
{{ $entity := . }}
{{range .Opts.Cores }}
func New{{ . }}{{ $entity.Name }}(ctx context.Context, db `+dbObjectIdent+`) {{ $entity.Name }} {
	return {{ $entity.Name | printf $entity.Opts.TypeNamePattern }}{ Repo: dbie.New{{ . }}[{{ $entity.ModelIdent}}](ctx, db) }
}
{{ end }}
{{ end -}}

{{- if eq .Opts.Constr "factory" -}}
{{ $entity := . }}
{{range .Opts.Cores }}
func (factory {{ . }}) New{{ . }}{{ $entity.Name }}(ctx context.Context) {{ $entity.Name }} {
	return {{ $entity.Name | printf $entity.Opts.TypeNamePattern }}{ Repo: core{{ . }}.New[{{ $entity.ModelIdent}}](ctx, factory.DB) }
}
{{ end -}}
{{ end -}}
`)

var CoreImport = template.Tr[inspected.Opts](template.AtImport, `
{{- range .Cores }}
	core{{ . }} "github.com/iamgoroot/dbie/core/{{ . | toSnake }}"
	{{ if eq . "Bun" }}"github.com/uptrace/bun"{{ end -}}
	{{ if eq . "Gorm" }}"gorm.io/gorm"{{ end -}}
	{{ if eq . "Bee" }}"github.com/beego/beego/v2/client/orm"{{ end -}}
{{- end -}}
`)

var Struct = template.Tr[inspected.Entity](template.AtTypes, `
type {{ .Name | printf .Opts.TypeNamePattern }} struct {
	dbie.Repo[{{ .ModelIdent }}]
}
`)

var SelectByField = template.Tr[inspected.SelectorMethod](template.AtFunc,
	"{{if .OrderBy}}"+
		"var sort{{ .MethodName }}Setting = []dbie.Sort{\n"+
		"{{ range $val := .OrderBy }}\t{Field: `{{ $val.Field }}`, Order: {{ $val.Order }}},\n{{ end }}"+
		"}{{ end }}"+`

func (g {{ .Name | printf .Opts.TypeNamePattern }}) {{ .MethodName }}(
{{- if .PageParamName }}{{ .PageParamName }} dbie.Page, {{ end -}}
{{- .SelectFieldName }} {{ .SelectFieldType -}}
) ({{ .ModelIdent }}, error) {
	return g.Repo.{{ .CallMethod }}({{ if .PageParamName }}{{ .PageParamName }}, {{ end -}}`+
		"\"{{ .SelectFieldName }}\""+
		`, dbie.{{ .Op }}, {{ .SelectFieldName }}{{if .OrderBy}}, sort{{ .MethodName }}Setting...{{ end }})		
}
`)
