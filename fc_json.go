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
func fc_json_node_rdr(c_n **C.struct_fc_node, fname *C.char) *C.fc_err {
	rdr, err := os.Open(C.GoString(fname))
	if err != nil {
		return c_err(err)
	}
	n := nodeutil.ReadJSONIO(rdr)
	(*c_n) = (*C.struct_fc_node)(C.calloc(0, C.size_t(unsafe.Sizeof(*c_n))))
	(*c_n).mem_id = C.long(mem.Add(n, free(unsafe.Pointer(c_n))))
	return nil
}
