package emeta

import (
	"fmt"
	"io"

	"github.com/freeconf/yang/meta"
	"github.com/ugorji/go/codec"
)

func Encode(m *meta.Module, out io.Writer) error {
	var hnd codec.CborHandle
	hnd.EncodeOptions.StructToArray = true
	e := &encoder{}
	data := e.def(m)
	pack := codec.NewEncoder(out, &hnd)
	return pack.Encode(data)
}

type encoder struct {
}

func (e *encoder) defs(c meta.HasDataDefinitions) []interface{} {
	r := make([]interface{}, len(c.DataDefinitions()))
	for i, def := range c.DataDefinitions() {
		r[i] = e.def(def)
	}
	return r
}

func (e *encoder) def(m meta.Definition) any {
	switch x := m.(type) {
	case *meta.Module:
		return Module{
			x.Ident(),
			x.Description(),
			[]Extension{},
			DefTypeModule,
			x.Namespace(),
			x.Prefix(),
			x.Contact(),
			x.Organization(),
			x.Reference(),
			x.Version(),
			e.defs(x),
		}
	case *meta.Leaf:
		return Leaf{
			m.Ident(),
			x.Description(),
			[]Extension{},
			DefTypeLeaf,
			boolPtr(x.IsConfigSet(), x.Config()),
			boolPtr(x.IsMandatorySet(), x.Mandatory()),
		}
	}
	panic(fmt.Sprintf("not supported yet %T", m))
}

func boolPtr(isSet bool, v bool) *bool {
	if isSet {
		return &v
	}
	return nil
}
