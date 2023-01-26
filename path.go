package main

/*
#include "freeconf.h"
*/
import "C"
import "unsafe"

//export fc_new_path
func fc_new_path(parent *C.fc_path, meta *C.fc_meta) *C.fc_path {
	var path *C.fc_path
	path = (*C.fc_path)(C.malloc(C.size_t(unsafe.Sizeof(*path))))
	path.parent = parent
	path.meta = meta
	return path
}
