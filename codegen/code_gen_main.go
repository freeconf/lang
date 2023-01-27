//go:build ignore
// +build ignore

// core_gen generates boilerplate functions for meta structs looking
// at specfic field names and generating functions based on that.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/freeconf/lang/codegen"
)

var dirs = []string{
	"..",
}

func main() {
	metas, err := codegen.ParseSource("./meta.go")
	chkerr(err)
	for _, dir := range dirs {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !strings.HasSuffix(path, ".in") {
				return nil
			}
			destFname := path[:len(path)-3]
			fmt.Printf("%s => %s\n", path, destFname)
			dest, err := os.Create(destFname)
			if err != nil {
				return err
			}
			defer dest.Close()
			err = codegen.GenerateSource(metas, path, dest)
			return err
		})
		chkerr(err)
	}
}

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}
