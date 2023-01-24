package main

/*
#cgo CFLAGS: -Ic/include

#include <freeconf/selection.h>
*/
import "C"
import (
	"context"
	"errors"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

//export fc_upsert_from
func fc_upsert_from(c_sel C.fc_selection, c_node C.fc_node) *C.fc_node_error {
	sel, valid := pool.Get(int64(c_sel.node.pool_id)).(node.Selection)
	if !valid {
		return newNodeErr(errors.New("no selection found"))
	}
	next := sel.UpdateInto(CNode{c_sel, c_node})
	if next.LastErr != nil {
		return newNodeErr(next.LastErr)
	}
	return nil
}

type CNode struct {
	c_sel  C.fc_selection
	c_node C.fc_node
}

func (n CNode) Child(r node.ChildRequest) (child node.Node, err error) {
	var next C.fc_node
	meta := C.fc_find_meta(n.c_sel.meta, r.Meta.Ident())
	c_r := C.fc_child_request{
		context: n.c_node.context,
		meta:    meta,
	}
	c_err := n.c_node.on_child(n.c_node.context, c_r, &next)
	if c_err != nil {
		return nil, cerrToGo(c_err)
	}
	if next == nil {
		return nil, nil
	}
	c_path := C.fc_new_path(n.c_sel.path, meta)
	next_sel := C.fc_selection{c_path, next}
	return CNode{next_sel, next}, nil
}

func (n CNode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n CNode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	return nil
}

func (n CNode) Choose(sel node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error) {
	return nil, nil
}

func (n CNode) BeginEdit(r node.NodeRequest) error {
	return nil
}

func (n CNode) EndEdit(r node.NodeRequest) error {
	return nil
}

func (n CNode) Action(r node.ActionRequest) (output node.Node, err error) {
	return nil, nil
}

func (n CNode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	return nil, nil
}

func (n CNode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil

}

func (n CNode) Context(sel node.Selection) context.Context {
	return sel.Context
}
