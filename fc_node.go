package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"context"
	"unsafe"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

type gnode struct {
	csel  *C.struct_fc_select
	cnode *C.fc_node
}

func (n gnode) Child(r node.ChildRequest) (child node.Node, err error) {
	var next *C.fc_node
	ident := C.CString(r.Meta.Ident())
	defer C.free(unsafe.Pointer(ident))
	meta := C.fc_meta_find(n.csel.path.meta, ident)
	c_r := C.fc_node_child_req{
		context: n.cnode.context,
		meta:    meta,
	}
	c_err := C.fc_node_child(n.cnode, n.cnode.context, c_r, &next)
	if c_err != nil {
		return nil, goErr(c_err)
	}
	if next == nil {
		return nil, nil
	}
	c_path := C.fc_meta_path_new(n.csel.path, meta)
	next_sel := &C.struct_fc_select{path: c_path, node: next}
	return gnode{next_sel, next}, nil
}

func (n gnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n gnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	return nil
}

func (n gnode) Choose(sel node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error) {
	return nil, nil
}

func (n gnode) BeginEdit(r node.NodeRequest) error {
	return nil
}

func (n gnode) EndEdit(r node.NodeRequest) error {
	return nil
}

func (n gnode) Action(r node.ActionRequest) (output node.Node, err error) {
	return nil, nil
}

func (n gnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	return nil, nil
}

func (n gnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil

}

func (n gnode) Context(sel node.Selection) context.Context {
	return sel.Context
}
