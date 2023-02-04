#include <stdio.h>
#include <freeconf.h>

fc_pack fc_yang_parse_pack(char* ypathPtr, char* yfilePtr);

fc_pack_err fc_yang_parse(fc_meta_module** m, char* ypath, char* filename) {
	fc_pack pack = fc_yang_parse_pack(ypath, filename);
    if (pack.serialized == NULL) {
        return FC_EMPTY_BUFFER;
    }
	fc_pack_err err = fc_unpack_fc_meta(m, pack.serialized, pack.serialized_len);
    free(pack.serialized);
    return err;
}