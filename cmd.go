package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func parseCmd() []string {
	log.SetFlags(0)
	log.SetPrefix("dbietool: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")
	return types
}

func pwd() string {
	args := flag.Args()
	if len(args) == 0 {
		return filepath.Dir(".")
	}
	if len(args) == 1 && isDirectory(args[0]) {
		return args[0]
	}
	return filepath.Dir(args[0])
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
