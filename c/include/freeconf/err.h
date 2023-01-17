#ifndef FC_ERR_H
#define FC_ERR_H

typedef enum fc_error_e {
    FC_ERR_NONE,
    FC_BAD_ENCODING,
    FC_EMPTY_BUFFER,
    FC_UNEXPECTED_ENCODING,
    FC_NOT_IMPLEMENTED,    
} fc_error;

#endif
