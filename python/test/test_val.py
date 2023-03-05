#!/usr/bin/env python3
import unittest
import fc.val
import pb.meta_pb2

class TestVal(unittest.TestCase):

    def test_val(self):
        pb_type = pb.meta_pb2.INT32
        meta_type = fc.val.Format(pb_type)
        self.assertTrue(fc.val.Format.INT32 == meta_type)
        v = fc.val.Val(fc.val.Format.INT32, 10)
        x = fc.val.proto_encode(v)
        self.assertEqual(x.i32, 10)

if __name__ == '__main__':
    unittest.main()