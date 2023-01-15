package emeta

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
	"github.com/ugorji/go/codec"
)

func TestMetaEncoder(t *testing.T) {
	ypath := source.Dir("../../../yang")
	m := parser.RequireModule(ypath, "testme")
	var buf bytes.Buffer
	fc.AssertEqual(t, nil, Encode(m, &buf))

	var hnd codec.CborHandle
	d := codec.NewDecoder(&buf, &hnd)
	var data []interface{}
	fc.AssertEqual(t, nil, d.Decode(&data))

	var dump bytes.Buffer
	dumper := json.NewEncoder(&dump)
	dumper.SetIndent("", "  ")
	dumper.Encode(data)
	t.Log(dump.String())
}
