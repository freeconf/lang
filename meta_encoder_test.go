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
	z := x.Definitions[0].GetContainer()
	fc.AssertEqual(t, "z", z.Ident)
	fc.AssertEqual(t, "a", z.Definitions[0].GetLeaf().Ident)
	z2 := z.Definitions[1].GetContainer()
	fc.AssertEqual(t, "z", z2.Ident)
	fc.AssertEqual(t, "", z2.RecursivePath)
	fc.AssertEqual(t, true, z2.IsRecursive)
	p3 := z2.Definitions[1].GetPtr()
	fc.AssertEqual(t, "z/z", p3.Path)
}
