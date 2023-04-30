from enum import IntEnum
from pb import val_pb2
from pprint import pprint

class Format(IntEnum):
    BINARY = 1
    BITS = 2
    BOOL = 3
    DECIMAL64 = 4
    EMPTY = 5
    ENUM = 6
    IDENTITY_REF = 7
    INT8 = 9
    INT16 = 10
    INT32 = 11
    INT64 = 12
    STRING = 14
    UINT8 = 15
    UINT16 = 16
    UINT32 = 17
    UINT64 = 18
    BINARY_LIST = 1025
    BITS_LIST = 1026
    BOOL_LIST = 1027
    DECIMAL64_LIST = 1028
    EMPTY_LIST = 1029
    ENUM_LIST = 1030
    IDENTITY_REF_LIST = 1031
    INT8_LIST = 1033
    INT16_LIST = 1034
    INT32_LIST = 1035
    INT64_LIST = 1036
    STRING_LIST = 1038
    UINT8_LIST = 1039
    UINT16_LIST = 1040
    UINT32_LIST = 1041
    UINT64_LIST = 1042

class Val():

    def __init__(self, format, v):
        self.format = format
        self.v = v

def proto_encode(val):
    if val == None:
        return None
    if val.format == Format.BINARY:
        return val_pb2.Val(format=val_pb2.BINARY, value=val_pb2.ValUnion(binary_val=val.v))
    if val.format == Format.BITS:
        return val_pb2.Val(format=val_pb2.BITS, value=val_pb2.ValUnion(bits_val=val.v))
    if val.format == Format.BOOL:
        return val_pb2.Val(format=val_pb2.BOOL, value=val_pb2.ValUnion(bool_val=val.v))
    if val.format == Format.DECIMAL64:
        return val_pb2.Val(format=val_pb2.DECIMAL64, value=val_pb2.ValUnion(decimal64_val=val.v))
    if val.format == Format.EMPTY:
        return val_pb2.Val(format=val_pb2.EMPTY, value=val_pb2.ValUnion(empty_val=val.v))
    if val.format == Format.ENUM:
        return val_pb2.Val(format=val_pb2.ENUM, value=val_pb2.ValUnion(enum_val=val.v))
    if val.format == Format.IDENTITY_REF:
        return val_pb2.Val(format=val_pb2.IDENTITY_REF, value=val_pb2.ValUnion(identity_ref_val=val.v))
    if val.format == Format.INT8:
        return val_pb2.Val(format=val_pb2.INT8, value=val_pb2.ValUnion(int8_val=val.v))
    if val.format == Format.INT16:
        return val_pb2.Val(format=val_pb2.INT16, value=val_pb2.ValUnion(int16_val=val.v))
    if val.format == Format.INT32:
        return val_pb2.Val(format=val_pb2.INT32, value=val_pb2.ValUnion(int32_val=val.v))
    if val.format == Format.INT64:
        return val_pb2.Val(format=val_pb2.INT64, value=val_pb2.ValUnion(int64_val=val.v))
    if val.format == Format.STRING:
        return val_pb2.Val(format=val_pb2.STRING, value=val_pb2.ValUnion(string_val=val.v))
    if val.format == Format.UINT8:
        return val_pb2.Val(format=val_pb2.UINT8, value=val_pb2.ValUnion(uint8_val=val.v))
    if val.format == Format.UINT16:
        return val_pb2.Val(format=val_pb2.UINT16, value=val_pb2.ValUnion(uint16_val=val.v))
    if val.format == Format.UINT32:
        return val_pb2.Val(format=val_pb2.UINT32, value=val_pb2.ValUnion(uint32_val=val.v))
    if val.format == Format.UINT64:
        return val_pb2.Val(format=val_pb2.UINT64, value=val_pb2.ValUnion(uint64_val=val.v))
    if val.format == Format.BINARY_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(binary_val=x_val))
        return val_pb2.Val(format=val_pb2.BINARY_LIST, list_value=vals)
    if val.format == Format.BITS_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(bits_val=x_val))
        return val_pb2.Val(format=val_pb2.BITS_LIST, list_value=vals)
    if val.format == Format.BOOL_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(bool_val=x_val))
        return val_pb2.Val(format=val_pb2.BOOL_LIST, list_value=vals)
    if val.format == Format.DECIMAL64_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(decimal64_val=x_val))
        return val_pb2.Val(format=val_pb2.DECIMAL64_LIST, list_value=vals)
    if val.format == Format.EMPTY_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(empty_val=x_val))
        return val_pb2.Val(format=val_pb2.EMPTY_LIST, list_value=vals)
    if val.format == Format.ENUM_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(enum_val=x_val))
        return val_pb2.Val(format=val_pb2.ENUM_LIST, list_value=vals)
    if val.format == Format.IDENTITY_REF_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(identity_ref_val=x_val))
        return val_pb2.Val(format=val_pb2.IDENTITY_REF_LIST, list_value=vals)
    if val.format == Format.INT8_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(int8_val=x_val))
        return val_pb2.Val(format=val_pb2.INT8_LIST, list_value=vals)
    if val.format == Format.INT16_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(int16_val=x_val))
        return val_pb2.Val(format=val_pb2.INT16_LIST, list_value=vals)
    if val.format == Format.INT32_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(int32_val=x_val))
        return val_pb2.Val(format=val_pb2.INT32_LIST, list_value=vals)
    if val.format == Format.INT64_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(int64_val=x_val))
        return val_pb2.Val(format=val_pb2.INT64_LIST, list_value=vals)
    if val.format == Format.STRING_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(string_val=x_val))
        return val_pb2.Val(format=val_pb2.STRING_LIST, list_value=vals)
    if val.format == Format.UINT8_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(uint8_val=x_val))
        return val_pb2.Val(format=val_pb2.UINT8_LIST, list_value=vals)
    if val.format == Format.UINT16_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(uint16_val=x_val))
        return val_pb2.Val(format=val_pb2.UINT16_LIST, list_value=vals)
    if val.format == Format.UINT32_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(uint32_val=x_val))
        return val_pb2.Val(format=val_pb2.UINT32_LIST, list_value=vals)
    if val.format == Format.UINT64_LIST:
        vals = []
        for x_val in val.v:
            vals.append(val_pb2.ValUnion(uint64_val=x_val))
        return val_pb2.Val(format=val_pb2.UINT64_LIST, list_value=vals)
    raise Exception(f'unimplemented value encoder {val.format}')


def proto_decode(proto_val):
    if proto_val == None:
        return None
    if proto_val.format == val_pb2.BINARY:
        return Val(Format.BINARY, proto_val.value.binary_val)
    if proto_val.format == val_pb2.BITS:
        return Val(Format.BITS, proto_val.value.bits_val)
    if proto_val.format == val_pb2.BOOL:
        return Val(Format.BOOL, proto_val.value.bool_val)
    if proto_val.format == val_pb2.DECIMAL64:
        return Val(Format.DECIMAL64, proto_val.value.decimal64_val)
    if proto_val.format == val_pb2.EMPTY:
        return Val(Format.EMPTY, proto_val.value.empty_val)
    if proto_val.format == val_pb2.ENUM:
        return Val(Format.ENUM, proto_val.value.enum_val)
    if proto_val.format == val_pb2.IDENTITY_REF:
        return Val(Format.IDENTITY_REF, proto_val.value.identity_ref_val)
    if proto_val.format == val_pb2.INT8:
        return Val(Format.INT8, proto_val.value.int8_val)
    if proto_val.format == val_pb2.INT16:
        return Val(Format.INT16, proto_val.value.int16_val)
    if proto_val.format == val_pb2.INT32:
        return Val(Format.INT32, proto_val.value.int32_val)
    if proto_val.format == val_pb2.INT64:
        return Val(Format.INT64, proto_val.value.int64_val)
    if proto_val.format == val_pb2.STRING:
        return Val(Format.STRING, proto_val.value.string_val)
    if proto_val.format == val_pb2.UINT8:
        return Val(Format.UINT8, proto_val.value.uint8_val)
    if proto_val.format == val_pb2.UINT16:
        return Val(Format.UINT16, proto_val.value.uint16_val)
    if proto_val.format == val_pb2.UINT32:
        return Val(Format.UINT32, proto_val.value.uint32_val)
    if proto_val.format == val_pb2.UINT64:
        return Val(Format.UINT64, proto_val.value.uint64_val)
    if proto_val.format == val_pb2.BINARY_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.binary_val)
        return Val(Format.BINARY_LIST, vals)
    if proto_val.format == val_pb2.BITS_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.bits_val)
        return Val(Format.BITS_LIST, vals)
    if proto_val.format == val_pb2.BOOL_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.bool_val)
        return Val(Format.BOOL_LIST, vals)
    if proto_val.format == val_pb2.DECIMAL64_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.decimal64_val)
        return Val(Format.DECIMAL64_LIST, vals)
    if proto_val.format == val_pb2.EMPTY_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.empty_val)
        return Val(Format.EMPTY_LIST, vals)
    if proto_val.format == val_pb2.ENUM_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.enum_val)
        return Val(Format.ENUM_LIST, vals)
    if proto_val.format == val_pb2.IDENTITY_REF_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.identity_ref_val)
        return Val(Format.IDENTITY_REF_LIST, vals)
    if proto_val.format == val_pb2.INT8_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.int8_val)
        return Val(Format.INT8_LIST, vals)
    if proto_val.format == val_pb2.INT16_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.int16_val)
        return Val(Format.INT16_LIST, vals)
    if proto_val.format == val_pb2.INT32_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.int32_val)
        return Val(Format.INT32_LIST, vals)
    if proto_val.format == val_pb2.INT64_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.int64_val)
        return Val(Format.INT64_LIST, vals)
    if proto_val.format == val_pb2.STRING_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.string_val)
        return Val(Format.STRING_LIST, vals)
    if proto_val.format == val_pb2.UINT8_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.uint8_val)
        return Val(Format.UINT8_LIST, vals)
    if proto_val.format == val_pb2.UINT16_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.uint16_val)
        return Val(Format.UINT16_LIST, vals)
    if proto_val.format == val_pb2.UINT32_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.uint32_val)
        return Val(Format.UINT32_LIST, vals)
    if proto_val.format == val_pb2.UINT64_LIST:
        vals = []
        for p_val in proto_val.list_value:
            vals.append(p_val.uint64_val)
        return Val(Format.UINT64_LIST, vals)
    raise Exception(f'unimplemented list value decoder {pprint(proto_val)}')
