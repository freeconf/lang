#ifndef FC_PATH_H
#define FC_PATH_H

#include <freeconf/meta.h>

typedef struct fc_path {
    struct fc_path* parent;
    fc_meta* meta;
} fc_path;

fc_path* fc_new_path(fc_path* parent, fc_meta* meta);

#endif