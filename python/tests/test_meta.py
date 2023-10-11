#!/usr/bin/env python3
import sys
import unittest
from freeconf import driver, parser, source, meta

class TestPath(unittest.TestCase):

    def test_basic(self):
        d = driver.Driver()
        d.load()

        # load a module as test that driver is working
        ypath = source.path("testdata", driver=d)
        m = parser.load_module_file(ypath, 'car2', driver=d)

        wear = meta.Path.find(m, "tire/wear")
        self.assertEqual('wear', wear.ident)

        size = meta.Path.find(wear, "../size")
        self.assertEqual('size', size.ident)

        tire = meta.Path.find(wear, "/tire")
        self.assertEqual('tire', tire.ident)
        d.unload()


if __name__ == '__main__':
    unittest.main()