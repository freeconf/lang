#include <stdio.h>
#include "freeconf.h"

fc_select* fc_select_new(long mem_id, fc_node* node, fc_meta_path* path) {
	fc_select* sel = (fc_select*)calloc(0, sizeof(fc_select));
	sel->mem_id = mem_id;
	sel->node = node;
	sel->path = path;
	return sel;
}

fc_select* fc_select_new_err_msg(char* msg) {
	fc_select* sel = (fc_select*)calloc(0, sizeof(fc_select));
	sel->last_err = fc_err_new(msg);
	return sel;
}

fc_select* fc_select_new_err(fc_err* err) {
	fc_select* sel = (fc_select*)calloc(0, sizeof(fc_select));
	sel->last_err = err;
	return sel;
}

fc_err* fc_select_child(fc_node_child_req r, fc_node** child) {
	printf("child node=%p sel=%p\n", r.selection->node, r.selection);
	return r.selection->node->on_child(r.selection->node, r, child);
}

fc_err* fc_select_field(fc_node_field_req r, fc_val* val) {
	printf("field node=%p sel=%p\n", r.selection->node, r.selection);
	//return NULL;
	return r.selection->node->on_field(r.selection->node, r, val);
}