package emeta

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"

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
	g := &visitor{}
	ast.Walk(g, f)
	return g.Meta, nil
}

type structDef struct {
	Name   string
	Fields []*fieldDef
}

type fieldDef struct {
	Name string
	Type string
	Tags []string
}

// Yes, this is meta data about the Meta.
type MetaMeta struct {
	Definitions     []*structDef
	DefinitionTypes []string
}

type visitor struct {
	Meta          MetaMeta
	currentStruct *structDef
	currentField  *fieldDef
}

func (g *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch x := n.(type) {
	case *ast.TypeSpec:
		switch x.Type.(type) {
		case *ast.StructType:
			def := &structDef{
				Name: x.Name.Name,
			}
			g.Meta.Definitions = append(g.Meta.Definitions, def)
			g.currentStruct = def
		}
	case *ast.ValueSpec:
		if strings.HasPrefix(x.Names[0].Name, "DefType") {
			g.Meta.DefinitionTypes = append(g.Meta.DefinitionTypes, x.Names[0].Name)
		}
	case *ast.Field:
		if g.currentStruct != nil {
			def := &fieldDef{
				Name: x.Names[0].Name,
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

func cFieldType(f *fieldDef) string {
	switch f.Type {
	case "string":
		return "char*"
	case "*bool":
		return "bool*"
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
