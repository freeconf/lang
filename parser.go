package main

/*

#include <stdlib.h>

typedef struct fc_encoded_module {
	long pool_id;
	void* serialized;
	int   serialized_len;
} fc_encoded_module;
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/freeconf/lang/emeta"
	"github.com/freeconf/yang/meta"
	p "github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

//export fc_parse_into_encoded_module
func fc_parse_into_encoded_module(ypathPtr *C.char, yfilePtr *C.char) C.struct_fc_encoded_module {
	ypath := C.GoString(ypathPtr)
	yfile := C.GoString(yfilePtr)
	mod, err := p.LoadModule(source.Path(ypath), yfile)
	if err != nil {
		return C.struct_fc_encoded_module{}
	}
	var buf bytes.Buffer
	if err = emeta.Encode(mod, &buf); err != nil {
		return C.struct_fc_encoded_module{}
	}
	serialized := buf.Bytes()
	m := C.struct_fc_encoded_module{
		serialized:     C.CBytes(serialized),
		serialized_len: C.int(len(serialized)),
	}
	m.pool_id = C.long(pool.Add(modRef{mod: mod}, freeParser(m)))
	return m
}

type modRef struct {
	mod        *meta.Module
	serialized []byte
}

func freeParser(m C.struct_fc_encoded_module) func() {
	return func() {
		if m.serialized != nil {
			C.free(unsafe.Pointer(m.serialized))
		}
	}
}
