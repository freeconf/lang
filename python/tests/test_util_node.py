#!/usr/bin/env python3
import unittest 
from freeconf import nodeutil, node, parser


mstr = """
module x {
    leaf z {
        type int32;
    }
    container y {
        leaf q {
            type string;        
        }
    }
    container g {
        leaf b {
            type int32;
        }
        leaf h {
            type string;
        }
    }
    list p {
        key f;
        leaf f {
            type string;
        }
        leaf b {
            type int32;
        }
    }
    list t {
        key f;
        leaf f {
            type string;
        }
        leaf b {
            type int32;
        }    
    }
}
"""

class G():
    def __init__(self):
        self.b = 99
        self.h = "H"
        self.do_something_called = False

    def set_h(self, h):
        self.h = h

    def get_h(self):
        return self.h

    def do_something(self):
        self.do_something_called = True


class TestUtilNode(unittest.TestCase):

    def setUp(self):
        self.m = parser.load_module_str(None, mstr)

    def test_read(self):

        obj = {
            "z":100, 
            "y" : {"q": "yo"}, 
            "g": G(),
            "p":[{"f":"ONE", "b":1},{"f":"TWO", "b":2}],
            "t":{"ONE":{"f":"ONE", "b":1},"TWO":{"f":"TWO", "b":2}},
        }
        n = nodeutil.Node(obj)
        b = node.Browser(self.m, n)
        actual = nodeutil.json_write_str(b.root())
        expected = '{"z":100,"y":{"q":"yo"},"g":{"b":99,"h":"H"},"p":[{"f":"ONE","b":1},{"f":"TWO","b":2}],"t":[{"f":"ONE","b":1},{"f":"TWO","b":2}]}'
        self.assertEqual(expected, actual)

        p_two = nodeutil.json_write_str(b.root().find("p=TWO"))
        self.assertEqual('{"f":"TWO","b":2}', p_two)

        t_two = nodeutil.json_write_str(b.root().find("t=TWO"))
        self.assertEqual('{"f":"TWO","b":2}', t_two)

    def test_read_empty(self):
        obj = {}
        n = nodeutil.Node(obj)
        b = node.Browser(self.m, n)
        actual = nodeutil.json_write_str(b.root())
        self.assertEqual('{}', actual)


    def test_write(self):

        cfg = """{
          "z": 888,
          "y":{"q":"boo!"},
          "g": {"b":444,"h":"Haytch"},
          "p":[{"f":"ONE","b":1},{"f":"TWO","b":2}],
          "t":[{"f":"ONE","b":1},{"f":"TWO","b":2}]
        }
        """
        obj = {
            "g": G(),
            "p":[],
            "t":{},
        }
        n = nodeutil.Node(obj)
        b = node.Browser(self.m, n)
        root = b.root()
        root.upsert_from(nodeutil.json_read_str(cfg))

        self.assertEqual(888, obj["z"])
        self.assertEqual("boo!", obj["y"]["q"])
        self.assertEqual(444, obj["g"].b)
        self.assertEqual("Haytch", obj["g"].get_h())
        self.assertEqual(2, len(obj["p"]))
        self.assertEqual("ONE", obj["p"][0]["f"])
        self.assertEqual(2, len(obj["t"]))
        self.assertEqual("ONE", obj["t"]["ONE"]["f"])

        b.root().find("p=TWO").delete()
        self.assertEqual(1, len(obj["p"]))

        b.root().find("t=TWO").delete()
        self.assertEqual(1, len(obj["t"]))

        root.release()

    def test_options(self):
        class X:
            def __init__(self):
                self.ps = [{"f":"ONE","b":1},{"f":"TWO","b":2}]

        n = nodeutil.Node(
            X(),
            options = nodeutil.NodeOptions(
                try_plural_on_lists=True,
            )
        )
        b = node.Browser(self.m, n)
        self.assertIsNotNone(b.root().find("p=ONE"))


        n = nodeutil.Node(
            {"y":{}},
            options = nodeutil.NodeOptions(
                ignore_empty=True,
            )
        )
        b = node.Browser(self.m, n)
        self.assertIsNone(b.root().find("y"))

if __name__ == '__main__':
    unittest.main()    