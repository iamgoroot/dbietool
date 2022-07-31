package main

import (
	"flag"
	"github.com/iamgoroot/dbietool/generator"
	"github.com/iamgoroot/dbietool/generator/inspect/handler"
	"log"
	"strings"
)

var (
	output = flag.String("output", "generated.go", "output file name")
	cores  = flag.String("core", "Bun", "core: Bun,Gorm,Beego")
	constr = flag.String("constr", "factory", "constr: factory,func")

	typeNamePattern = flag.String("typeName", "%sImpl", "type name pattern")
)

func main() {
	flag.Parse()

	err := generator.New(
		handler.DbieToolHandler{
			Opts: map[string]interface{}{
				"Cores":           strings.Split(*cores, ","),
				"Constr":          *constr,
				"TypeNamePattern": *typeNamePattern,
			},
		},
	).Out(*output)

	if err != nil {
		log.Fatalln(err)
	}
}
