#!/usr/bin/env python3
import io
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil

class Dump(fc.nodeutil.Basic):

    def __init__(self, out, indent=''):
        super(Dump, self).__init__()
        self.out = out
        self.indent = indent

    def child(self, req):
        self.out.write(f'{self.indent}{req.meta.ident}:\n')
        return Dump(self.out, self.indent + '  ')

    def field(self, req, write_val):
        self.out.write(f'{self.indent}{req.meta.ident}:{write_val.v}\n')
        return None


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
        dumper = Dump(actual)
        b = fc.node.Browser(d, m, dumper)
        rdr = fc.nodeutil.json_rdr(d, "testdata/testme-sample-1.json")
        b.root().upsert_from(rdr)
        expected = """
z:
  q:99
x:hello
"""
        self.assertMultiLineEqual(expected, actual.getvalue())
        d.unload()

if __name__ == '__main__':
    unittest.main()