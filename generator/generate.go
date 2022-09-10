package generator

import (
	"bytes"
	"fmt"
	"github.com/iamgoroot/dbietool/generator/inspect"
	"github.com/iamgoroot/dbietool/generator/inspect/inspected"
	"github.com/iamgoroot/dbietool/template"
	"go/ast"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type generator struct {
	inspect.TypeHandler[*inspected.Result]
	*goPackage
}

type Generator interface {
	Out(string) error
}

func New(typeHandler inspect.TypeHandler[*inspected.Result]) Generator {
	g := generator{
		TypeHandler: typeHandler,
		goPackage:   nil,
	}
	g.parsePackageSources()
	return &g
}

func (g *generator) Out(output string) error {
	src, errs := g.inspectAndGenerate(&inspected.Result{})
	if len(errs) > 0 {
		log.Fatalln(errs)
	}
	src = g.fmt(src)
	g.save(src, output)
	return nil //TODO: err handling
}

func (g *generator) inspectAndGenerate(results *inspected.Result) ([]byte, []error) {
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
		inspector.TypeHandler.Once(),
	)
	return generateSrc(results)
}

func generateSrc(result *inspected.Result) ([]byte, []error) {
	resultBuf := bytes.Buffer{}
	var errs []error
	for _, snippet := range result.GetSnippets() {
		resultBuf.Write(snippet.Bytes())
		if err := snippet.Error(); err != nil {
			errs = append(errs, err)
		}
	}
	return resultBuf.Bytes(), errs
}

func (g *generator) fmt(src []byte) []byte {
	fmt.Println("generated:\n", string(src))
	fm, err := format.Source(src)
	fmt.Println("formatted:\n", string(fm))
	if err != nil {
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return src
	}
	return src
}

func (g *generator) save(src []byte, output string) {
	outputName := filepath.Join(pwd(), output)
	_ = os.Mkdir(filepath.Dir(outputName), 0750)
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
