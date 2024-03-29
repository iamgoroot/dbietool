package generator

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
)

type gofile struct {
	*goPackage
	*ast.File
	TypesSet []string
}

type goPackage struct {
	Name  string
	Defs  map[*ast.Ident]types.Object
	Files []*gofile
}

func (g *generator) parsePackageSources() {
	cfg := &packages.Config{
		Mode:  packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedImports | packages.NeedTypes,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

func (g *generator) addPackage(pkg *packages.Package) {
	g.goPackage = &goPackage{
		Defs:  pkg.TypesInfo.Defs,
		Files: make([]*gofile, len(pkg.Syntax)),
	}
	for i, f := range pkg.Syntax {
		g.Files[i] = &gofile{
			File:      f,
			goPackage: g.goPackage,
		}
	}
}
