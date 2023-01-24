#ifndef FC_NODE_H
#define FC_NODE_H

#include <stdlib.h>
#include <freeconf/err.h>
#include <freeconf/val.h>
#include <freeconf/meta.h>

#define MAX_ERR_MSG_LEN 512

typedef struct fc_node_error {
    char message[MAX_ERR_MSG_LEN];
} fc_node_error;

fc_node_error* new_node_error(char* msg);

typedef void* fc_context;

typedef struct fc_field_request {    
    fc_context context;
    fc_meta* meta;
    bool write;
} fc_field_request;

typedef struct fc_child_request {
    fc_context context;
    fc_has_definitions* meta;
} fc_child_request;

typedef struct fc_val {
    fc_val_type type;
    void* data;
    size_t size;
} fc_val;

typedef fc_val* fc_val_ptr;

typedef struct fc_node fc_node;

struct fc_node {
    long pool_id;
    void* context;
    fc_node_error* (*on_field)(void* context, fc_field_request r, fc_val_ptr* val);
    fc_node_error* (*on_child)(void* context, fc_child_request r, fc_node* child);
};

#endif