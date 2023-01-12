package driver

import (
	"fmt"
	"io"

	"github.com/freeconf/yang/meta"
	"github.com/ugorji/go/codec"
)

func Encode2(m *meta.Module, out io.Writer) error {
	var hnd codec.CborHandle
	e := codec.NewEncoder(out, &hnd)
	return e.Encode(meta2driver(m))
}

type DataDef struct {
	Ident string
	Type  string
}

type Module struct {
	Ident       string
	Description string
	DataDef     []*DataDef
}

func meta2driver(m *meta.Module) *Module {
	d := &Module{
		Ident:       m.Ident(),
		Description: m.Description(),
	}
	d.DataDef = make([]*DataDef, len(m.DataDefinitions()))
	for i, def := range m.DataDefinitions() {
		d.DataDef[i] = &DataDef{
			Ident: def.Ident(),
			Type:  fmt.Sprintf("%T", def),
		}
	}
	return d
}

func Decode2(in io.Reader) (*meta.Module, error) {
	var hnd codec.CborHandle
	d := codec.NewDecoder(in, &hnd)
	ref := &Module{}
	err := d.Decode(ref)
	b := new(meta.Builder)
	m := b.Module(ref.Ident, nil)
	b.Description(m, ref.Description)
	for _, def := range ref.DataDef {
		switch def.Type {
		case "*meta.Leaf":
			b.Leaf(m, def.Ident)
		}
	}
	return m, err
}
