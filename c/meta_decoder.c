#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <cbor.h>
#include <stdio.h>
#include "meta.h"

char *fc_decode_alloc_str(cbor_item_t *item) {
    char *src = (char *)cbor_string_handle(item);
    size_t len = cbor_string_length(item);
    char *dst = malloc(len + 1);
    memcpy(dst, src, len);
    dst[len] = 0;
    return dst;
}

int fc_decode_expect(cbor_item_t* item, int expected)  {
    if (cbor_typeof(item) != expected) {
        printf("expected %d but got %d\n", expected, cbor_typeof(item));
        return -1;
    }    
    return 0;
}

int fc_decode_module_def(fc_module* def, cbor_item_t* item) {
    if (fc_decode_expect(item, CBOR_TYPE_ARRAY)) {
        return -1;
    }
    cbor_item_t** array = cbor_array_handle(item);
    def->ident = fc_decode_alloc_str(array[0]);
    def->description = fc_decode_alloc_str(array[1]);
    // ext
    // int type = cbor_get_uint32(array[3]);
    // defs   
    return 0;
}

int fc_decode_module(fc_module *def, void* buffer, int len) {
    struct cbor_load_result result;
    cbor_item_t* item = cbor_load(buffer, len, &result);
    if (item == NULL) {
        if (result.error.code != CBOR_ERR_NONE) {
            return -1;
        }
        return -2;
    }
    // cbor_describe(item, stdout);
    int rc = fc_decode_module_def(def, item);
    cbor_decref(&item);
    return rc;
}
