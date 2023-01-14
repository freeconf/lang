#include <stdbool.h>

typedef enum def_type
{
    DEF_MODULE,
    DEF_CONTAINER,
    DEF_LIST,
    DEF_LEAF
} Def_type;

typedef struct mod Mod;
typedef struct leaf Leaf;
typedef struct extension Extension;

struct extension {
    char *name;
};

typedef void* datadef_ptr;

struct mod {
    datadef_ptr parent;
    char *ident;
    char *description;
    char *namespace;
    char *prefix;
    char *contact;
    char *org;
    char *ref;
    char *ver;
    Extension** extensions;
    datadef_ptr defs;    
};

struct leaf {
    datadef_ptr parent;
    char *ident;
    char *description;
    char *ref;
    char *units;
    bool *config;
    bool *mandatory;
    int defaultType;
    void* defaultData;
    Extension** extensions;
};