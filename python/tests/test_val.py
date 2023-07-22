#!/usr/bin/env python3
import unittest
import freeconf.val
import freeconf.pb.val_pb2

class TestVal(unittest.TestCase):

    def test_val(self):
        v = freeconf.val.Val(10, freeconf.val.Format.INT32)
        x = freeconf.val.proto_encode(v)
        self.assertEqual(10, x.value.int32_val)
        self.assertEqual(freeconf.pb.val_pb2.INT32, x.format)
        rt = freeconf.val.proto_decode(x)
        self.assertEqual(v.format, rt.format)
        self.assertEqual(v.v, rt.v)

    def test_val_list(self):
        v = freeconf.val.Val([10, 12, 14], freeconf.val.Format.INT32_LIST)
        x = freeconf.val.proto_encode(v)
        self.assertEqual(freeconf.pb.val_pb2.INT32_LIST, x.format)
        self.assertEqual(10, x.list_value[0].int32_val)
        self.assertEqual(3, len(x.list_value))
        rt = freeconf.val.proto_decode(x)
        self.assertEqual(v.format, rt.format)
        self.assertEqual(v.v, rt.v)        

    def test_auto_pick(self):
        v = freeconf.val.Val(10)
        self.assertEqual(freeconf.val.Format.INT32, v.format)

        v = freeconf.val.Val("hello")
        self.assertEqual(freeconf.val.Format.STRING, v.format)

        v = freeconf.val.Val([10, 12, 14])
        self.assertEqual(freeconf.val.Format.INT32_LIST, v.format)


if __name__ == '__main__':
    unittest.main()