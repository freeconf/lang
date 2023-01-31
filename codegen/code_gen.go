package codegen

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:generate go run code_gen_main.go

type Vars struct {
	Meta MetaMeta
	Val  ValMeta
}

func ParseDefs(homeDir string) (vars Vars, err error) {
	if vars.Meta, err = ParseMetaDefs(homeDir); err != nil {
		return
	}
	if vars.Val, err = ParseValDefs(homeDir); err != nil {
		return
	}
	return
}

func GenerateSource(vars Vars, tmpl string, out io.Writer) error {
	tmplFile, err := os.Open(tmpl)
	if err != nil {
		panic(err)
	}
	tmplSrc, err := ioutil.ReadAll(tmplFile)
	if err != nil {
		panic(err)
	}
	funcs := template.FuncMap{
		"lc":    strings.ToLower,
		"uc":    strings.ToUpper,
		"snake": strcase.ToSnake,
	}
	cTemplateFuncs(funcs)
	t, err := template.New("code_gen").Funcs(funcs).Parse(string(tmplSrc))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(out, vars); err != nil {
		panic(err)
	}
	return nil
}
