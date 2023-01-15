package emeta

import (
	"bytes"
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestCodeGen(t *testing.T) {
	structs, err := ParseSource("./meta.go")
	fc.AssertEqual(t, nil, err)
	var buf bytes.Buffer
	err = GenerateSource(structs, "./templates/meta.c.tmpl", &buf)
	fc.AssertEqual(t, nil, err)
	t.Log(buf.String())
}

func TestReadGoSource(t *testing.T) {
	meta, err := ParseSource("./meta.go")
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, 5, len(meta.Definitions))
	fc.AssertEqual(t, "Module", meta.Definitions[0].Name)
	fc.AssertEqual(t, "Ident", meta.Definitions[0].Fields[0].Name)
	fc.AssertEqual(t, "string", meta.Definitions[0].Fields[0].Type)
}
