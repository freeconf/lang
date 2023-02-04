import fc
import ctypes
import numpy
from .unpack_meta import unpack_module

class Module():

    def __init__(self):
        self.ident = "todo"


class Pack(ctypes.Structure):
    _fields_ = [
        ("serialized", ctypes.c_void_p),
        ("serialized_len", ctypes.c_int),
    ]


def parser(ypath, fname):
    pack = fc.library.fc_yang_parse_pack(cstr(ypath), cstr(fname))
    try:
        ptr = ctypes.cast(pack.serialized, ctypes.POINTER(ctypes.c_ubyte))
        # works but i'd like to find a way to do this w/o bringing in all of numpy
        buf = numpy.ctypeslib.as_array(ptr, (pack.serialized_len,))
        return unpack_module(buf)
    finally:
        fc.library.fc_yang_pack_free(pack)


def cstr(s):
    return s.encode('utf-8')

fc.library.fc_yang_parse_pack.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
fc.library.fc_yang_parse_pack.restype = Pack
