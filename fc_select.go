package main

/*
#include "freeconf.h"
*/
import "C"
import (
	"github.com/freeconf/yang/node"
)

//export fc_select_upsert_from
func fc_select_upsert_from(c_sel *C.fc_select, c_node *C.fc_node) *C.fc_err {
	mem_id := int64(c_sel.mem_id)
	sel, valid := mem.Get(mem_id).(node.Selection)
	if !valid {
		return c_err_msg("no selection found")
	}
	var g_node node.Node
	if c_node.mem_id != 0 {
		// If node is a go node, use it directly
		g_node = mem.Get(int64(c_node.mem_id)).(node.Node)
	} else {
		g_node = gnode{c_node: c_node, c_path: c_sel.path}
	}
	resp_sel := sel.UpdateFrom(g_node)
	if resp_sel.LastErr != nil {
		return c_err(resp_sel.LastErr)
	}
	return nil
}
