package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"errors"

	"github.com/freeconf/yang/node"
)

//export fc_select_upsert_from
func fc_select_upsert_from(c_sel *C.fc_select, c_node *C.fc_node) *C.fc_err {
	sel, valid := mem.Get(int64(c_sel.node.mem_id)).(node.Selection)
	if !valid {
		return ceeErr(errors.New("no selection found"))
	}
	next := sel.UpdateInto(gnode{c_sel, c_node})
	if next.LastErr != nil {
		return ceeErr(next.LastErr)
	}
	return nil
}
