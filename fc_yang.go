package main

/*
#include "freeconf.h"
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

//export fc_yang_parse
func fc_yang_parse(m **C.fc_meta_module, ypath *C.char, filename *C.char) C.fc_pack_err {
	pack := fc_yang_parse_pack(ypath, filename)
	defer freePack(pack)
	return C.fc_unpack_fc_meta(m, pack.serialized, pack.serialized_len)
}

type modRef struct {
	mod        *meta.Module
	serialized []byte
}

func freePack(m C.fc_pack) {
	if m.serialized != nil {
		C.free(unsafe.Pointer(m.serialized))
	}
}
