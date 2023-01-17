#ifndef FC_PARSER_H
#define FC_PARSER_H

#include <freeconf/err.h>
#include <freeconf/meta.h>

fc_error fc_parse_yang(fc_module** m, char *ypath, char *filename);

#endif
