#include "freeconf.h"

// cgo cannot call C function pointers directly so we wrap them here
// https://pkg.go.dev/cmd/cgo

fc_node_error* fc_node_on_child_x(fc_node* node, void* context, fc_child_request r, fc_node** child) {
	return node->on_child(context, r, child);
}

