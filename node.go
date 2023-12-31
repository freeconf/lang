package lang

import (
	"context"
	"fmt"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
)

type NodeService struct {
	pb.UnimplementedNodeServer
	d *Driver
}

func (s *NodeService) NewBrowser(ctx context.Context, in *pb.NewBrowserRequest) (*pb.NewBrowserResponse, error) {
	m := s.d.handles.Require(in.ModuleHnd).(*meta.Module)
	browserHnd := s.d.handles.Reserve()
	nodeSrc := func() node.Node {
		req := pb.XNodeSourceRequest{BrowserHnd: browserHnd}
		resp, err := s.d.xnodes.XNodeSource(context.Background(), &req)
		if err != nil {
			panic(err)
		}
		return s.d.handles.Require(resp.NodeHnd).(node.Node)
	}
	b := node.NewBrowserSource(m, nodeSrc)
	s.d.handles.Record(b, browserHnd)
	return &pb.NewBrowserResponse{BrowserHnd: browserHnd}, nil
}

func (s *NodeService) BrowserRoot(ctx context.Context, in *pb.BrowserRootRequest) (*pb.BrowserRootResponse, error) {
	b := s.d.handles.Require(in.BrowserHnd).(*node.Browser)
	root := b.Root()
	var resp pb.BrowserRootResponse
	resp.SelHnd = resolveSelection(s.d, root)
	return &resp, nil
}

func (s NodeService) ReleaseSelection(ctx context.Context, in *pb.ReleaseSelectionRequest) (*pb.ReleaseSelectionResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	sel.Release()
	return &pb.ReleaseSelectionResponse{}, nil
}

func (s NodeService) Find(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	found, err := sel.Find(in.Path)
	if err != nil {
		return nil, err
	}
	var resp pb.FindResponse
	if found != nil {
		resp.SelHnd = resolveSelection(s.d, found)
	}
	return &resp, err
}

func (s *NodeService) SelectionEdit(ctx context.Context, in *pb.SelectionEditRequest) (*pb.SelectionEditResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	var n node.Node
	if in.NodeHnd != 0 {
		n = s.d.handles.Require(in.NodeHnd).(node.Node)
	}
	var err error
	switch in.Op {
	case pb.SelectionEditOp_UPSERT_INTO:
		err = sel.UpsertInto(n)
	case pb.SelectionEditOp_UPSERT_FROM:
		err = sel.UpsertFrom(n)
	case pb.SelectionEditOp_INSERT_INTO:
		err = sel.InsertInto(n)
	case pb.SelectionEditOp_INSERT_FROM:
		err = sel.InsertFrom(n)
	case pb.SelectionEditOp_UPSERT_INTO_SET_DEFAULTS:
		err = sel.UpsertIntoSetDefaults(n)
	case pb.SelectionEditOp_UPSERT_FROM_SET_DEFAULTS:
		err = sel.UpsertFromSetDefaults(n)
	case pb.SelectionEditOp_UPDATE_INTO:
		err = sel.UpdateInto(n)
	case pb.SelectionEditOp_UPDATE_FROM:
		err = sel.UpdateFrom(n)
	case pb.SelectionEditOp_REPLACE_FROM:
		err = sel.ReplaceFrom(n)
	case pb.SelectionEditOp_DELETE:
		err = sel.Delete()
		if err != nil {
			sel.Release()
		}
	}
	return &pb.SelectionEditResponse{}, err
}

func (s *NodeService) NewNode(ctx context.Context, in *pb.NewNodeRequest) (*pb.NewNodeResponse, error) {
	n := &xnode{d: s.d}
	n.nodeHnd = s.d.handles.Put(n)
	return &pb.NewNodeResponse{NodeHnd: n.nodeHnd}, nil
}

func (s *NodeService) Action(ctx context.Context, in *pb.ActionRequest) (*pb.ActionResponse, error) {
	var input node.Node
	if in.InputNodeHnd != 0 {
		input = s.d.handles.Require(in.InputNodeHnd).(node.Node)
	}
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	output, err := sel.Action(input)
	if err != nil {
		return nil, err
	}
	var resp pb.ActionResponse
	if output != nil {
		resp.OutputSelHnd = resolveSelection(s.d, output)
	}
	return &resp, nil
}

func resolveSelection(d *Driver, sel *node.Selection) uint64 {
	hnd := d.handles.Hnd(sel)
	return hnd
}

func (s *NodeService) GetSelection(ctx context.Context, in *pb.GetSelectionRequest) (*pb.GetSelectionResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	_, isNotRemote := sel.Node.(*xnode) // here "remote" is from perspective of X impl
	resp := pb.GetSelectionResponse{
		NodeHnd:    s.d.handles.Hnd(sel.Node),
		RemoteNode: !isNotRemote,
	}
	s.d.handles.Hnd(sel.Browser.Meta)
	resp.Path = buildPath(s.d, sel.Path)
	if sel.InsideList {
		resp.InsideList = true
	}
	resp.BrowserHnd = s.d.handles.Hnd(sel.Browser)
	return &resp, nil
}

func buildPath(d *Driver, p *node.Path) *pb.Path {
	// there might be an opportunity to optimize this given that caller likely
	// has most of the path from previous calls and we're sending full path copy back
	// each time.  Maybe caller can send hint as to what piece they are
	// missing and we can send just those pieces.
	protoSegs := make([]*pb.PathSegment, len(p.Segments())-1)
	root := p.Segments()[0].Meta
	if _, isMod := root.(*meta.Module); !isMod {
		panic(fmt.Sprintf("path %s does not begin with a module", p.String()))
	}
	protoPath := &pb.Path{
		Segments:  protoSegs,
		ModuleHnd: d.handles.Hnd(root),
	}
	for i, seg := range p.Segments()[1:] {
		protoSegs[i] = &pb.PathSegment{
			MetaIdent: seg.Meta.Ident(),
		}
		if meta.IsAction(seg.Meta) {
			protoSegs[i].Type = pb.PathSegmentType_RPC
		} else if meta.IsNotification(seg.Meta) {
			protoSegs[i].Type = pb.PathSegmentType_NOTIFICATION
		} else if _, match := seg.Meta.(*meta.RpcInput); match {
			protoSegs[i].Type = pb.PathSegmentType_RPC_INPUT
		} else if _, match := seg.Meta.(*meta.RpcOutput); match {
			protoSegs[i].Type = pb.PathSegmentType_RPC_OUTPUT
		} else {
			protoSegs[i].Type = pb.PathSegmentType_DATA_DEF
		}
		if seg.Key != nil {
			protoSegs[i].Key = encodeVals(seg.Key)
		}
	}
	return protoPath
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
		Module: NewMetaEncoder().Encode(m),
	}, nil
}

func (s *NodeService) Notification(in *pb.NotificationRequest, srv pb.Node_NotificationServer) error {
	closed := false
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	closer, err := sel.Notifications(func(n node.Notification) {
		resp := pb.NotificationResponse{
			SelHnd: resolveSelection(s.d, n.Event),
			When:   n.EventTime.Unix(),
		}
		if err := srv.Send(&resp); err != nil && !closed {
			panic(err)
		}
	})
	if err != nil {
		return err
	}
	<-srv.Context().Done()
	closed = true
	closer()
	return nil
}
