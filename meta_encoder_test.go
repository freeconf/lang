package lang

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestBasicMetaEncoder(t *testing.T) {
	ypath := source.Dir("./test/testdata/yang")
	m := parser.RequireModule(ypath, "testme")
	x := NewMetaEncoder().Encode(m)
	fc.AssertEqual(t, "testme", x.Ident)
	fc.AssertEqual(t, len(m.DataDefinitions()), len(x.Definitions))
	fc.AssertEqual(t, "x", x.Definitions[0].GetLeaf().Ident)
	fc.AssertEqual(t, "z", x.Definitions[1].GetContainer().Ident)
}

var testFiles = []string{
	"car",
	"echo",
	"meta",
	"testme",
	"advanced",
}

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("./test/testdata/yang")
	for _, f := range testFiles {
		t.Log(f)
		m := parser.RequireModule(ypath, f)
		NewMetaEncoder().Encode(m)
	}
}

func TestRecurse(t *testing.T) {
	ypath := source.Dir("./test/testdata/yang")
	m := parser.RequireModule(ypath, "recurse")
	x := NewMetaEncoder().Encode(m)
	zdef := x.Definitions[0]
	z := zdef.GetContainer()
	fc.AssertEqual(t, "z", z.Ident)
	fc.AssertEqual(t, "a", z.Definitions[0].GetLeaf().Ident)
	fc.AssertEqual(t, "z", z.Definitions[1].GetContainer().Ident)

	//fc.AssertEqual(t, true, z == z.Definitions[1].GetContainer())
}
