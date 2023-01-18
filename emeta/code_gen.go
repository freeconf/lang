package emeta

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
	Name       string
	EncodingId int
	Fields     []*fieldDef
}

type fieldDef struct {
	Name string
	Type string
	Tags []string
}

// Yes, this is meta data about the Meta.
type MetaMeta struct {
	Definitions []*structDef
}

type visitor struct {
	Meta          MetaMeta
	encodingIds   map[string]int
	currentStruct *structDef
	currentField  *fieldDef
}

func newVistor() *visitor {
	return &visitor{
		encodingIds: make(map[string]int),
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
			defTypeName := "EncodingId" + name
			encodingId, valid := g.encodingIds[defTypeName]
			if !valid {
				panic(fmt.Sprintf("'%s' not defined", defTypeName))
			}
			def := &structDef{
				Name:       name,
				EncodingId: encodingId,
			}
			g.Meta.Definitions = append(g.Meta.Definitions, def)
			g.currentStruct = def
		}
	case *ast.ValueSpec:
		if strings.HasPrefix(x.Names[0].Name, "EncodingId") {
			g.encodingIds[x.Names[0].Name] = len(g.encodingIds)
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
		"lc":         strings.ToLower,
		"uc":         strings.ToUpper,
		"snake":      strcase.ToSnake,
		"cDefType":   cDefType,
		"cFieldType": cFieldType,
		"cName":      cName,
		"cNameSafe":  cNameSafe,
	}
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

func cDefType(d *structDef) string {
	return "fc_" + whisperingSnake(d.Name)
}

func cName(f *fieldDef) string {
	return whisperingSnake(f.Name)
}

func whisperingSnake(s string) string {
	return strings.ToLower(strcase.ToSnake(s))
}

func cNameSafe(name string) string {
	return strings.ReplaceAll(name, "*", "_ptr")
}

func cFieldType(f *fieldDef) string {
	switch f.Type {
	case "string":
		return "char*"
	case "[]string":
		return "char**"
	case "*bool":
		return "bool*"
	case "bool":
		return "bool"
	case "int":
		return "int"
	}
	switch f.Name {
	case "Definitions":
		return "fc_datadef_ptr*"
	}
	ctype := whisperingSnake(f.Type)
	if strings.HasPrefix(f.Type, "[]") {
		ctype = ctype[2:] + "**"
	}
	if strings.HasPrefix(f.Type, "*") {
		ctype = ctype[1:] + "*"
	}

	return "fc_" + ctype
}
