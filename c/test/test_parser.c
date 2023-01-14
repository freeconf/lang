#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "libfc.h"

#include "meta.c"
#include "pkj_meta.c"
// #include "qc_meta.c"


int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    Module m = parser(ypath, "testme");
    assert(strcmp("testme", m.ident) == 0);

    // PKJ
    Mod m1;
    int rc1 = pkj_decode(&m1, m.serialized, m.serialized_len);
    printf("m1.ident=%s, m1.desc=%s, rc1=%d\n", m1.ident, m1.description, rc1);

    // // QCBOR
    // DataDef m2;
    // int rc2 = qc_decode(&m2, m.serialized, m.serialized_len);
    // printf("m2.ident=%s, m2.desc=%s, rc2=%d, dur=%d\n", m2.ident, m2.description, rc2);
}