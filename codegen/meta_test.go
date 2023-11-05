package codegen

import (
	"testing"

	"github.com/freeconf/yang/fc"
)

func TestParseMetaProto(t *testing.T) {
	defs, err := ParseProtos("..")
	fc.AssertEqual(t, nil, err)
	fc.AssertEqual(t, "container", defs.MessagesByName["DataDef"].OneOfs["def_oneof"][0].Name)
	fc.AssertEqual(t, "BINARY", defs.ValEnums[0].Ident)
}
