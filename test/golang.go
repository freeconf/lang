package test

import (
	"os"

	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

type golang struct {
}

func (d *golang) dump(sel node.Selection, fname string) error {
	out, err := os.Create(fname)
	if err != nil {
		return err
	}
	null := nodeutil.ReflectChild(make(map[string]any))
	return sel.UpsertInto(nodeutil.Dump(null, out)).LastErr
}
