package main

/*
#cgo CFLAGS: -Iinclude

#include <stdlib.h>
#include <freeconf/meta.h>

typedef struct fc_path {
    struct fc_path* parent;
    fc_meta* meta;
} fc_path;

*/
import "C"
import "unsafe"

//export fc_new_path
func fc_new_path(parent *C.fc_path, meta *C.fc_meta) *C.fc_path {
	path := (*C.fc_path)(C.malloc(unsafe.Sizeof(fc_path)))
	path.parent = parent
	path.meta = meta
	return path
}
