#include <string.h>
#include "libfc.h"
#include "meta.h"
#include "meta_decoder.h"

int fc_parse_yang(char* ypath, char* filename, fc_module* m) {
    fc_encoded_module encoded = fc_parse_into_encoded_module(ypath, filename);
    memset(m, 0, sizeof(fc_module));
    int rc = fc_decode_module(m, encoded.serialized, encoded.serialized_len);
    return rc;
}
