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
	m := s.d.handles.Get(in.ModuleHnd).(*meta.Module)
	browserHnd := s.d.handles.Reserve()
	nodeSrc := func() node.Node {
		req := pb.XNodeSourceRequest{BrowserHnd: browserHnd}
		resp, err := s.d.xnodes.XNodeSource(context.Background(), &req)
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
	resp.SelHnd = resolveSelection(s.d, root)
	return &resp, nil
}

func (s NodeService) Find(ctx context.Context, in *pb.FindRequest) (*pb.FindResponse, error) {
	sel := s.d.handles.Get(in.SelHnd).(*node.Selection)
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
	sel := s.d.handles.Get(in.SelHnd).(*node.Selection)
	n := s.d.handles.Get(in.NodeHnd).(node.Node)
	var err error
	switch in.Op {
	case pb.SelectionEditOp_UPSERT_INTO:
		err = sel.UpsertInto(n)
	case pb.SelectionEditOp_UPSERT_FROM:
		err = sel.UpsertFrom(n)
	case pb.SelectionEditOp_INSERT_INTO:
		fmt.Printf("go.insert_info[%d], path=%s\n", in.SelHnd, sel.Path.String())
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
		input = s.d.handles.Get(in.InputNodeHnd).(node.Node)
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
	return d.handles.Hnd(sel)
}

func (s *NodeService) GetSelection(ctx context.Context, in *pb.GetSelectionRequest) (*pb.GetSelectionResponse, error) {
	sel := s.d.handles.Require(in.SelHnd).(*node.Selection)
	resp := pb.GetSelectionResponse{
		NodeHnd: s.d.handles.Hnd(sel.Node),
	}
	resp.Path = &pb.PathSegment{
		MetaIdent: sel.Path.Meta.Ident(),
	}
	if sel.InsideList {
		resp.Path.Key = encodeVals(sel.Path.Key)
		resp.Path.Type = pb.PathSegmentType_LIST_ITEM
		resp.InsideList = true
	} else if meta.IsAction(sel.Path.Meta) {
		resp.Path.Type = pb.PathSegmentType_RPC
	} else if meta.IsNotification(sel.Path.Meta) {
		resp.Path.Type = pb.PathSegmentType_RPC
	} else if _, isRpcInput := sel.Path.Meta.(*meta.RpcInput); isRpcInput {
		resp.Path.Type = pb.PathSegmentType_RPC_INPUT
	} else if _, isRpcOutput := sel.Path.Meta.(*meta.RpcOutput); isRpcOutput {
		resp.Path.Type = pb.PathSegmentType_RPC_OUTPUT
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
			SelHnd: resolveSelection(s.d, n.Event),
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
