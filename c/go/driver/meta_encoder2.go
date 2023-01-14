package driver

import (
	"io"

	"github.com/freeconf/yang/meta"
	"github.com/ugorji/go/codec"
)

func Encode2(m *meta.Module, out io.Writer) error {
	var hnd codec.CborHandle
	hnd.EncodeOptions.StructToArray = true
	e := &encoder{}
	data := e.module(m)
	pack := codec.NewEncoder(out, &hnd)
	return pack.Encode(data)
}

type encoder struct {
}

type Module struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     int
	Definitions []interface{}
}

type DefType int

const (
	DefTypeModule = iota
	DefTypeContainer
	DefTypeLeaf
)

type Extension struct {
	Name string
}

type Leaf struct {
	Ident       string
	Description string
	Extensions  []Extension
	DefType     int
}

func (e *encoder) module(m *meta.Module) Module {
	return Module{
		m.Ident(),
		m.Description(),
		[]Extension{},
		DefTypeModule,
		e.defs(m),
	}
}

func (e *encoder) defs(c meta.HasDataDefinitions) []interface{} {
	r := make([]interface{}, len(c.DataDefinitions()))
	for i, def := range c.DataDefinitions() {
		r[i] = e.def(def)
	}
	return r
}

func (e *encoder) description(i interface{}) string {
	if d, valid := i.(meta.Describable); valid {
		return d.Description()
	}
	return ""
}

func (e *encoder) def(m meta.Definition) Leaf {
	return Leaf{
		m.Ident(),
		e.description(m),
		[]Extension{},
		DefTypeLeaf,
	}
}
