package lang

import (
	"github.com/freeconf/yang/node"
)

// Provides access to some private methods and data in lang package
// meant for test package only.
type TestHarnessAccess struct {
	d *Driver
}

func NewTestHarnessClient(d *Driver) *TestHarnessAccess {
	return &TestHarnessAccess{
		d: d,
	}
}

func (c *TestHarnessAccess) ResolveSelection(sel node.Selection) uint64 {
	return resolveSelection(c.d, &sel)
}

func (c *TestHarnessAccess) ResolveHnd(hnd uint64) any {
	return c.d.handles.Get(hnd)
}

func (c *TestHarnessAccess) CreateNode(nodeHnd uint64) node.Node {
	return &xnode{d: c.d, nodeHnd: nodeHnd}
}

func (c *TestHarnessAccess) HandleCount() int {
	return len(c.d.handles.handles)
}
