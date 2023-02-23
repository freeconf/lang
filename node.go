package lang

import (
	"context"
	"fmt"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

type NodeService struct {
	pb.UnimplementedNodeServer
	Client pb.XNodeClient
}

func (s *NodeService) NewBrowser(ctx context.Context, in *pb.NewBrowserRequest) (*pb.NewBrowserResponse, error) {
	m := Handles.Get(in.GModuleHnd).(*meta.Module)
	n := &gnode{client: s.Client, xNodeHnd: in.XNodeHnd, xBrowserHnd: in.XBrowserHnd}
	b := node.NewBrowser(m, n)
	return &pb.NewBrowserResponse{GBrowserHnd: Handles.Put(b)}, nil
}

func (s *NodeService) BrowserRoot(ctx context.Context, in *pb.BrowserRootRequest) (*pb.BrowserRootResponse, error) {
	b := Handles.Get(in.GBrowserHnd).(*node.Browser)
	resp := pb.BrowserRootResponse{GSelHnd: Handles.Put(b.Root())}
	return &resp, nil
}

func (s *NodeService) UpsertFrom(ctx context.Context, in *pb.UpsertFromRequest) (*pb.UpsertFromResponse, error) {
	sel := Handles.Get(in.GSelHnd).(node.Selection)
	var err error
	var n node.Node
	if in.XNodeHnd > 0 {
		req := pb.SplitRequest{
			GSelHnd:     in.GSelHnd,
			XNodeHnd:    in.XNodeHnd,
			ModuleIdent: sel.Browser.Meta.Ident(),
			MetaPath:    sel.Path.StringNoModule(),
		}
		_, err := s.Client.Split(ctx, &req)
		if err != nil {
			return nil, err
		}
		n = &gnode{client: s.Client, xNodeHnd: in.XNodeHnd}
	} else if in.GNodeHnd > 0 {
		n = Handles.Get(in.GNodeHnd).(node.Node)
	}
	err = sel.UpdateFrom(n).LastErr
	return &pb.UpsertFromResponse{}, err
}

/**
 * gnode wraps the xlang's implementation of a node
 */
type gnode struct {
	client pb.XNodeClient
	//xSelHnd uint64
	xNodeHnd    uint64
	xBrowserHnd uint64
}

type xSelHndKey int

const xSelHndContextKey = xSelHndKey(0)

func (n *gnode) Context(s node.Selection) context.Context {
	req := pb.SelectRequest{
		GSelHnd:   Handles.Put(s),
		MetaIdent: s.Meta().Ident(),
		XNodeHnd:  n.xNodeHnd,
	}
	if s.Parent == nil {
		req.XBrowserHnd = n.xBrowserHnd
	} else {
		req.XSelHnd = n.xSelHnd(s.Context)
	}
	fmt.Printf("adding client=%v, x_sel_hnd=%d, x_browser_hnd=%d\n", n.client, req.XSelHnd, req.XBrowserHnd)
	resp, err := n.client.Select(s.Context, &req)
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(s.Context, xSelHndContextKey, resp.XSelHnd)
	return ctx
}

func (n *gnode) xSelHnd(ctx context.Context) uint64 {
	return ctx.Value(xSelHndContextKey).(uint64)
}

func (n *gnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.ChildRequest{
		XSelHnd:   n.xSelHnd(r.Selection.Context),
		MetaIdent: r.Meta.Ident(),
		New:       r.New,
		Delete:    r.Delete,
	}
	resp, err := n.client.Child(r.Selection.Context, &req)
	if err != nil || resp.XNodeHnd == 0 {
		return nil, err
	}
	return &gnode{client: n.client, xNodeHnd: resp.XNodeHnd, xBrowserHnd: n.xBrowserHnd}, err
}

func (n *gnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n *gnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
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
