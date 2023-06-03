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

func (d *golang) advancedNode(data map[string]interface{}) node.Node {
	return nodeutil.ReflectChild(data)
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

func (d *golang) parseModule(dir string, module string, dumpFilename string) error {
	ypath := source.Dir(dir)
	m, err := parser.LoadModule(ypath, module)
	if err != nil {
		return err
	}
	dumpFile, err := os.Create(dumpFilename)
	if err != nil {
		return err
	}
	defer dumpFile.Close()
	wtr := nodeutil.NewJSONWtr(dumpFile)
	wtr.Pretty = true
	b := nodeutil.Schema(d.loadYangModule(), m)
	err = b.Root().UpsertInto(wtr.Node())
	return err
}
