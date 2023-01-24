#include "libfc.h"
#include <freeconf/node.h>
#include <freeconf/node_util.h>

fc_node_error* new_json_rdr(fc_node* n, char* fname) {
    return fc_json_node_rdr(n, fname);
}