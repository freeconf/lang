#include <stdbool.h>

typedef enum def_type
{
    DEF_UNKNOWN,
    DEF_MODULE,
    DEF_CONTAINER,
    DEF_LIST,
    DEF_LEAF
} Def_type;

typedef struct dataDef DataDef;
typedef struct mod Mod;
typedef struct leaf Leaf;
typedef struct extension Extension;

struct dataDef {
    DataDef *parent;
    char *ident;
    char *description;
    Def_type type;

    // union?
    Mod *module;
    Leaf *leaf;

    Extension *extensions[];
};

struct extension {
    char *name;
};

struct mod {
    DataDef *parent;
    DataDef def;
    char *namespace;
    char *prefix;
    char *contact;
    char *org;
    char *ref;
    char *ver;    
    DataDef *defs[];
};

struct leaf {
    DataDef* parent;
    DataDef def;
    char *ref;
    char *units;
    bool *config;
    bool *mandatory;
    int defaultType;
    void* defaultData;
};