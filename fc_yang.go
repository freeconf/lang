package main

/*
#include <stdlib.h>
#include <freeconf_lang.h>
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/freeconf/lang/codegen"
	"github.com/freeconf/yang/meta"
	p "github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

//export fc_yang_parse_pack
func fc_yang_parse_pack(ypathPtr *C.char, yfilePtr *C.char) C.fc_pack {
	ypath := C.GoString(ypathPtr)
	yfile := C.GoString(yfilePtr)
	mod, err := p.LoadModule(source.Path(ypath), yfile)
	if err != nil {
		return C.fc_pack{}
	}
	mem_id := mem.Add(mod, nil)
	var buf bytes.Buffer
	if err = codegen.Encode(mod, mem_id, &buf); err != nil {
		return C.fc_pack{}
	}
	serialized := buf.Bytes()
	m := C.fc_pack{
		serialized:     C.CBytes(serialized),
		serialized_len: C.int(len(serialized)),
	}
	return m
}

//export fc_yang_pack_free
func fc_yang_pack_free(p C.fc_pack) {
	if p.serialized != nil {
		C.free(unsafe.Pointer(p.serialized))
	}
}

type modRef struct {
	mod        *meta.Module
	serialized []byte
}
