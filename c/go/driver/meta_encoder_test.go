package driver

import (
	"bytes"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("../../../yang")
	m := parser.RequireModule(ypath, "testme")
	var buf bytes.Buffer
	fc.AssertEqual(t, nil, new(Encoder).Encode(m, &buf))
	t.Logf("buf len %d", buf.Len())

	rt, err := new(Decoder).Decode(&buf)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "testme", rt.Ident())
	fc.AssertEqual(t, "testing loading yang files", rt.Description())
}

func TestMetaEncoder2(t *testing.T) {
	ypath := source.Dir("../../../yang")
	m := parser.RequireModule(ypath, "testme")
	var buf bytes.Buffer
	fc.AssertEqual(t, nil, Encode2(m, &buf))
	t.Logf("buf len %d", buf.Len())

	rt, err := Decode2(&buf)
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "testme", rt.Ident())
	fc.AssertEqual(t, "testing loading yang files", rt.Description())
	fc.AssertEqual(t, 1, len(rt.DataDefinitions()))
}
