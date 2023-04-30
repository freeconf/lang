package codegen

import (
	"bytes"
	"flag"
	"testing"

	"github.com/freeconf/yang/fc"
)

var update = flag.Bool("update", false, "update gold files instead of testing against them")

func TestMetaDefs(t *testing.T) {
	meta, err := ParseProtos("../")
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, true, len(meta.AllMessages) > 0)
}

func TestCodeGen(t *testing.T) {
	vars, err := ParseProtos("../")
	fc.AssertEqual(t, nil, err)
	var buf bytes.Buffer
	err = GenerateSource(vars, "./testdata/testme.txt.in", &buf)
	fc.AssertEqual(t, nil, err)
	fc.Gold(t, *update, buf.Bytes(), "./testdata/gold/testme.txt")
}

func TestFuncs(t *testing.T) {
	fc.AssertEqual(t, "SomeCase", camel("SOME_CASE"))
}
