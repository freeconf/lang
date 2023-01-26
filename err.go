package main

/*
#include "freeconf.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func newNodeErr(err error) *C.fc_node_error {
	cstr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(cstr))
	var cerr *C.fc_node_error
	cerr = (*C.fc_node_error)(C.malloc(C.size_t(unsafe.Sizeof(*cerr))))
	C.strcpy(&cerr.message[0], cstr)
	return cerr
}

func cerrToGo(c_err *C.fc_node_error) error {
	return errors.New(C.GoString(&c_err.message[0]))
}
