#include <stdio.h>
#include "freeconf.h"

fc_select* fc_select_new(long mem_id, fc_node* node, fc_meta_path* path) {
	fc_select* sel = (fc_select*)calloc(1, sizeof(fc_select));
	sel->mem_id = mem_id;
	sel->node = node;
	sel->path = path;
	return sel;
}

void fc_select_delete(fc_select *sel) {
	free(sel->path);
	free(sel);
}

fc_err* fc_select_child(fc_node_child_req r, fc_node** child) {
	return r.selection->node->on_child(r.selection->node, r, child);
}

fc_err* fc_select_field(fc_node_field_req r, fc_val* val) {
	return r.selection->node->on_field(r.selection->node, r, val);
}