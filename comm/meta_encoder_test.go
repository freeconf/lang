package comm

import (
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("../test/yang")
	m := parser.RequireModule(ypath, "testme-1")
	x := new(metaEncoder).encode(m)
	fc.AssertEqual(t, "testme-1", x.Ident)
}
