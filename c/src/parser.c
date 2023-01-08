#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "../out/fc-c.h"

int main(int argc, char **argv) {
    char* ypath = getenv("YANGPATH");
    Module m = parser(ypath, "testme");
}