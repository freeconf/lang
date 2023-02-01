package main

/*
#include <freeconf.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/freeconf/yang/val"
)

func cee_val(v val.Value) (c_val C.fc_val) {
	switch v.Format() {
	case val.FmtString:
		c_val.str = C.CString(v.String())
		c_val.data_type = C.FC_VAL_STRING
	case val.FmtInt32:
		x := v.Value().(int)
		c_val.i32 = C.int(x)
		c_val.data_type = C.FC_VAL_INT_32
	case val.FmtDecimal64:
		x := v.Value().(float64)
		c_val.decimal = C.float(x)
		c_val.data_type = C.FC_VAL_DECIMAL_64
	default:
		panic(fmt.Sprintf("%s not implemented", v.Format().String()))
	}
	return
}

func go_val(c_val C.fc_val) val.Value {
	switch c_val.data_type {
	case C.FC_VAL_STRING:
		return val.String(C.GoString(c_val.str))
	case C.FC_VAL_INT_32:
		return val.Int32(int(c_val.i32))
	case C.FC_VAL_DECIMAL_64:
		return val.Decimal64(float64(c_val.decimal))
	default:
		panic(fmt.Sprintf("%d not implemented", c_val.data_type))
	}
}

func free_val(c_val C.fc_val) {
	switch c_val.data_type {
	case C.FC_VAL_STRING:
		C.free(unsafe.Pointer(c_val.str))
	}
}
