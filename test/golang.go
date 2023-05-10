package test

import (
	"os"

	"github.com/freeconf/lang/pb"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
)

type golang struct {
	traceFile *os.File
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
		n = &nodeutil.Basic{
			OnAction: func(r node.ActionRequest) (node.Node, error) {
				data, err := nodeutil.WriteJSON(r.Input)
				if err != nil {
					return nil, err
				}
				return nodeutil.ReadJSON(data), nil
			},
		}
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

func (d *golang) finalizeTestCase() error {
	return d.traceFile.Close()
}
