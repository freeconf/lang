//go:build ignore
// +build ignore

// core_gen generates boilerplate functions for meta structs looking
// at specfic field names and generating functions based on that.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/freeconf/lang/codegen"
)

var homeDir = flag.String("codegen_dir", ".", "File path to directory containing _defs.go files")

func main() {
	flag.Parse()
	vars, err := codegen.ParseDefs(*homeDir)
	chkerr(err)
	for _, path := range flag.Args() {
		if !strings.HasSuffix(path, ".in") {
			chkerr(fmt.Errorf("expected .in file ext on '%s'", path))
		}
		destFname := path[:len(path)-3]
		fmt.Printf("%s => %s\n", path, destFname)
		dest, err := os.Create(destFname)
		if err != nil {
			chkerr(err)
		}
		defer dest.Close()
		err = codegen.GenerateSource(vars, path, dest)
		chkerr(err)
	}
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}
