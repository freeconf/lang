#!/usr/bin/env python3
import unittest
import freeconf.nodeutil
import freeconf.meta
import freeconf.driver
import freeconf.parser


class dummy:
    x = "Hi"
    q = 100
    z = None

class request:
    write = False
    delete = False
    new = False
    meta = None
    clear = False

class TestReflect(unittest.TestCase):

    def setUp(self):
        d = freeconf.driver.Driver()
        d.load()
        p = freeconf.parser.Parser(driver=d)
        self.module = p.load_module_file('testdata', 'testme-1')
    
    def test_dict_field(self):
        obj = {'x' : "X"}        
        n = freeconf.nodeutil.Reflect(obj)
        r = request()
        r.meta = freeconf.meta.get_def(self.module, 'x')

        # field

        # read
        self.assertEqual("X", n.field(r, None).v)

        # write
        z = freeconf.val.Val(r.meta.type.format, "Z")
        r.write = True
        n.field(r, z)
        self.assertEqual("Z", obj["x"])

    def test_dict_child(self):
        obj = {'x' : "X"}        
        n = freeconf.nodeutil.Reflect(obj)
        r = request()
        r.meta = freeconf.meta.get_def(self.module, 'x')

        # field

        # read
        self.assertEqual("X", n.field(r, None).v)

        # write
        z = freeconf.val.Val(r.meta.type.format, "Z")
        r.write = True
        n.field(r, z)
        self.assertEqual("Z", obj["x"])

    def test_cls_field(self):
        obj = dummy()
        n = freeconf.nodeutil.Reflect(obj)
        r = request()
        r.meta = freeconf.meta.get_def(self.module, 'x')

        # read
        self.assertEqual("Hi", n.field(r, None).v)

        # write
        r.write = True
        n.field(r, freeconf.val.Val(r.meta.type.format, "Bye"))
        self.assertEqual("Bye", obj.x)

    def test_cls_child(self):
        obj = dummy()
        n = freeconf.nodeutil.Reflect(obj)
        r = request()
        r.meta = freeconf.meta.get_def(self.module, 'z')

        # read
        self.assertEqual(None, n.child(r))
        obj.z = {}
        self.assertEqual({}, n.child(r).obj)

        # write
        r.new = True
        obj.z = None
        self.assertEqual({}, n.child(r).obj)
        self.assertEqual({}, obj.z)      


if __name__ == '__main__':
    unittest.main()