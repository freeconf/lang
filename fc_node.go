package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"context"
	"fmt"
	"unsafe"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

type gnode struct {
	c_node *C.struct_fc_node
	c_path *C.struct_fc_meta_path
}

func (n gnode) Child(r node.ChildRequest) (child node.Node, err error) {
	var next *C.fc_node
	ident := C.CString(r.Meta.Ident())
	defer C.free(unsafe.Pointer(ident))

	// replica of *potentially* what originally requested child
	c_sel := C.fc_select_new(0, n.c_node, n.c_path)
	defer C.free(unsafe.Pointer(c_sel))

	c_meta := C.fc_meta_find(n.c_path.meta, ident)

	c_r := C.fc_node_child_req{
		selection: c_sel,
		meta:      c_meta,
	}
	c_err := C.fc_select_child(c_r, &next)
	if c_err != nil {
		return nil, go_err(c_err)
	}
	if next == nil {
		return nil, nil
	}
	fmt.Printf("!! here %s\n", r.Meta.Ident())
	child_path := C.fc_meta_path_new(n.c_path, c_meta)
	return gnode{c_node: next, c_path: child_path}, nil
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
