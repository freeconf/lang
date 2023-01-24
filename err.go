package main

/*
#cgo CFLAGS: -Ic/include

#include <stdlib.h>
#include <freeconf/err.h>
#include <freeconf/node.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func newNodeErr(err error) *C.fc_node_error {
	cstr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(cstr))
	return C.new_node_error(cstr)
}

func cerrToGo(c_err *C.fc_node_error) error {
	return errors.New(C.GoString(c_err.message))
}
