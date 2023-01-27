package codegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:generate go run code_gen_main.go

func ParseSource(src string) (MetaMeta, error) {
	var empty MetaMeta
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		return empty, err
	}
	g := newVistor()
	ast.Walk(g, f)
	return g.Meta, nil
}

type structDef struct {
	Name   string
	MetaId int
	Fields []*fieldDef
}

type fieldDef struct {
	Name string
	Type string
	Tags []string
}

func (def *structDef) HasField(name string) bool {
	for _, f := range def.Fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

// Yes, this is meta data about the Meta.
type MetaMeta struct {
	Definitions []*structDef
}

type visitor struct {
	Meta          MetaMeta
	metaIds       map[string]int
	currentStruct *structDef
	currentField  *fieldDef
}

func newVistor() *visitor {
	return &visitor{
		metaIds: make(map[string]int),
	}
}

func (g *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch x := n.(type) {
	case *ast.TypeSpec:
		switch x.Type.(type) {
		case *ast.StructType:
			name := x.Name.Name
			defTypeName := "MetaId" + name
			metaId, valid := g.metaIds[defTypeName]
			if !valid {
				panic(fmt.Sprintf("'%s' not defined", defTypeName))
			}
			def := &structDef{
				Name:   name,
				MetaId: metaId,
			}
			g.Meta.Definitions = append(g.Meta.Definitions, def)
			g.currentStruct = def
		}
	case *ast.ValueSpec:
		if strings.HasPrefix(x.Names[0].Name, "MetaId") {
			g.metaIds[x.Names[0].Name] = len(g.metaIds)
		}
	case *ast.Field:
		if g.currentStruct != nil {
			name := x.Names[0].Name
			def := &fieldDef{
				Name: name,
				Type: types.ExprString(x.Type),
			}
			if x.Tag != nil {
				def.Tags = []string{x.Tag.Value}
			}
			g.currentStruct.Fields = append(g.currentStruct.Fields, def)
			g.currentField = def
		}
		// default:
		// 	fmt.Printf("unaccounted for %T, %v\n", n, n)
	}
	return g
}

func GenerateSource(src MetaMeta, tmpl string, out io.Writer) error {
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
	vars := struct {
		Meta MetaMeta
	}{
		Meta: src,
	}
	if err := t.Execute(out, vars); err != nil {
		panic(err)
	}
	return nil
}
