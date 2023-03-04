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
	m := Handles.Get(in.ModuleHnd).(*meta.Module)
	n := Handles.Get(in.NodeHnd).(node.Node)
	b := node.NewBrowser(m, n)
	browserHnd := Handles.Put(b)
	return &pb.NewBrowserResponse{BrowserHnd: browserHnd}, nil
}

func (s *NodeService) BrowserRoot(ctx context.Context, in *pb.BrowserRootRequest) (*pb.BrowserRootResponse, error) {
	b := Handles.Get(in.BrowserHnd).(*node.Browser)
	root := b.Root()
	var resp pb.BrowserRootResponse
	resp.SelHnd = lazyGetSelectionHnd(&root)
	return &resp, nil
}

func (s NodeService) Find(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	sel := Handles.Get(in.SelHnd).(*node.Selection)
	found := sel.Find(in.Path)
	var resp pb.FindResponse
	if !found.IsNil() {
		resp.SelHnd = lazyGetSelectionHnd(&found)
	}
	return &resp, nil
}

func (s *NodeService) UpsertFrom(ctx context.Context, in *pb.UpsertFromRequest) (*pb.UpsertFromResponse, error) {
	sel := Handles.Get(in.SelHnd).(*node.Selection)
	n := s.resolveNode(in.NodeHnd)
	fmt.Printf("go: upsertFrom, from=%d\n", in.NodeHnd)
	err := sel.UpsertFrom(n).LastErr
	fmt.Printf("go: upsertFrom, err=%s\n", err)
	return &pb.UpsertFromResponse{}, err
}

func (s *NodeService) resolveNode(nodeHnd uint64) node.Node {
	if nodeHnd == 0 {
		return nil
	}
	n, found := Handles.Get(nodeHnd).(node.Node)
	if !found {
		n = &gnode{client: s.Client, nodeHnd: nodeHnd}
		Handles.Record(n, nodeHnd)
	}
	return n
}

func (s *NodeService) NewNode(ctx context.Context, in *pb.NewNodeRequest) (*pb.NewNodeResponse, error) {
	n := &gnode{client: s.Client}
	n.nodeHnd = Handles.Put(n)
	return &pb.NewNodeResponse{NodeHnd: n.nodeHnd}, nil
}

func (s *NodeService) Action(ctx context.Context, in *pb.ActionRequest) (*pb.ActionResponse, error) {
	inputNode := s.resolveNode(in.InputNodeHnd)
	sel := Handles.Require(in.SelHnd).(*node.Selection)
	output := sel.Action(inputNode)
	var resp pb.ActionResponse
	if !output.IsNil() {
		resp.OutputSelHnd = lazyGetSelectionHnd(&output)
	}
	return &resp, output.LastErr
}

func lazyGetSelectionHnd(sel *node.Selection) uint64 {
	if sel.Hnd != 0 {
		return sel.Hnd
	}
	sel.Hnd = Handles.Put(sel)
	return sel.Hnd
}

func (s *NodeService) GetSelection(ctx context.Context, in *pb.GetSelectionRequest) (*pb.GetSelectionResponse, error) {
	sel := Handles.Require(in.SelHnd).(*node.Selection)
	resp := pb.GetSelectionResponse{
		MetaIdent: sel.Path.Meta.Ident(),
		NodeHnd:   Handles.Hnd(sel.Node),
	}
	if sel.Parent == nil {
		resp.BrowserHnd = Handles.Hnd(sel.Browser)
	} else {
		resp.ParentHnd = lazyGetSelectionHnd(sel.Parent)
	}
	return &resp, nil
}

func (s *NodeService) GetBrowser(ctx context.Context, in *pb.GetBrowserRequest) (*pb.GetBrowserResponse, error) {
	browser := Handles.Require(in.BrowserHnd).(*node.Browser)
	resp := pb.GetBrowserResponse{
		ModuleHnd: Handles.Hnd(browser.Meta),
	}
	return &resp, nil
}

func (s *NodeService) GetModule(ctx context.Context, in *pb.GetModuleRequest) (*pb.GetModuleResponse, error) {
	m := Handles.Require(in.ModuleHnd).(*meta.Module)
	return &pb.GetModuleResponse{
		Module: new(MetaEncoder).Encode(m),
	}, nil
}

/**
 * gnode wraps the xlang's implementation of a node
 */
type gnode struct {
	client  pb.XNodeClient
	nodeHnd uint64
}

func (n *gnode) Context(s node.Selection) context.Context {
	req := pb.SelectRequest{
		MetaIdent: s.Meta().Ident(),
		NodeHnd:   n.nodeHnd,
	}
	if s.Parent != nil {
		req.ParentSelHnd = lazyGetSelectionHnd(s.Parent)
	} else {
		req.BrowserHnd = Handles.Hnd(s.Browser)
	}
	resp, err := n.client.Select(s.Context, &req)
	if err != nil {
		panic(err)
	}
	resp.SelHnd = Handles.Put(&s)
	return s.Context
}

func (n *gnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.ChildRequest{
		SelHnd:    lazyGetSelectionHnd(&r.Selection),
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
		SelHnd:    lazyGetSelectionHnd(&r.Selection),
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
		if hnd.Val != nil {
			fmt.Printf("read val, hnd=%d %s. %s\n", n.nodeHnd, r.Meta.Ident(), hnd.Val.String())
		} else {
			fmt.Printf("read val, hnd=%d %s.  nil\n", n.nodeHnd, r.Meta.Ident())
		}
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
	req := pb.XActionRequest{
		SelHnd:    lazyGetSelectionHnd(&r.Selection),
		MetaIdent: r.Meta.Ident(),
	}
	if !r.Input.IsNil() {
		req.InputSelHnd = r.Input.Hnd
	}
	resp, err := n.client.Action(r.Selection.Context, &req)
	if err != nil || resp.OutputNodeHnd == 0 {
		return nil, err
	}
	return Handles.Get(resp.OutputNodeHnd).(node.Node), nil
}

func (n *gnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	return nil, nil
}

func (n *gnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil
}
