#ifndef FC_BROWSER_H
#define FC_BROWSER_H

#include <freeconf/err.h>
#include <freeconf/meta.h>
#include <freeconf/node.h>
#include <freeconf/selection.h>

typedef struct fc_browser {
    fc_module* module;
    fc_node node;
} fc_browser;

fc_browser fc_new_browser(fc_module *m, fc_node n) {
    return (fc_browser) { .module = m, .node = n };
}

void fc_root_selection(fc_selection* root, fc_browser b);

#endif
