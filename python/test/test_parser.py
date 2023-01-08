import unittest 
import os
import fc.parser

class TestParser(unittest.TestCase):
    def test_load(self):
        m = fc.parser.parser(os.environ['YANGPATH'], 'testme')
        self.assertEqual(m.ident, b'testme')

if __name__ == '__main__':
    unittest.main()