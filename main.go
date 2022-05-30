package main

import (
	"flag"
	"fmt"
	"github.com/iamgoroot/dbietool/inspect"
	"github.com/iamgoroot/dbietool/parse"
	"os"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name")
	buildTags = flag.String("tags", "", "comma-separated list of build tags to apply")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of stringer:\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tstringer [flags] -type T files... # Must be a single package\n")
	flag.PrintDefaults()
}

func main() {
	types, tags, args := parseCmd()
	g := parse.Gocrastinate{
		FileInspector: func(lookup map[string]string) parse.FileInspector {
			return &inspect.SingleFile{WithTypes: types}
		},
	}
	dir := pwd(args, tags)
	g.ParsePackage(args, tags)
	g.Run(types, dir, *output)
}
