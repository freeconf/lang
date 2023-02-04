package codegen

import (
	"fmt"
	"io"

	"github.com/freeconf/yang/meta"
	"github.com/ugorji/go/codec"
)

func Encode(m *meta.Module, memId int64, out io.Writer) error {
	var hnd codec.CborHandle
	hnd.EncodeOptions.StructToArray = true
	e := &encoder{memId: memId}
	data := e.def(m)
	pack := codec.NewEncoder(out, &hnd)
	return pack.Encode(data)
}

type encoder struct {
	memId int64
}

func (e *encoder) defs(c meta.HasDataDefinitions) []DataDefinitions {
	r := make([]DataDefinitions, len(c.DataDefinitions()))
	for i, def := range c.DataDefinitions() {
		r[i] = e.def(def)
	}
	return r
}

func (e *encoder) extensions(m meta.HasExtensions) []Extension {
	src := m.Extensions()
	dest := make([]Extension, len(src))
	for i, x := range src {
		dest[i] = Extension{
			MetaIdExtension,
			x.Ident(),
			x.Prefix(),
			x.Keyword(),
			x.Def().Ident(),
			x.Arguments(),
		}
	}
	return dest
}

func (e *encoder) def(m meta.Definition) any {
	switch x := m.(type) {
	case *meta.Module:
		return Module{
			MetaIdModule,
			x.Ident(),
			x.Description(),
			e.extensions(m),
			e.defs(x),
			e.memId,
			x.Namespace(),
			x.Prefix(),
			x.Contact(),
			x.Organization(),
			x.Reference(),
			x.Version(),
		}
	case *meta.List:
		return List{
			MetaIdList,
			x.Ident(),
			x.Description(),
			e.extensions(m),
			e.defs(x),
			boolPtr(x.IsConfigSet(), x.Config()),
			boolPtr(x.IsMandatorySet(), x.Mandatory()),
		}
	case *meta.Container:
		return Container{
			MetaIdContainer,
			x.Ident(),
			x.Description(),
			e.extensions(m),
			e.defs(x),
			boolPtr(x.IsConfigSet(), x.Config()),
			boolPtr(x.IsMandatorySet(), x.Mandatory()),
		}
	case *meta.LeafList:
		return LeafList{
			MetaIdLeafList,
			m.Ident(),
			x.Description(),
			e.extensions(m),
			boolPtr(x.IsConfigSet(), x.Config()),
			boolPtr(x.IsMandatorySet(), x.Mandatory()),
		}
	case *meta.Leaf:
		return Leaf{
			MetaIdLeaf,
			m.Ident(),
			x.Description(),
			e.extensions(m),
			boolPtr(x.IsConfigSet(), x.Config()),
			boolPtr(x.IsMandatorySet(), x.Mandatory()),
		}
	}
	panic(fmt.Sprintf("not supported yet %T", m))
}

func boolPtr(isSet bool, v bool) OptionalBool {
	if isSet {
		if v {
			return BoolTrue
		}
		return BoolFalse
	}
	return BoolNotSet
}
