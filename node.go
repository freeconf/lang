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

//export fc_new_node_err
func fc_new_node_err(msg *C.char) *C.fc_node_error {
	var err *C.fc_node_error
	err = (*C.fc_node_error)(C.malloc(C.size_t(unsafe.Sizeof(*err))))
	C.strcpy(&err.message[0], msg)
	return err
}

type gnode struct {
	csel  *C.struct_fc_selection
	cnode *C.fc_node
}

func (n gnode) Child(r node.ChildRequest) (child node.Node, err error) {
	var next *C.fc_node
	ident := C.CString(r.Meta.Ident())
	defer C.free(unsafe.Pointer(ident))
	meta := C.fc_find_meta(n.csel.path.meta, ident)
	c_r := C.fc_child_request{
		context: n.cnode.context,
		meta:    meta,
	}
	c_err := C.fc_node_on_child_x(n.cnode, n.cnode.context, c_r, &next)
	if c_err != nil {
		return nil, cerrToGo(c_err)
	}
	if next == nil {
		return nil, nil
	}
	c_path := fc_new_path(n.csel.path, meta)
	next_sel := &C.struct_fc_selection{path: c_path, node: next}
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
