package generator

import (
	"bytes"
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
	"github.com/iancoleman/strcase"
	"go/ast"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Generator struct {
	inspect.TypeHandler[*inspected.Result]
	*goPackage
}

func (g *Generator) Run(types []string, dir string, outputPattern string) {
	g.parsePackageSources()
	src := g.inspectAndGenerate(&inspected.Result{})
	src = g.fmt(src)
	g.save(types, dir, src, outputPattern)
}

func (g *Generator) inspectAndGenerate(results *inspected.Result) []byte {
	inspector := inspect.Inspector{TypeHandler: g.TypeHandler, ImportLookup: results.ImportLookup}
	for _, file := range g.Files {
		if file.File != nil {
			ast.Inspect(file.File, func(node ast.Node) bool {
				inspector.Parse(node, results)
				return false
			})
		}
	}
	results.Add(
		template.Header.With(results.Pkg),
		template.Import.With(results.GetImports()),
	)
	return g.generateSrc(results)
}

func (g *Generator) fmt(src []byte) []byte {
	fm, err := format.Source(src)
	fmt.Println("formatted:\n", string(fm))
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return src
	}
	return src
}

func (g *Generator) save(types []string, dir string, src []byte, outputNamePattern string) {
	if outputNamePattern == "" {
		outputNamePattern = "%_generated.go"
	}
	baseName := fmt.Sprintf(outputNamePattern, strcase.ToSnake(types[0]))
	outputName := filepath.Join(dir, baseName)
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) generateSrc(result *inspected.Result) []byte {
	resultBuf := bytes.Buffer{}
	for _, snippet := range result.GetRenderers() {
		resultBuf.Write(snippet.Bytes())
	}
	return resultBuf.Bytes()
}
