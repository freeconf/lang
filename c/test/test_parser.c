#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include <freeconf/err.h>
#include <freeconf/meta.h>
#include <freeconf/parser.h>

int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    fc_module* m;
    fc_error err = fc_parse_yang(&m, ypath, "testme");
    printf("err=%d\n", err);
    assert(err == FC_ERR_NONE);
    assert(m != NULL);
    // assert(strcmp("testme", m->ident) == 0);
    printf("m1.ident=%s, m.description=%s\n", m->ident, m->description);
}