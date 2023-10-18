package test

import (
	"container/list"
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
	"github.com/freeconf/yang/val"
)

type golang struct {
	ypath      source.Opener
	yangModule *meta.Module
	traceFile  *os.File
	recievers  *list.List
}

func newGolang(ypath source.Opener) *golang {
	return &golang{
		ypath:     ypath,
		recievers: list.New(),
	}
}

func (d *golang) Name() string {
	return "go"
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
	case pb.TestCase_ADVANCED:
		n = d.advancedNode(make(map[string]interface{}))
	case pb.TestCase_VAL_TYPES:
		n = d.valTypes()
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

func (d *golang) valTypes() node.Node {
	data := map[string]interface{}{
		"a":     "hello",
		"a-ref": "hello",
	}
	return &nodeutil.Extend{
		Base: &nodeutil.Node{Object: data},
		OnField: func(p node.Node, r node.FieldRequest, hnd *node.ValueHandle) error {
			switch r.Meta.Ident() {
			case "a-ref":
				hnd.Val = val.String(data["a"].(string))
				return nil
			}
			return p.Field(r, hnd)
		},
	}
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

func (d *golang) advancedNode(data map[string]interface{}) node.Node {
	return &nodeutil.Node{Object: data}
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
	if d.traceFile != nil {
		d.traceFile.Close()
	}
	return nil
}

func (d *golang) loadYangModule() *meta.Module {
	if d.yangModule == nil {
		var err error
		if d.yangModule, err = parser.LoadModule(d.ypath, "fc-yang"); err != nil {
			panic(err)
		}
	}
	return d.yangModule
}

func (d *golang) parseModule(dir string, module string) (node.Node, error) {
	ypath := source.Dir(dir)
	m, err := parser.LoadModule(ypath, module)
	if err != nil {
		return nil, err
	}
	return nodeutil.Schema2(m), nil
}

func (d *golang) handleCount() int {
	return 0
}
