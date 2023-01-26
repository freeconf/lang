#ifndef FC_META_DECODER_H
#define FC_META_DECODER_H

#include <freeconf/err.h>

typedef struct fc_encoded_module {
	long pool_id;
	void* serialized;
	int   serialized_len;
} fc_encoded_module;

fc_error fc_decode_module(fc_module **m, void* buffer, int len);

#endif
