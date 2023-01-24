package main

/*
#cgo CFLAGS: -Iinclude

#include <freeconf/meta.h>

typedef struct fc_browser {
    fc_module* module;
    struct fc_node* node;
} fc_browser;

*/
import "C"

//export fc_new_browser
func fc_new_browser(m *C.fc_module, n *C.struct_fc_node) C.fc_browser {
	return C.fc_browser{
		module: m,
		node:   n,
	}
}

//export fc_root_selection
func fc_root_selection(root *C.struct_fc_selection, b *C.fc_browser) {
	// TODO
}
