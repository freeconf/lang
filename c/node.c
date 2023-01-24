#include <string.h>
#include <freeconf/node.h>

fc_node_error* new_node_error(char* msg) {
    fc_node_error* err = malloc(sizeof(fc_node_error));
    strncpy(err->message, msg, MAX_ERR_MSG_LEN);
    return err;
}
