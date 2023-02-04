package main

/*
#include <freeconf.h>
*/
import "C"

import (
	"unsafe"

	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
)

//export fc_browser_new
func fc_browser_new(m *C.fc_meta_module, c_node *C.fc_node) *C.fc_browser {
	gm := mem.Get(int64(m.mem_id)).(*meta.Module)
	var b *C.fc_browser
	b = (*C.fc_browser)(C.calloc(1, C.size_t(unsafe.Sizeof(*b))))
	b.module = m
	b.node = c_node
	c_path := C.fc_meta_path_new_root(m)
	n := &gnode{c_node: c_node, c_path: c_path}
	gb := node.NewBrowser(gm, n)
	b.mem_id = C.long(mem.Add(gb, nil))
	return b
}

//export fc_browser_root_select
func fc_browser_root_select(b *C.fc_browser) *C.fc_select {
	gb := mem.Get(int64(b.mem_id)).(*node.Browser)
	root := gb.Root()
	mem_id := mem.Add(root, nil)
	path := C.fc_meta_path_new_root(b.module)
	return C.fc_select_new(C.long(mem_id), b.node, path)
}
