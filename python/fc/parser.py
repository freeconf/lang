import fc
import ctypes

class Module(ctypes.Structure):
    _fields_ = [
        ("poolId", ctypes.c_long),
        ("ident", ctypes.c_char_p),
        ("desc", ctypes.c_char_p),
    ]

    def __del__(self):
        fc.library.destruct(self.poolId)

def parser(ypath, fname):
    return fc.library.parser(cstr(ypath), cstr(fname))

def cstr(s):
    return s.encode('utf-8')

fc.library.parser.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
fc.library.parser.restype = Module
