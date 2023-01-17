#include <string.h>
#include "libfc.h"
#include <freeconf/meta.h>
#include "meta_decoder.h"

fc_error fc_parse_yang(fc_module** m, char* ypath, char* filename) {
    fc_encoded_module encoded = fc_parse_into_encoded_module(ypath, filename);
    return fc_decode_module(m, encoded.serialized, encoded.serialized_len);
}
