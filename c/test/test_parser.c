#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>

#include "libfc.h"

int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    Module m = parser(ypath, "testme");
    assert(strcmp("testme", m.ident) == 0);
    destruct(m.poolId);
}