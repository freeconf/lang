#!/usr/bin/env python3
import sys
import unittest
import freeconf.driver
import freeconf.parser
import freeconf.source

class TestParser(unittest.TestCase):

    def test_basic(self):
        d = freeconf.driver.Driver()
        d.load()

        # load a module as test that driver is working
        ypath = freeconf.source.path("testdata", driver=d)
        m = freeconf.parser.load_module_file(ypath, 'testme-1', driver=d)
        self.assertEqual('testme-1', m.ident)
        self.assertEqual(3, len(m.definitions))
        self.assertEqual('z', m.definitions[0].ident)
        self.assertEqual('x', m.definitions[1].ident)
        z = m.definitions[0]
        self.assertEqual('q', z.definitions[0].ident)
        d.unload()

    def test_car(self):
        d = freeconf.driver.Driver()
        d.load()

        # load a module as test that driver is working
        ypath = freeconf.source.path("testdata", driver=d)
        m = freeconf.parser.load_module_file(ypath, 'car', driver=d)
        self.assertEqual('car', m.ident)
        self.assertEqual(2, len(m.actions))
        start = m.actions['start']
        self.assertEqual('start', start.ident)
        d.unload()

    def test_module_str(self):
        d = freeconf.driver.Driver()
        d.load()

        # load a module as test that driver is working
        mstr = """
module x {
    container c { }
}
        """
        m = freeconf.parser.load_module_str(None, mstr, driver=d)
        self.assertEqual('x', m.ident)
        self.assertEqual(1, len(m.definitions))
        d.unload()

    def test_decoder(self):
        d = freeconf.driver.Driver()
        d.load()

        ypath = freeconf.source.path("testdata", driver=d)

        files = ['car2']
        for f in files:
            freeconf.parser.load_module_file(ypath, f, driver=d)

        d.unload()


if __name__ == '__main__':
    unittest.main()