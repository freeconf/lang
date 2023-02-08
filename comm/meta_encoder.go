package comm

import (
	"github.com/freeconf/lang/comm/pb"
	"github.com/freeconf/yang/meta"
)

type metaEncoder struct {
}

func (e *metaEncoder) encode(from *meta.Module) *pb.Module {
	return &pb.Module{
		Ident:       from.Ident(),
		Description: from.Description(),
	}
}
