package lang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestMeta(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "../yang/meta/core.go", nil, 0)
	fc.AssertEqual(t, nil, err)
	g := gen{structs: make(map[string]*structDef)}
	ast.Walk(g, f)
	t.Logf("%v", g.structs)
}

type structDef struct {
	Name   string
	Fields []*fieldDef
}

type fieldDef struct {
	Name string
}

type gen struct {
	structs       map[string]*structDef
	currentStruct *structDef
	currentField  *fieldDef
}

func (g gen) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch x := n.(type) {
	case *ast.TypeSpec:
		def := &structDef{Name: x.Name.Name}
		g.structs[def.Name] = def
		g.currentStruct = def
	case *ast.Field:
		if g.currentStruct != nil {
			def := &fieldDef{Name: x.Names[0].Name}
			g.currentStruct.Fields = append(g.currentStruct.Fields, def)
			g.currentField = def
		}
	}
	return g
}
