package lang

import (
	"context"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/val"
)

/**
 * xnode wraps the xlang's implementation of a node
 */
type xnode struct {
	d       *Driver
	nodeHnd uint64
}

func (n *xnode) Context(s node.Selection) context.Context {
	// TODO
	return s.Context
}

func (n *xnode) Child(r node.ChildRequest) (node.Node, error) {
	req := pb.XChildRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
		New:       r.New,
		Delete:    r.Delete,
	}
	resp, err := n.d.xnodes.XChild(r.Selection.Context, &req)
	if err != nil || resp.NodeHnd == 0 {
		return nil, err
	}
	return n.d.handles.Get(resp.NodeHnd).(node.Node), nil
}

func (n *xnode) Next(r node.ListRequest) (node.Node, []val.Value, error) {
	req := pb.XNextRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
		New:       r.New,
		Row:       r.Row64,
		First:     r.First,
		Delete:    r.Delete,
	}
	if len(r.Key) > 0 {
		req.Key = make([]*pb.Val, len(r.Key))
		for i, v := range r.Key {
			req.Key[i] = encodeVal(v)
		}
	}
	resp, err := n.d.xnodes.XNext(r.Selection.Context, &req)
	if err != nil || resp.NodeHnd == 0 {
		return nil, nil, err
	}
	var key []val.Value
	if len(resp.Key) > 0 {
		key = decodeVals(resp.Key)
	}
	return n.d.handles.Get(resp.NodeHnd).(node.Node), key, nil
}

func (n *xnode) Field(r node.FieldRequest, hnd *node.ValueHandle) error {
	req := pb.XFieldRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
		Write:     r.Write,
		Clear:     r.Clear,
	}
	if r.Write {
		req.ToWrite = encodeVal(hnd.Val)
	}
	resp, err := n.d.xnodes.XField(r.Selection.Context, &req)
	if err != nil {
		return err
	}
	if !r.Write {
		hnd.Val = decodeVal(resp.FromRead)
	}

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
	req := pb.XActionRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
	}
	if !r.Input.IsNil() {
		req.InputSelHnd = r.Input.Hnd
	}
	resp, err := n.d.xnodes.XAction(r.Selection.Context, &req)
	if err != nil || resp.OutputNodeHnd == 0 {
		return nil, err
	}
	return n.d.handles.Get(resp.OutputNodeHnd).(node.Node), nil
}

func (n *xnode) Notify(r node.NotifyRequest) (node.NotifyCloser, error) {
	req := pb.XNotificationRequest{
		SelHnd:    resolveSelection(n.d, &r.Selection),
		MetaIdent: r.Meta.Ident(),
	}
	var recvErr error
	client, err := n.d.xnodes.XNotification(r.Selection.Context, &req)
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

func (n *xnode) Peek(sel node.Selection, consumer interface{}) interface{} {
	return nil
}