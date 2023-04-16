from enum import IntEnum
from pb import fc_x_pb2
from pprint import pprint

class Format(IntEnum):
    INT32 = 11
    INT64 = 13
    STRING = 21

class Val():

    def __init__(self, format, v):
        self.format = format
        self.v = v

def proto_encode(val):
    if not val:
        return None
    if val.format.value == Format.STRING:
        return fc_x_pb2.Val(str=val.v)
    if val.format.value == Format.INT32:
        return fc_x_pb2.Val(i32=val.v)
    if val.format.value == Format.INT64:
        return fc_x_pb2.Val(i64=val.v)
    raise Exception(f'unimplemented value encoder {val.format}')

def proto_decode(proto_val):
    if proto_val == None:
        return None
    which = proto_val.WhichOneof('value')
    if which == 'str':
        return Val(Format.STRING, proto_val.str)
    if which == 'i32':
        return Val(Format.INT32, proto_val.i32)
    if which == 'i64':
        return Val(Format.INT64, proto_val.i64)
    raise Exception(f'unimplemented value decoder {pprint(proto_val)}')
