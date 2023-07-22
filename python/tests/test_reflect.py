#!/usr/bin/env python3
import unittest
import freeconf.nodeutil
import freeconf.meta
import freeconf.driver
import freeconf.parser
import freeconf.source

class dummy:
    x = "Hi"
    q = 100
    z = None
    y = "one"

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
        ypath = freeconf.source.path("testdata", driver=d)
        self.module = freeconf.parser.load_module_file(ypath, 'testme-1', driver=d)
    
    def test_dict_field(self):
        obj = {'x' : "X"}        
        n = freeconf.nodeutil.Reflect(obj)
        r = request()
        r.meta = freeconf.meta.get_def(self.module, 'x')

        # field

        # read
        self.assertEqual("X", n.field(r, None).v)

        # write
        z = freeconf.val.Val("Z", r.meta.type.format)
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
        z = freeconf.val.Val("Z", r.meta.type.format)
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
        n.field(r, freeconf.val.Val("Bye", r.meta.type.format))
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