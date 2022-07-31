package generator

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

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
