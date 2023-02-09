package comm

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("../test/yang")
	m := parser.RequireModule(ypath, "testme")
	x := new(MetaEncoder).Encode(m)
	fc.AssertEqual(t, "testme", x.Ident)
	fc.AssertEqual(t, len(m.DataDefinitions()), len(x.Definitions))
}
