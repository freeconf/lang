#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil

class Dump(fc.node.Node):

    def child(self, req):
        print("Child from python")
        return self

    def field(self, req):
        print("Field from python")
        return None


class TestNode(unittest.TestCase):
    def test_load(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'testme')
        self.assertEqual(m.ident, 'testme')

        node = fc.node.NodeService(d)
        util = fc.nodeutil.NodeUtilService(d)

        b = node.new_browser(m, Dump())
        rdr = util.json_rdr("../test/testdata/testme-sample.json")
        b.root().upsert_from(rdr)

        d.unload()

if __name__ == '__main__':
    unittest.main()