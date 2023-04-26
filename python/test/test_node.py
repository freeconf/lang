#!/usr/bin/env python3
import io
import unittest 
import fc.driver
import fc.parser
import fc.node
from fc.nodeutil import reflect, dump, json

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
        dumper = dump.Dump(reflect.Reflect({}), actual)
        b = fc.node.Browser(d, m, dumper)
        rdr = json.rdr(d, "testdata/testme-sample-1.json")
        b.root().upsert_from(rdr)
        expected = """
z:
found=false
z:
found=true
  ->q=(int32.99)
->x=(string.hello)
"""
        self.assertMultiLineEqual(expected, actual.getvalue())
        d.unload()

if __name__ == '__main__':
    unittest.main()