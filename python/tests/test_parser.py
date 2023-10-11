#!/usr/bin/env python3
import sys
import unittest
from freeconf import driver, parser, source

class TestParser(unittest.TestCase):

    def test_basic(self):
        d = driver.Driver()
        d.load()

        # load a module as test that driver is working
        ypath = source.path("testdata", driver=d)
        m = parser.load_module_file(ypath, 'testme-1', driver=d)
        self.assertEqual('testme-1', m.ident)
        self.assertEqual(3, len(m.definitions))
        self.assertEqual('z', m.definitions[0].ident)
        self.assertEqual('x', m.definitions[1].ident)
        z = m.definitions[0]
        self.assertEqual('q', z.definitions[0].ident)
        d.unload()

    def test_car(self):
        d = driver.Driver()
        d.load()

        # load a module as test that driver is working
        ypath = source.path("testdata", driver=d)
        m = parser.load_module_file(ypath, 'car', driver=d)
        self.assertEqual('car', m.ident)
        self.assertEqual(2, len(m.actions))
        start = m.actions['start']
        self.assertEqual('start', start.ident)
        d.unload()

    def test_module_str(self):
        d = driver.Driver()
        d.load()

        # load a module as test that driver is working
        mstr = """
module x {
    container c { }
}
        """
        m = parser.load_module_str(None, mstr, driver=d)
        self.assertEqual('x', m.ident)
        self.assertEqual(1, len(m.definitions))
        d.unload()

    def test_decoder(self):
        d = driver.Driver()
        d.load()

        ypath = source.path("../../test/testdata/yang", driver=d)

        #files = ['car', 'recurse', 'advanced']
        files = ['recurse']
        for f in files:
            parser.load_module_file(ypath, f, driver=d)


if __name__ == '__main__':
    unittest.main()