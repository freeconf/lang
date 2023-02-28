#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil

class Dump(fc.node.Node):

    def __init__(self, store):
        super(Dump, self).__init__()
        self.store = store

    def child(self, req):
        new_store = {}
        self.store[req.meta.ident] = new_store
        return Dump(new_store)

    def field(self, req, write_val):
        self.store[req.meta.ident] = write_val.v
        return None


class TestNode(unittest.TestCase):
    def test_load(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'testme-1')
        self.assertEqual(m.ident, 'testme-1')

        node = fc.node.NodeService(d)
        util = fc.nodeutil.NodeUtilService(d)

        dumper = Dump({})
        b = node.new_browser(m, dumper)
        rdr = util.json_rdr("../test/testdata/testme-sample-1.json")
        b.root().upsert_from(rdr)
        expected = {'x':'hello','z': {'q': 99}}
        self.assertEqual(expected, dumper.store)
        d.unload()

if __name__ == '__main__':
    unittest.main()