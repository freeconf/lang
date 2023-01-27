#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "libfc.h"
#include "freeconf.h"

int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    fc_meta_module* m;
    fc_pack_err err = fc_yang_parse(&m, ypath, "testme");
    printf("err=%d\n", err);
    assert(err == FC_ERR_NONE);
    assert(m != NULL);
    // assert(strcmp("testme", m->ident) == 0);
    printf("m1.ident=%s, m.description=%s\n", m->ident, m->description);
}