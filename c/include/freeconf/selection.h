#ifndef FC_SELECTION_H
#define FC_SELECTION_H

#include <freeconf/node.h>
#include <freeconf/meta.h>
#include <freeconf/path.h>

typedef struct fc_selection {
    fc_path* path;
    fc_node* node;
    long pool_id;
} fc_selection;


fc_node_error* fc_upsert_from(fc_selection sel, fc_node node);

#endif