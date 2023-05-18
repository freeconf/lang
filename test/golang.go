package test

import (
	"container/list"
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

type golang struct {
	traceFile *os.File
	recievers *list.List
}

func newGolang() *golang {
	return &golang{
		recievers: list.New(),
	}
}

func (d *golang) Close() error {
	return nil
}

func (d *golang) Connect() error {
	return nil
}

func (d *golang) createTestCase(c pb.TestCase, tracefile string) (node.Node, error) {
	var n node.Node
	switch c {
	case pb.TestCase_BASIC:
		n = nodeutil.ReflectChild(make(map[string]any))
	case pb.TestCase_ECHO:
		n = d.echoNode()
	default:
		panic("test case not implemented")
	}
	f, err := os.Create(tracefile)
	if err != nil {
		return nil, err
	}
	d.traceFile = f
	return nodeutil.Trace(n, f), nil
}

func (d *golang) echoNode() node.Node {
	return &nodeutil.Basic{
		OnAction: func(r node.ActionRequest) (node.Node, error) {
			switch r.Meta.Ident() {
			case "echo":
				data, err := nodeutil.WriteJSON(r.Input)
				if err != nil {
					return nil, err
				}
				return nodeutil.ReadJSON(data), nil
			case "send":
				d.send(r.Input.Node)
			}
			return nil, nil
		},
		OnNotify: func(r node.NotifyRequest) (node.NotifyCloser, error) {
			switch r.Meta.Ident() {
			case "recv":
				sub := d.subscribe(func(msg node.Node) {
					r.Send(msg)
				})
				return sub.Close, nil
			}
			return nil, nil
		},
	}
}

type reciever func(msg node.Node)

func (d *golang) send(msg node.Node) {
	p := d.recievers.Front()
	for p != nil {
		p.Value.(reciever)(msg)
		p = p.Next()
	}
}

func (d *golang) subscribe(r reciever) nodeutil.Subscription {
	return nodeutil.NewSubscription(d.recievers, d.recievers.PushBack(r))
}

func (d *golang) finalizeTestCase() error {
	return d.traceFile.Close()
}
