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
