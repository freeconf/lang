#include <freeconf/path.h>

fc_path* fc_new_path(fc_path* parent, fc_meta* meta) {
    fc_path* path = malloc(sizeof(fc_path));
    path->parent = parent;
    path->meta = meta;
    return path;
}
