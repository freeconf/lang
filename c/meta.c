#include <stdlib.h>
#include <freeconf/meta.h>

fc_meta* fc_find_meta(fc_has_definitions* p, char* ident) {
    for (int i = 0; i < p->definitions.length; i++) {
        fc_meta* meta = (fc_meta*)p->definitions.metas[i];
        if (strcmp(meta->ident, ident) == 0) {
            return meta;
        }
    }
    return NULL;
}
