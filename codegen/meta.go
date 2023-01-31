package codegen

/**
* Parse meta_defs.go and use it to generate meta data about the definitions that can be used
* to generate defs in other languages
 */

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

// Yes, this is meta data about the Meta.
type MetaMeta struct {
	Definitions []*structDef
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

func ParseMetaDefs(homeDir string) (MetaMeta, error) {
	var empty MetaMeta
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath.Join(homeDir, "meta_defs.go"), nil, 0)
	if err != nil {
		return empty, err
	}
	g := newVistor()
	ast.Walk(g, f)
	return g.Meta, nil
}

func (def *structDef) HasField(name string) bool {
	for _, f := range def.Fields {
		if f.Name == name {
			return true
		}
	}
	return false
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
