package lang

import (
	"context"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

type NodeService struct {
	pb.UnimplementedNodeServer

	Client pb.XNodeClient
}

func (s *NodeService) NewBrowser(ctx context.Context, in *pb.NewBrowserRequest) (*pb.Handle, error) {
	m := Handles.Get(in.ModuleHandle).(*meta.Module)
	n := &xnode{nodeHandle: in.NodeHandle, client: s.Client}
	b := node.NewBrowser(m, n)
	return &pb.Handle{Handle: Handles.Put(b)}, nil
}

type xnode struct {
	nodeHandle uint64
	client     pb.XNodeClient
}

func (n *xnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.ChildRequest{
		Meta: r.Meta.Ident(),
	}
	resp, err := n.client.Child(r.Selection.Context, &req)
	if err != nil || resp.Handle == 0 {
		return nil, err
	}

	// c_sel, c_meta := n.new_select_and_meta(r.Meta)
	// defer C.free(unsafe.Pointer(c_sel))
	// c_r := C.fc_node_child_req{
	// 	selection: c_sel,
	// 	meta:      c_meta,
	// }
	// var next *C.fc_node
	// c_err := C.fc_select_child(c_r, &next)
	// if c_err != nil {
	// 	return nil, go_err(c_err)
	// }
	// if next == nil {
	// 	return nil, nil
	// }
	// child_path := C.fc_meta_path_new(n.c_path, c_meta)
	// return &xnode{c_node: next, c_path: child_path}, nil
	child := &xnode{nodeHandle: resp.Handle, client: n.client}
	return child, err
}

func (n *xnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

// func (n *xnode) new_select_and_meta(m meta.Identifiable) (*C.fc_select, *C.fc_meta) {
// 	ident := C.CString(m.Ident())
// 	defer C.free(unsafe.Pointer(ident))

// 	// replica of *potentially* what originally requested child
// 	c_sel := C.fc_select_new(0, n.c_node, n.c_path)
// 	c_meta := C.fc_meta_find(n.c_path.meta, ident)
// 	return c_sel, c_meta
// }

func (n *xnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	// c_sel, c_meta := n.new_select_and_meta(r.Meta)

	// defer C.fc_select_delete(c_sel)
	// c_r := C.fc_node_field_req{
	// 	selection: c_sel,
	// 	meta:      c_meta,
	// 	write:     C.bool(r.Write),
	// }
	// var c_val C.fc_val
	// if r.Write {
	// 	c_val = cee_val(hnd.Val)
	// }
	// c_err := C.fc_select_field(c_r, &c_val)
	// defer free_val(c_val)
	// if !r.Write {
	// 	hnd.Val = go_val(c_val)
	// }
	// if c_err != nil {
	// 	return go_err(c_err)
	// }
	return nil
}

func (n *xnode) Choose(sel node.Selection, choice *meta.Choice) (m *meta.ChoiceCase, err error) {
	return nil, nil
}

func (n *xnode) BeginEdit(r node.NodeRequest) error {
	return nil
}

func (n *xnode) EndEdit(r node.NodeRequest) error {
	return nil
}

func (n *xnode) Action(r node.ActionRequest) (output node.Node, err error) {
	return nil, nil
}

func (n *xnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	return nil, nil
}

func (n *xnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil
}

func (n *xnode) Context(sel node.Selection) context.Context {
	return sel.Context
}
