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

func (n *gnode) Child(r node.ChildRequest) (child node.Node, err error) {
	fmt.Printf("child, self.c_node=%p\n", n.c_node)
	c_sel, c_meta := n.new_select_and_meta(r.Meta)
	defer C.free(unsafe.Pointer(c_sel))
	c_r := C.fc_node_child_req{
		selection: c_sel,
		meta:      c_meta,
	}
	var next *C.fc_node
	c_err := C.fc_select_child(c_r, &next)
	if c_err != nil {
		return nil, go_err(c_err)
	}
	if next == nil {
		return nil, nil
	}
	child_path := C.fc_meta_path_new(n.c_path, c_meta)
	return &gnode{c_node: next, c_path: child_path}, nil
}

func (n *gnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n *gnode) new_select_and_meta(m meta.Identifiable) (*C.fc_select, *C.fc_meta) {
	ident := C.CString(m.Ident())
	defer C.free(unsafe.Pointer(ident))

	// replica of *potentially* what originally requested child
	c_sel := C.fc_select_new(0, n.c_node, n.c_path)
	c_meta := C.fc_meta_find(n.c_path.meta, ident)
	return c_sel, c_meta
}

func (n *gnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	fmt.Printf("field, self.c_node=%p\n", n.c_node)
	c_sel, c_meta := n.new_select_and_meta(r.Meta)

	fmt.Printf("field, c_sel.node=%p\n", c_sel.node)
	defer C.free(unsafe.Pointer(c_sel))
	c_r := C.fc_node_field_req{
		selection: c_sel,
		meta:      c_meta,
		write:     C.bool(r.Write),
	}
	var c_val C.fc_val
	if r.Write {
		c_val = cee_val(hnd.Val)
	}
	// HACK: Find out why this is nec
	c_r.selection.node = n.c_node

	fmt.Printf("field, c_r.selection.node=%p\n", c_r.selection.node)
	c_err := C.fc_select_field(c_r, &c_val)
	defer free_val(c_val)
	if !r.Write {
		hnd.Val = go_val(c_val)
	}
	if c_err != nil {
		return go_err(c_err)
	}
	return nil
}

func (n *gnode) Choose(sel node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error) {
	return nil, nil
}

func (n *gnode) BeginEdit(r node.NodeRequest) error {
	return nil
}

func (n *gnode) EndEdit(r node.NodeRequest) error {
	return nil
}

func (n *gnode) Action(r node.ActionRequest) (output node.Node, err error) {
	return nil, nil
}

func (n *gnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	return nil, nil
}

func (n *gnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil
}

func (n *gnode) Context(sel node.Selection) context.Context {
	return sel.Context
}
