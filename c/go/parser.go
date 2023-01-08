package main

/*

#include <stdlib.h>

typedef struct Module {
	long poolId;
	char* ident;
	char* desc;
} Module;
*/
import "C"

import (
	"unsafe"

	p "github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

//export parser
func parser(ypathPtr *C.char, yfilePtr *C.char) C.struct_Module {
	ypath := C.GoString(ypathPtr)
	yfile := C.GoString(yfilePtr)
	mod, err := p.LoadModule(source.Path(ypath), yfile)
	if err != nil {
		return C.struct_Module{}
	}
	m := C.struct_Module{
		ident: C.CString(mod.Ident()),
		desc:  C.CString(mod.Description()),
	}
	m.poolId = C.long(pool.Add(mod, freeParser(m)))
	return m
}

func freeParser(m C.struct_Module) func() {
	return func() {
		if m.ident != nil {
			C.free(unsafe.Pointer(m.ident))
		}
		if m.desc != nil {
			C.free(unsafe.Pointer(m.desc))
		}
	}
}
