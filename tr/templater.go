package tr

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"log"
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
	tmpl, err := template.New("dbie_example").Parse(render.Template)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return templateRendererResult[Data]{
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
