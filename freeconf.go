package main

/*
#include "freeconf.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func ceeErr(err error) *C.fc_err {
	cstr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(cstr))
	return C.fc_err_new(cstr)
}

func goErr(c_err *C.fc_err) error {
	return errors.New(C.GoString(&c_err.message[0]))
}
