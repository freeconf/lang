package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"os"
	"unsafe"

	"github.com/freeconf/yang/nodeutil"
)

//export fc_json_node_rdr
func fc_json_node_rdr(c_n *C.struct_fc_node, fname *C.char) *C.fc_node_error {
	rdr, err := os.Open(C.GoString(fname))
	if err != nil {
		return newNodeErr(err)
	}
	n := nodeutil.ReadJSONIO(rdr)
	c_n = (*C.struct_fc_node)(C.malloc(C.size_t(unsafe.Sizeof(*c_n))))
	c_n.pool_id = C.long(pool.Add(n, freeNode(c_n)))
	return nil
}

func freeNode(c_n *C.struct_fc_node) func() {
	return func() {
		C.free(unsafe.Pointer(c_n))
	}
}
