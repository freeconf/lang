#include "freeconf.h"

fc_err* fc_select_child(fc_select selection, fc_node_child_req r, fc_node** child) {
	r.selection = selection;
	return selection.node.on_child(selection.node.context, r, child);
}
