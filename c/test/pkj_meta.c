// Uses PJK/libcbor: CBOR protocol implementation for C
// https://github.com/PJK/libcbor
#include <cbor.h>
#include <stdio.h>

char *pkj_alloc_str(cbor_item_t *item) {
    char *src = (char *)cbor_string_handle(item);
    size_t len = cbor_string_length(item);
    char *dst = malloc(len + 1);
    memcpy(dst, src, len);
    dst[len] = 0;
    return dst;
}

int pkj_expect(cbor_item_t* item, int expected)  {
    if (cbor_typeof(item) != expected) {
        printf("expected %d but got %d\n", expected, cbor_typeof(item));
        return -1;
    }    
    return 0;
}

int pkj_decode_module(Mod* def, cbor_item_t* item) {
    if (pkj_expect(item, CBOR_TYPE_ARRAY)) {
        return -1;
    }
    cbor_item_t** array = cbor_array_handle(item);
    def->ident = pkj_alloc_str(array[0]);
    def->description = pkj_alloc_str(array[1]);
    // ext
    // int type = cbor_get_uint32(array[3]);
    // defs
    
    return 0;
}


int pkj_decode(Mod *def, void* buffer, int len) {
    struct cbor_load_result result;
    cbor_item_t* item = cbor_load(buffer, len, &result);
    if (item == NULL) {
        if (result.error.code != CBOR_ERR_NONE) {
            return -1;
        }
        return -2;
    }
    // cbor_describe(item, stdout);
    int rc = pkj_decode_module(def, item);
    cbor_decref(&item);
    return rc;
}
