#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include <libfc.h>

#include <freeconf/err.h>
#include <freeconf/node.h>
#include <freeconf/browser.h>
#include <freeconf/parser.h>
#include <freeconf/node_util.h>


fc_node_error* dump_child(void* context, fc_child_request r, fc_node* child) {
    fc_container* meta = (fc_container*)r.meta;
    int rc = fprintf((FILE*)(context), "CHILD %s\n", meta->ident);
    if (rc != 0) {
        return new_node_error("file write error");
    }
    return NULL; 
}

fc_node_error* dump_field(void* context, fc_field_request r, fc_val_ptr* val) {
    fc_leaf* meta = (fc_leaf*)r.meta;
    int rc = fprintf((FILE*)(context), "FIELD %s\n", meta->ident);
    if (rc != 0) {
        return new_node_error("file write error");
    }
    return NULL; 
}


int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    fc_module* m;
    fc_error err = fc_parse_yang(&m, ypath, "testme");
    assert(err == FC_ERR_NONE);

    // assert(strcmp("testme", m->ident) == 0);
    fc_node json;
    fc_node_error* nerr = new_json_rdr(&json, "./test/testdata/testme-sample.json");
    if (nerr != NULL) {
        printf("%s\n", nerr->message);
        assert(nerr == FC_ERR_NONE);
    }
    fc_node dump = {
        .context = stdout,
        .on_field = &dump_field,
        .on_child = &dump_child,
    };
    fc_browser b = fc_new_browser(m, dump);
    fc_selection root;
    fc_root_selection(&root, b);
    nerr = fc_upsert_from(root, json);
    if (nerr != NULL) {
        printf("%s\n", nerr->message);
        assert(nerr == FC_ERR_NONE);
    }
    return 0;
}