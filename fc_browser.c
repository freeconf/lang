#include "freeconf.h"

fc_browser fc_browser_new(fc_meta_module *m, fc_node *n) {
	fc_browser b = {
		.module = m,
		.node = n
	};
    return b;
}

fc_select fc_browser_root_select(fc_browser b) {
    fc_select sel = {
        .path = fc_meta_path_new(NULL, (fc_meta*)b.module),
        .node = b.node
    };
    return sel;
}
