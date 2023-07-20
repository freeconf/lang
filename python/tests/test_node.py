#!/usr/bin/env python3
import io
import sys
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.node
import freeconf.nodeutil
import freeconf.source

sys.path.append(".")
import gold

class TestNode(unittest.TestCase):
    
    def test_load(self):
        d = freeconf.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        ypath = freeconf.source.path("testdata", driver=d)
        m = freeconf.parser.load_module_file(ypath, 'testme-1', driver=d)
        self.assertEqual(m.ident, 'testme-1')

        actual = io.StringIO()
        actual.write('\n')        
        dumper = freeconf.nodeutil.Trace(freeconf.nodeutil.Reflect({}), actual)
        b = freeconf.node.Browser(m, dumper, driver=d)
        cfg = """
            {
                "x" : "hello",
                "z" : {
                    "q": 99
                }
            }
        """
        rdr = freeconf.nodeutil.json_read_str(cfg, driver=d)
        print("about to upsert")
        b.root().upsert_from(rdr)
        gold.assert_equal(self, actual.getvalue(), "testdata/gold/node.trace")
        d.unload()

if __name__ == '__main__':
    gold.parse_flags()
    unittest.main()