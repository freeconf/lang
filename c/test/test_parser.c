#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "libfc.h"
#include "meta.h"
#include "parser.h"

int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    fc_module m;
    int rc = fc_parse_yang(ypath, "testme", &m);
    assert(strcmp("testme", m.ident) == 0);
    printf("m1.ident=%s, m1.desc=%s, rc1=%d\n", m.ident, m.description, rc);

    // // QCBOR
    // DataDef m2;
    // int rc2 = qc_decode(&m2, m.serialized, m.serialized_len);
    // printf("m2.ident=%s, m2.desc=%s, rc2=%d, dur=%d\n", m2.ident, m2.description, rc2);
}