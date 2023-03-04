#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser

class TestParser(unittest.TestCase):

    def test_basic(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'testme-1')
        self.assertEqual('testme-1', m.ident)
        self.assertEqual(2, len(m.definitions))
        self.assertEqual('z', m.definitions[0].ident)
        self.assertEqual('x', m.definitions[1].ident)
        z = m.definitions[0]
        self.assertEqual('q', z.definitions[0].ident)
        p.release_module(m)
        
        d.unload()

    def test_car(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'car')
        self.assertEqual('car', m.ident)
        self.assertEqual(2, len(m.actions))
        start = m.actions['start']
        self.assertEqual('start', start.ident)
        
        d.unload()        

if __name__ == '__main__':
    unittest.main()