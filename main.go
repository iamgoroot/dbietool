package main

import (
	"flag"
	"fmt"
	"github.com/iamgoroot/dbietool/generator"
	"github.com/iamgoroot/dbietool/generator/inspect/handler"
	"os"
)

var (
	typeNames   = flag.String("type", "", "comma-separated list of type names; must be set?")
	filePattern = flag.String("filepattern", "%s_generated.go", "output file patten '%s_generated.go'")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of dbietool:\n")
	flag.PrintDefaults()
}

func main() {
	types := parseCmd()
	g := generator.Generator{TypeHandler: handler.DbieToolHandler{}}
	g.Run(types, pwd(), *filePattern)
}
