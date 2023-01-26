package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"errors"

	"github.com/freeconf/yang/node"
)

//export fc_selection_upsert_from
func fc_selection_upsert_from(c_sel *C.fc_selection, c_node *C.fc_node) *C.fc_node_error {
	sel, valid := pool.Get(int64(c_sel.node.pool_id)).(node.Selection)
	if !valid {
		return newNodeErr(errors.New("no selection found"))
	}
	next := sel.UpdateInto(gnode{c_sel, c_node})
	if next.LastErr != nil {
		return newNodeErr(next.LastErr)
	}
	return nil
}
