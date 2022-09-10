package template

import (
	"bytes"
	"fmt"
	"github.com/iancoleman/strcase"
	"hash/fnv"
	"strings"
	"text/template"
)

type CodePosition int

const (
	AtPackage CodePosition = 0 + iota
	AtImport
	AtConstants
	AtTypes
	AtConstructor
	AtMethod
	AtFunc
	AtEOF
)

func Tr[Data any](w CodePosition, t string) templateRenderer[Data] {
	return templateRenderer[Data]{
		Template: t,
		Weight:   int(w),
	}
}

type Renderer[Data any] interface {
	With(d Data) RendererResult
}

type RendererResult interface {
	Bytes() []byte
	Weight() int
	ID() string
	Imports() []string
	Error() error
}

type templateRenderer[Data any] struct {
	Template string
	Weight   int
	Imports  []string
}

func (render templateRenderer[Data]) Import(name string) templateRenderer[Data] {
	render.Imports = append(render.Imports, name)
	return render
}

func (render templateRenderer[Data]) With(data Data) RendererResult {
	buf := bytes.Buffer{}
	tmpl, err := template.New("dbie").
		Funcs(map[string]any{
			"toTitle":      strings.ToTitle,
			"toLower":      strings.ToLower,
			"toUpper":      strings.ToUpper,
			"trim":         strings.Trim,
			"split":        strings.Split,
			"join":         strings.Join,
			"replace":      strings.Replace,
			"replaceAll":   strings.ReplaceAll,
			"repeat":       strings.Repeat,
			"contains":     strings.Contains,
			"containsAny":  strings.ContainsAny,
			"containsRune": strings.ContainsRune,
			"count":        strings.Count,
			"hasPrefix":    strings.HasPrefix,
			"hasSuffix":    strings.HasSuffix,
			"index":        strings.Index,
			"indexAny":     strings.IndexAny,
			"lastIndex":    strings.LastIndex,
			"lastIndexAny": strings.LastIndexAny,
			"toCamel":      strcase.ToCamel,
			"toSnake":      strcase.ToSnake,
			"toKebab":      strcase.ToKebab,
			"toLowerCamel": strcase.ToLowerCamel,
			"toDelimited":  strcase.ToDelimited,
		},
		).Parse(render.Template)
	if err == nil {
		err = tmpl.Execute(&buf, data)
	}
	return templateRendererResult[Data]{
		Err:      err,
		bytes:    buf.Bytes(),
		weight:   render.Weight,
		uniqueID: fmt.Sprint(data, hash(render.Template)),
	}
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type templateRendererResult[Data any] struct {
	uniqueID string
	bytes    []byte
	weight   int
	imports  []string
	Err      error
}

func (t templateRendererResult[Data]) Error() error {
	return t.Err
}

func (t templateRendererResult[Data]) Imports() []string {
	return t.imports
}

func (t templateRendererResult[Data]) Bytes() []byte {
	return t.bytes
}

func (t templateRendererResult[Data]) Weight() int {
	return t.weight
}

func (t templateRendererResult[Data]) ID() string {
	return t.uniqueID
}
