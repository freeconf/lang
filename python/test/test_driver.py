import unittest 
import fc.driver
import fc.parser

class TestDriver(unittest.TestCase):
    def test_load(self):
        d = fc.driver.Driver()
        d.load()
        
        # load a module as test that driver is working
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'testme')
        self.assertEqual(m.ident, 'testme')

        s = fc.node.NodeService(d)
        s.new_browser(m, fc.node.Node())

        d.unload()

if __name__ == '__main__':
    unittest.main()