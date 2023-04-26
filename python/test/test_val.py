#!/usr/bin/env python3
import unittest
import fc.val
import pb.meta_pb2

class TestVal(unittest.TestCase):

    def test_val(self):
        v = fc.val.Val(fc.val.Format.INT32, 10)
        x = fc.val.proto_encode(v)
        self.assertEqual(10, x.value.int32_val)
        self.assertEqual(pb.meta_pb2.INT32, x.format)
        rt = fc.val.proto_decode(x)
        self.assertEqual(v.format, rt.format)
        self.assertEqual(v.v, rt.v)

    def test_val_list(self):
        v = fc.val.Val(fc.val.Format.INT32_LIST, [10, 12, 14])
        x = fc.val.proto_encode(v)
        self.assertEqual(pb.meta_pb2.INT32_LIST, x.format)
        self.assertEqual(10, x.list_value[0].int32_val)
        self.assertEqual(3, len(x.list_value))
        rt = fc.val.proto_decode(x)
        self.assertEqual(v.format, rt.format)
        self.assertEqual(v.v, rt.v)        

if __name__ == '__main__':
    unittest.main()