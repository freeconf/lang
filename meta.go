package main

/*
#cgo CFLAGS: -Iinclude

#include <stdlib.h>
#include <freeconf/meta.h>

*/
import "C"

//export fc_find_meta
func fc_find_meta(hasDefs *C.fc_has_definitions, ident *C.char) *C.fc_meta {
	for i := 0; i < hasDefs.definitions.length; i++ {
		meta := hasDefs.definitions.metas[i]
		if strcmp(meta.ident, ident) == 0 {
			return meta
		}
	}
	return nil
}
