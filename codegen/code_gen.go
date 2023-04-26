package codegen

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

type Vars struct {
	Meta MetaMeta
}

func ParseDefs(homeDir string) (vars Vars, err error) {
	if vars.Meta, err = ParseMetaDefs(homeDir); err != nil {
		return
	}
	return
}

func title(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func whisperingSnake(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}

func camel(s string) string {
	return title(strcase.ToLowerCamel(strings.ToLower(s)))
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
		"camel": camel,
	}
	t, err := template.New("code_gen").Funcs(funcs).Parse(string(tmplSrc))
	if err != nil {
		panic(err)
	}
	if err := t.Execute(out, vars); err != nil {
		panic(err)
	}
	return nil
}
