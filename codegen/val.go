package codegen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

/**
* Parse val_defs.go and use it to generate meta data about the definitions that can be used
* to generate defs in other languages
 */

type ValMeta struct {
	Definitions []*valDef
}

type valDef struct {
	Name   string
	IsList bool
	TypeId int
}

func ParseValDefs(homeDir string) (ValMeta, error) {
	var empty ValMeta
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, homeDir+"codegen/val_defs.go", nil, 0)
	if err != nil {
		return empty, err
	}
	g := newValDefParser()
	ast.Walk(g, f)
	return g.Meta, nil
}

type valDefParser struct {
	Meta ValMeta
}

func newValDefParser() *valDefParser {
	return &valDefParser{}
}

func (g *valDefParser) Visit(n ast.Node) ast.Visitor {
	typeTypePrefix := "ValType"
	if n == nil {
		return nil
	}
	switch x := n.(type) {
	case *ast.ValueSpec:
		specName := x.Names[0].Name
		if strings.HasPrefix(specName, typeTypePrefix) {
			name := specName[len(typeTypePrefix):]
			g.Meta.Definitions = append(g.Meta.Definitions, &valDef{
				Name:   name,
				TypeId: len(g.Meta.Definitions),
				IsList: strings.HasSuffix(name, "List"),
			})
		}
	}
	return g
}
