package codegen

import (
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestParseMetaProto(t *testing.T) {
	defs, err := ParseMetaDefs("..")
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "Module", defs.Definitions[0].Name)
	fc.AssertEqual(t, "ident", defs.Definitions[0].Fields[0].Name)
	fc.AssertEqual(t, "string", defs.Definitions[0].Fields[0].Type)
	fc.AssertEqual(t, "container", defs.DataDefs[0].Name)
}
