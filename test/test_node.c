#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include <libfc.h>
#include <freeconf.h>


fc_err* dump_child(fc_node* self, fc_node_child_req r, fc_node** child) {
    fc_meta_container* meta = (fc_meta_container*)r.meta;
    int rc = fprintf((FILE*)(self->context), "CHILD %s\n", meta->ident);
    if (rc != 0) {
        return fc_err_new("file write error");
    }
    *child = self;
    return NULL;
}

fc_err* dump_field(fc_node* self, fc_node_field_req r, fc_val** val) {
    fc_meta_leaf* meta = (fc_meta_leaf*)r.meta;
    int rc = fprintf((FILE*)(self->context), "FIELD %s\n", meta->ident);
    if (rc != 0) {
        return fc_err_new("file write error");
    }
    return NULL; 
}

fc_node new_dumper() {
    fc_node dump = {
        .context = stdout,
        .on_field = &dump_field,
        .on_child = &dump_child,
    };
    return dump;
}


int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    fc_meta_module* m;
    fc_pack_err err = fc_yang_parse(&m, ypath, "testme");
    assert(err == FC_ERR_NONE);

    // assert(strcmp("testme", m->ident) == 0);
    fc_node json;
    fc_err* nerr = fc_json_node_rdr(&json, "./test/testdata/testme-sample.json");
    if (nerr != NULL) {
        printf("%s\n", nerr->message);
        assert(nerr == FC_ERR_NONE);
    }
    fc_node dump = {
        .context = stdout,
        .on_field = &dump_field,
        .on_child = &dump_child,
    };
    fc_browser b = fc_browser_new(m, dump);
    fc_select root = fc_browser_root_select(b);
    nerr = fc_select_upsert_from(root, json);
    if (nerr != NULL) {
        printf("%s\n", nerr->message);
        assert(nerr == FC_ERR_NONE);
    }
    return 0;
}