package main

/*

#include <stdlib.h>

typedef struct Module {
	long poolId;
	char* ident;
	char* desc;
	void* serialized;
	int   serialized_len;
} Module;
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/freeconf/lang/c/go/emeta"
	"github.com/freeconf/yang/meta"
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
	var buf bytes.Buffer
	if err = emeta.Encode(mod, &buf); err != nil {
		return C.struct_Module{}
	}
	serialized := buf.Bytes()
	m := C.struct_Module{
		ident:          C.CString(mod.Ident()),
		desc:           C.CString(mod.Description()),
		serialized:     C.CBytes(serialized),
		serialized_len: C.int(len(serialized)),
	}
	m.poolId = C.long(pool.Add(modRef{mod: mod}, freeParser(m)))
	return m
}

type modRef struct {
	mod        *meta.Module
	serialized []byte
}

func freeParser(m C.struct_Module) func() {
	return func() {
		if m.ident != nil {
			C.free(unsafe.Pointer(m.ident))
		}
		if m.desc != nil {
			C.free(unsafe.Pointer(m.desc))
		}
		if m.serialized != nil {
			C.free(unsafe.Pointer(m.serialized))
		}
	}
}
