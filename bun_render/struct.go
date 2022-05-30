package bun_render

import (
	"github.com/iamgoroot/dbietool/models"
	"github.com/iamgoroot/dbietool/tr"
)

var Header = tr.Tr[string](tr.AtPackage, `package {{ . }}
`)

var Import = tr.Tr[map[string]string](tr.AtImport, `
import (
	{{ range $key, $val := . }}
		 {{- $key }} {{ $val }}
	{{ end }}
)

`)

var Struct = tr.Tr[models.Entity](tr.AtTypes, `
type generatedUserRepo struct {
	dbie.BunBackend[{{ .ModelIdent }}]
}
`)

var Constr = tr.Tr[models.Entity](tr.AtConstructor, `
func New{{ .Name }}(db *bun.DB, ctx context.Context) {{ .Name }} {
	return generatedUserRepo{
		BunBackend: dbie.BunBackend[{{ .ModelIdent }}]{DB: db, Context: ctx},
	}
}
`)
