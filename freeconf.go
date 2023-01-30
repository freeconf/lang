package main

/*
#include "freeconf.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

func c_err_msg(s string) *C.fc_err {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	return C.fc_err_new(cstr)
}

func c_err(err error) *C.fc_err {
	cstr := C.CString(err.Error())
	defer C.free(unsafe.Pointer(cstr))
	return C.fc_err_new(cstr)
}

func go_err(c_err *C.fc_err) error {
	return errors.New(C.GoString(&c_err.message[0]))
}
