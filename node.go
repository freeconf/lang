package main

/*
#cgo CFLAGS: -Iinclude

#include <string.h>
#include <freeconf/meta.h>

typedef struct fc_field_request {
    void* context;
	struct fc_selection* selection;
    fc_meta* meta;
    bool write;
} fc_field_request;

typedef struct fc_child_request {
    void* context;
	struct fc_selection* selection;
    fc_has_definitions* meta;
} fc_child_request;

typedef struct fc_val {
    int type;
//    fc_val_type type;
    void* data;
    size_t size;
} fc_val;

typedef char fc_node_error[128];

typedef struct fc_node {
    long pool_id;
    void* context;
    fc_node_error* (*on_field)(void* context, fc_field_request r, fc_val** val);
    fc_node_error* (*on_child)(void* context, fc_child_request r, struct fc_node* child);
} fc_node;

struct fc_path* fc_new_path(struct fc_path* parent, fc_meta* meta);

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
	err := (*C.fc_node_error)(C.malloc(unsafe.Sizeof(C.fc_node_error)))
	C.strncpy(err.message, msg, unsafe.Sizeof(err.message))
	return err
}

type gnode struct {
	csel  C.struct_fc_selection
	cnode C.fc_node
}

func (n gnode) Child(r node.ChildRequest) (child node.Node, err error) {
	var next C.fc_node
	meta := C.fc_find_meta(n.csel.meta, r.Meta.Ident())
	c_r := C.fc_child_request{
		context: n.cnode.context,
		meta:    meta,
	}
	c_err := n.cnode.on_child(n.cnode.context, c_r, &next)
	if c_err != nil {
		return nil, cerrToGo(c_err)
	}
	if next == nil {
		return nil, nil
	}
	c_path := C.fc_new_path(n.csel.path, meta)
	next_sel := C.struct_fc_selection{c_path, next}
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
