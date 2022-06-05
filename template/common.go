package template

var Header = Tr[string](AtPackage, `package {{ . }}
`)

var Import = Tr[map[string]string](AtImport, `
import (
	{{ range $key, $val := . }}
		 {{- $key }} {{ $val }}
	{{ end }}
)

`)
