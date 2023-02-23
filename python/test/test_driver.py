#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser

class TestDriver(unittest.TestCase):
    def test_load(self):
        d = fc.driver.Driver()
        d.load()
        d.unload()

if __name__ == '__main__':
    unittest.main()