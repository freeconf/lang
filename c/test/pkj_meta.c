// Uses PJK/libcbor: CBOR protocol implementation for C
// https://github.com/PJK/libcbor
// 
// While I didn't find a problem, API appears fairly basic.
#include <cbor.h>
#include <stdio.h>

size_t pkj_peek_str(cbor_item_t *item, char *dst) {
    char *src = (char *)cbor_string_handle(item);
    size_t len = cbor_string_length(item);
    memcpy(dst, src, len);
    dst[len] = 0;
    return len;
}

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

/*
  int pkj_decode(ctx, def, item) {
    expect_map()
    if (decode_data_def(def, item)) {
        return ERR;
    }
    def->mod = calloc(size_of(Mod))
    Mod mod = calloc(size_of(Mod))
    if (decode_string(&mod->namespace)) {
        return ERR;
    }
}
*/

int pkj_decode_module(DataDef* def, cbor_item_t* item) {
    if (pkj_expect(item, CBOR_TYPE_MAP)) {
        return -1;
    }
    def->ident = pkj_alloc_str(cbor_map_handle(item)[0].value);
    def->description = pkj_alloc_str(cbor_map_handle(item)[1].value);

    return 0;
}


int pkj_decode(DataDef *def, void* buffer, int len) {
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
    // printf("item=%d\n", item->type);
    cbor_decref(&item);
    return rc;
}
