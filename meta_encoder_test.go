package lang

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("./test/yang")
	m := parser.RequireModule(ypath, "testme")
	x := new(MetaEncoder).Encode(m)
	fc.AssertEqual(t, "testme", x.Ident)
	fc.AssertEqual(t, len(m.DataDefinitions()), len(x.Definitions))
	fc.AssertEqual(t, "x", x.Definitions[0].GetLeaf().Ident)
	fc.AssertEqual(t, "z", x.Definitions[1].GetContainer().Ident)
}
