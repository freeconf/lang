#include "freeconf.h"

fc_err* fc_node_child(fc_node* node, void* context, fc_node_child_req r, fc_node** child) {
	return node->on_child(context, r, child);
}
