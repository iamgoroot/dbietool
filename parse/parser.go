package parse

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
)

type File struct {
	*Package           // Package to which this file belongs.
	*ast.File          // Parsed AST.
	TypesSet  []string // MethodName of the constant type.
}

type Package struct {
	Name  string
	Defs  map[*ast.Ident]types.Object
	Files []*File
}

func (g *Gocrastinate) ParsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedTypes,
		// TODO: Need to think about constants in dbie_example files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

func (g *Gocrastinate) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		Name:  pkg.Name,
		Defs:  pkg.TypesInfo.Defs,
		Files: make([]*File, len(pkg.Syntax)),
	}
	for i, file := range pkg.Syntax {
		g.pkg.Files[i] = &File{
			File:    file,
			Package: g.pkg,
		}
	}
}
