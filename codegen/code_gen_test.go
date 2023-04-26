package codegen

import (
	"bytes"
	"flag"
	"testing"

	"github.com/freeconf/yang/fc"
)

var update = flag.Bool("update", false, "update gold files instead of testing against them")

func TestMetaDefs(t *testing.T) {
	meta, err := ParseMetaDefs("../")
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 8, len(meta.Definitions))
	fc.AssertEqual(t, "ExtensionDefArg", meta.Definitions[0].Name)
	fc.AssertEqual(t, "MetaId", meta.Definitions[0].Fields[0].Name)
	fc.AssertEqual(t, "MetaId", meta.Definitions[0].Fields[0].Type)
}

func TestCodeGen(t *testing.T) {
	vars, err := ParseMetaDefs("../")
	fc.AssertEqual(t, nil, err)
	var buf bytes.Buffer
	err = GenerateSource(Vars{Meta: vars}, "./testdata/testme.txt.in", &buf)
	fc.AssertEqual(t, nil, err)
	fc.Gold(t, *update, buf.Bytes(), "./testdata/gold/testme.txt")
}

func TestFuncs(t *testing.T) {
	fc.AssertEqual(t, "SomeCase", camel("SOME_CASE"))
}
