#!/usr/bin/env python3
import io
import sys
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil
from test import gold

class TestNode(unittest.TestCase):
    
    def test_load(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('./testdata', 'testme-1')
        self.assertEqual(m.ident, 'testme-1')

        actual = io.StringIO()
        actual.write('\n')        
        dumper = fc.nodeutil.Trace(fc.nodeutil.Reflect({}), actual)
        b = fc.node.Browser(d, m, dumper)
        rdr = fc.nodeutil.json_read(d, "testdata/testme-sample-1.json")
        print("about to upsert")
        b.root().upsert_from(rdr)
        gold.assert_equal(self, actual.getvalue(), "testdata/gold/node.trace")
        d.unload()

if __name__ == '__main__':
    gold.parse_flags()
    unittest.main()