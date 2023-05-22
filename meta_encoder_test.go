package lang

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

func TestBasicMetaEncoder(t *testing.T) {
	ypath := source.Dir("./test/testdata/yang")
	m := parser.RequireModule(ypath, "testme")
	x := new(MetaEncoder).Encode(m)
	fc.AssertEqual(t, "testme", x.Ident)
	fc.AssertEqual(t, len(m.DataDefinitions()), len(x.Definitions))
	fc.AssertEqual(t, "x", x.Definitions[0].GetLeaf().Ident)
	fc.AssertEqual(t, "z", x.Definitions[1].GetContainer().Ident)
}

/*
parse all files just to flush out encoding errs
*/
func TestMetaEncoder(t *testing.T) {
	d := "./test/testdata/yang"
	ypath := source.Dir(d)
	files, err := ioutil.ReadDir(d)
	fc.RequireEqual(t, nil, err)
	for _, f := range files {
		fname := f.Name()
		if strings.HasSuffix(fname, ".yang") {
			name := fname[:len(fname)-5]
			t.Log(name)
			m := parser.RequireModule(ypath, name)
			new(MetaEncoder).Encode(m)
		}
	}
}
