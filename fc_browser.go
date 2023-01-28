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
func fc_browser_new(m *C.fc_meta_module, n C.fc_node) C.fc_browser {
	gm := mem.Get(int64(m.mem_id)).(*meta.Module)
	gb := node.NewBrowser(gm, gnode{cnode: n})
	b := C.fc_browser{
		//mem_id: C.long(mem.Add(gb, nil)),
		module: m,
		node:   n,
	}
	return b
}

//export fc_browser_root_select
func fc_browser_root_select(b C.fc_browser) C.fc_select {
	gb := mem.Get(int64(b.mem_id)).(*node.Browser)
	root := gb.Root()
	meta := (*C.fc_meta)(unsafe.Pointer(b.module))
	mem_id := mem.Add(root, nil)
	return C.fc_select{
		path:   C.fc_meta_path_new(nil, meta),
		node:   b.node,
		mem_id: C.long(mem_id),
	}
}
