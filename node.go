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
	d *Driver
}

func (s *NodeService) NewBrowser(ctx context.Context, in *pb.NewBrowserRequest) (*pb.NewBrowserResponse, error) {
	m := s.d.handles.Get(in.ModuleHnd).(*meta.Module)
	browserHnd := s.d.handles.Reserve()
	nodeSrc := func() node.Node {
		req := pb.NodeSourceRequest{BrowserHnd: browserHnd}
		resp, err := s.d.xnodes.NodeSource(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		return s.d.handles.Get(resp.NodeHnd).(node.Node)
	}
	b := node.NewBrowserSource(m, nodeSrc)
	s.d.handles.Record(b, browserHnd)
	return &pb.NewBrowserResponse{BrowserHnd: browserHnd}, nil
}

func (s *NodeService) BrowserRoot(ctx context.Context, in *pb.BrowserRootRequest) (*pb.BrowserRootResponse, error) {
	b := s.d.handles.Get(in.BrowserHnd).(*node.Browser)
	root := b.Root()
	var resp pb.BrowserRootResponse
	resp.SelHnd = resolveSelection(s.d, &root)
	return &resp, nil
}

func (s NodeService) Find(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	sel := s.d.handles.Get(in.SelHnd).(*node.Selection)
	found := sel.Find(in.Path)
	var resp pb.FindResponse
	if !found.IsNil() {
		resp.SelHnd = resolveSelection(s.d, &found)
	}
	return &resp, nil
}

func (s *NodeService) UpsertFrom(ctx context.Context, in *pb.UpsertFromRequest) (*pb.UpsertFromResponse, error) {
	sel := s.d.handles.Get(in.SelHnd).(*node.Selection)
	n := s.d.handles.Get(in.NodeHnd).(node.Node)
	err := sel.UpsertFrom(n).LastErr
	return &pb.UpsertFromResponse{}, err
}

func (s *NodeService) NewNode(ctx context.Context, in *pb.NewNodeRequest) (*pb.NewNodeResponse, error) {
	n := &gnode{d: s.d}
	n.nodeHnd = s.d.handles.Put(n)
	return &pb.NewNodeResponse{NodeHnd: n.nodeHnd}, nil
}

func (s *NodeService) Action(ctx context.Context, in *pb.ActionRequest) (*pb.ActionResponse, error) {
	var input node.Node
	if in.InputNodeHnd != 0 {
		input = s.d.handles.Get(in.InputNodeHnd).(node.Node)
	}
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	output := sel.Action(input)
	var resp pb.ActionResponse
	if !output.IsNil() {
		resp.OutputSelHnd = resolveSelection(s.d, &output)
	}
	return &resp, output.LastErr
}

func resolveSelection(d *Driver, sel *node.Selection) uint64 {
	if sel.Hnd != 0 {
		return sel.Hnd
	}
	sel.Hnd = d.handles.Put(sel)
	return sel.Hnd
}

func (s *NodeService) GetSelection(ctx context.Context, in *pb.GetSelectionRequest) (*pb.GetSelectionResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	resp := pb.GetSelectionResponse{
		MetaIdent: sel.Path.Meta.Ident(),
		NodeHnd:   s.d.handles.Hnd(sel.Node),
	}
	if sel.Parent == nil {
		resp.BrowserHnd = s.d.handles.Hnd(sel.Browser)
	} else {
		resp.ParentHnd = resolveSelection(s.d, sel.Parent)
	}
	return &resp, nil
}

func (s *NodeService) GetBrowser(ctx context.Context, in *pb.GetBrowserRequest) (*pb.GetBrowserResponse, error) {
	browser := s.d.handles.Require(in.BrowserHnd).(*node.Browser)
	resp := pb.GetBrowserResponse{
		ModuleHnd: s.d.handles.Hnd(browser.Meta),
	}
	return &resp, nil
}

func (s *NodeService) GetModule(ctx context.Context, in *pb.GetModuleRequest) (*pb.GetModuleResponse, error) {
	m := s.d.handles.Require(in.ModuleHnd).(*meta.Module)
	return &pb.GetModuleResponse{
		Module: new(MetaEncoder).Encode(m),
	}, nil
}

func (s *NodeService) Notification(in *pb.NotificationRequest, srv pb.Node_NotificationServer) error {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	closer, err := sel.Notifications(func(n node.Notification) {
		resp := pb.NotificationResponse{
			SelHnd: resolveSelection(s.d, &n.Event),
		}
		if err := srv.Send(&resp); err != nil {
			panic(err)
		}
	})
	if err != nil {
		return err
	}
	<-srv.Context().Done()
	closer()
	return nil
}

/**
 * gnode wraps the xlang's implementation of a node
 */
type gnode struct {
	d       *Driver
	nodeHnd uint64
}

func (n *gnode) Context(s node.Selection) context.Context {
	// TODO
	return s.Context
}

func (n *gnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.ChildRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
		New:       r.New,
		Delete:    r.Delete,
	}
	resp, err := n.d.xnodes.Child(r.Selection.Context, &req)
	if err != nil || resp.NodeHnd == 0 {
		return nil, err
	}
	return n.d.handles.Get(resp.NodeHnd).(node.Node), nil
}

func (n *gnode) Next(r node.ListRequest) (next node.Node, key []val.Value, err error) {
	return nil, nil, nil
}

func (n *gnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	req := pb.FieldRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
		Write:     r.Write,
		Clear:     r.Clear,
	}
	if r.Write {
		req.ToWrite = encodeVal(hnd.Val)
	}
	resp, err := n.d.xnodes.Field(r.Selection.Context, &req)
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
	req := pb.XActionRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
	}
	if !r.Input.IsNil() {
		req.InputSelHnd = r.Input.Hnd
	}
	resp, err := n.d.xnodes.Action(r.Selection.Context, &req)
	if err != nil || resp.OutputNodeHnd == 0 {
		return nil, err
	}
	return n.d.handles.Get(resp.OutputNodeHnd).(node.Node), nil
}

func (n *gnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	req := pb.XNotificationRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
	}
	var recvErr error
	client, err := n.d.xnodes.Notification(r.Selection.Context, &req)
	if err != nil {
		return nil, err
	}
	closer := func() error {
		if client != nil {
			if err := client.CloseSend(); err != nil && recvErr == nil {
				return err
			}
			client = nil
		}
		return recvErr
	}
	go func() {
		var resp *pb.XNotificationResponse
		for {
			resp, recvErr = client.Recv()
			if resp == nil || recvErr != nil {
				break
			}
			n := n.d.handles.Get(resp.NodeHnd).(node.Node)
			r.Send(n)
		}
	}()
	return closer, nil
}

func (n *gnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil
}
