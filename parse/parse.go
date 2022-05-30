package parse

import (
	"bytes"
	"fmt"
	"github.com/iamgoroot/dbietool/bun_render"
	"github.com/iamgoroot/dbietool/render"
	"go/ast"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type FileInspector interface {
	Parse(node ast.Node, result *render.Result)
}
type Gocrastinate struct {
	FileInspector func(map[string]string) FileInspector
	pkg           *Package
}

func (g *Gocrastinate) Run(types []string, dir string, output string) {
	src := g.inspectAndGenerate(&render.Result{})
	src = g.runFmt(src)
	g.persist(types, dir, src, output)
}

func (g *Gocrastinate) inspectAndGenerate(results *render.Result) []byte {
	for _, file := range g.pkg.Files {
		inspector := g.FileInspector(results.ImportLookup)
		if file.File != nil {
			ast.Inspect(file.File, func(node ast.Node) bool {
				inspector.Parse(node, results)
				return false
			})

		}
	}
	results.ImportByName(`"context"`)
	results.ImportByName(`"github.com/uptrace/bun"`)
	results.ImportByName(`"github.com/iamgoroot/dbie"`)
	results.Add(
		bun_render.Header.With(results.Pkg),
		bun_render.Import.With(results.GetImports()),
	)
	src := g.generate(results)
	return src
}

func (g *Gocrastinate) runFmt(src []byte) []byte {
	fm, err := format.Source(src)
	fmt.Println("formatted:\n", string(fm))
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return src
	}
	return src
}

func (g *Gocrastinate) persist(types []string, dir string, src []byte, outputName string) {
	if outputName == "" {
		baseName := fmt.Sprintf("%s_generated.go", types[0])
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Gocrastinate) generate(result *render.Result) []byte {
	resultBuf := bytes.Buffer{}
	for _, snippet := range result.GetRenderers() {
		resultBuf.Write(snippet.Bytes())
	}
	return resultBuf.Bytes()
}
