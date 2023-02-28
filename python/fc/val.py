from enum import Enum
from pb import fc_x_pb2
from pprint import pprint

class Format(Enum):
    STRING = 0
    INT32 = 1
    INT64 = 2

class Val():

    def __init__(self, format, v):
        self.format = format
        self.v = v

def proto_encode(val):
    if val.format == Format.STRING:
        return fc_x_pb2.Val(str=val.v)
    if val.format == Format.INT32:
        return fc_x_pb2.Val(i32=val.v)
    if val.format == Format.INT64:
        return fc_x_pb2.Val(i64=val.v)
    raise Exception(f'unimplemented value encoder {pprint(val)}')

def proto_decode(proto_val):
    if proto_val == None:
        return None
    if proto_val.HasField('str'):
        return Val(Format.STRING, proto_val.str)
    if proto_val.HasField('i32'):
        return Val(Format.INT32, proto_val.i32)
    if proto_val.HasField('i64'):
        return Val(Format.INT64, proto_val.i64)
    raise Exception(f'unimplemented value decoder {pprint(proto_val)}')


