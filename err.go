package main

/*

#include <stdlib.h>
#include <string.h>
#include <stdio.h>

typedef enum fc_error {
    FC_ERR_NONE,
    FC_BAD_ENCODING,
    FC_EMPTY_BUFFER,
    FC_UNEXPECTED_ENCODING,
    FC_NOT_IMPLEMENTED,
} fc_error;

typedef struct fc_node_error {
	char message[128];
} fc_node_error;

*/
import "C"

import (
	"errors"
	"unsafe"
)

func newNodeErr(err error) *C.fc_node_error {
	cstr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(cstr))
	csz := unsafe.Sizeof(C.fc_node_error)
	cerr := (*C.fc_node_error)(C.malloc(csz))
	C.strncpy(cerr, cstr, csz)
	return cerr
}

func cerrToGo(c_err *C.fc_node_error) error {
	return errors.New(C.GoString(&c_err))
}
