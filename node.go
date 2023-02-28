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

func (s *NodeService) NewBrowser(ctx context.Context, in *pb.NewBrowserRequest) (*pb.NewBrowserResponse, error) {
	m := Handles.Get(in.ModuleHnd).(*meta.Module)
	n := Handles.Get(in.NodeHnd).(node.Node)
	b := node.NewBrowser(m, n)
	browserHnd := Handles.Put(b)
	return &pb.NewBrowserResponse{BrowserHnd: browserHnd}, nil
}

func (s *NodeService) BrowserRoot(ctx context.Context, in *pb.BrowserRootRequest) (*pb.BrowserRootResponse, error) {
	b := Handles.Get(in.BrowserHnd).(*node.Browser)
	root := b.Root()
	selHndObj := root.Context.Value(selHndContextKey)
	var resp pb.BrowserRootResponse
	if selHndObj != nil {
		resp.SelHnd = selHndObj.(uint64)
	} else {
		resp.SelHnd = Handles.Put(root)
	}
	return &resp, nil
}

func (s *NodeService) UpsertFrom(ctx context.Context, in *pb.UpsertFromRequest) (*pb.UpsertFromResponse, error) {
	sel := Handles.Get(in.SelHnd).(node.Selection)
	var err error
	var n node.Node
	n, found := Handles.Get(in.NodeHnd).(node.Node)
	if !found {
		n = &gnode{client: s.Client, nodeHnd: in.NodeHnd}
		Handles.Record(n, in.NodeHnd)
	}
	// req := pb.SplitRequest{
	// 	NodeHnd:   in.NodeHnd,
	// 	ModuleHnd: Handles.Hnd(sel.Browser.Meta),
	// 	MetaPath:  sel.Path.StringNoModule(),
	// }
	// resp, err := s.Client.Split(ctx, &req)
	// split := sel.Split(n)
	// Handles.Record(split, resp.SelHnd)
	// if err != nil {
	// 	return nil, err
	// }
	err = sel.UpdateFrom(n).LastErr
	return &pb.UpsertFromResponse{}, err
}

func (s *NodeService) NewNode(ctx context.Context, in *pb.NewNodeRequest) (*pb.NewNodeResponse, error) {
	n := &gnode{client: s.Client}
	n.nodeHnd = Handles.Put(n)
	return &pb.NewNodeResponse{NodeHnd: n.nodeHnd}, nil
}

/**
 * gnode wraps the xlang's implementation of a node
 */
type gnode struct {
	client  pb.XNodeClient
	nodeHnd uint64
}

type selHndKey int

const selHndContextKey = selHndKey(0)

func (n *gnode) Context(s node.Selection) context.Context {
	req := pb.SelectRequest{
		MetaIdent: s.Meta().Ident(),
		NodeHnd:   n.nodeHnd,
	}
	if s.Parent != nil {
		req.ParentSelHnd = n.selHnd(s.Context)
	} else {
		req.BrowserHnd = Handles.Hnd(s.Browser)
	}
	resp, err := n.client.Select(s.Context, &req)
	if err != nil {
		panic(err)
	}
	s.Context = context.WithValue(s.Context, selHndContextKey, resp.SelHnd)
	Handles.Record(s, resp.SelHnd)
	return s.Context
}

func (n *gnode) selHnd(ctx context.Context) uint64 {
	return ctx.Value(selHndContextKey).(uint64)
}

func (n *gnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.ChildRequest{
		SelHnd:    n.selHnd(r.Selection.Context),
		MetaIdent: r.Meta.Ident(),
		New:       r.New,
		Delete:    r.Delete,
	}
	resp, err := n.client.Child(r.Selection.Context, &req)
	if err != nil || resp.NodeHnd == 0 {
		return nil, err
	}
	return Handles.Get(resp.NodeHnd).(node.Node), nil
}

func (n *gnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n *gnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	req := pb.FieldRequest{
		SelHnd:    n.selHnd(r.Selection.Context),
		MetaIdent: r.Meta.Ident(),
		Write:     r.Write,
		Clear:     r.Clear,
	}
	if r.Write {
		req.ToWrite = encodeVal(hnd.Val)
	}
	resp, err := n.client.Field(r.Selection.Context, &req)
	if err != nil {
		return err
	}
	if !r.Write {
		hnd.Val = decodeVal(resp.FromRead)
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
