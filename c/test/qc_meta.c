#include "qcbor/qcbor_encode.h"
#include "qcbor/qcbor_decode.h"
#include "qcbor/qcbor_spiffy_decode.h"

char *alloc_str(UsefulBufC buf) {
    char *copy = malloc(buf.len);
    memcpy(copy, buf.ptr, buf.len);
    return copy;
}

int qc_next_expect(QCBORDecodeContext *ctx, QCBORItem *item, int qc_type) {
    QCBORError err = QCBORDecode_GetNext(ctx, item);
    if (err == QCBOR_ERR_NO_MORE_ITEMS) {
        return -1;
    }
    if (item->uDataType != qc_type) {
        printf("expected %d, got %d\n", qc_type, item->uDataType);
        return -1;
    }
    printf("read %d\n", qc_type);
    return 0;
}

int qc_decode_module(QCBORDecodeContext *ctx, DataDef* def) {
    QCBORItem item;
    if (qc_next_expect(ctx, &item, QCBOR_TYPE_MAP) != 0) {
        return -1;
    }
    def->module = (Mod *)calloc(sizeof(Mod), 0);
    if (qc_next_expect(ctx, &item, QCBOR_TYPE_TEXT_STRING) != 0) {
        return -1;
    }
    def->ident = alloc_str(item.val.string);

    if (qc_next_expect(ctx, &item, QCBOR_TYPE_TEXT_STRING) != 0) {
        return -1;
    }
    def->description = alloc_str(item.val.string);

    return 0;
}

int qc_decode(DataDef *def, void* buffer, int len) {
    QCBORDecodeContext ctx;
    UsefulBufC debuf;
    int rc = 0;
    printf("===BEGIN QC==\n");

    debuf.ptr = buffer;
    debuf.len = (size_t)len;    
    QCBORDecode_Init(&ctx, debuf, QCBOR_DECODE_MODE_NORMAL);
    rc = qc_decode_module(&ctx, def);

    printf("===END QC==\n");
    return rc;
}

