#!/usr/bin/env python3
import unittest 
from freeconf import nodeutil, node, parser, source


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

class TestUtilNode(unittest.TestCase):

    def setUp(self):
        self.m = parser.load_module_str(None, mstr)

    def test_read(self):
        class G():
            def __init__(self):
                self.b = 99

        obj = {
            "z":100, 
            "y" : {"q": "yo"}, 
            "g": G(),
            "p":[{"f":"ONE", "b":1},{"f":"TWO", "b":2}],
            "t":{"ONE":{"f":"ONE", "b":1},"TWO":{"f":"TWO", "b":2}},
        }
        n = nodeutil.Node(
            object=obj,
        )
        b = node.Browser(self.m, n)
        actual = nodeutil.json_write_str(b.root())
        expected = '{"z":100,"y":{"q":"yo"},"g":{"b":99},"p":[{"f":"ONE","b":1},{"f":"TWO","b":2}],"t":[{"f":"ONE","b":1},{"f":"TWO","b":2}]}'
        self.assertEqual(expected, actual)

        p_two = nodeutil.json_write_str(b.root().find("p=TWO"))
        self.assertEqual('{"f":"TWO","b":2}', p_two)

        t_two = nodeutil.json_write_str(b.root().find("t=TWO"))
        self.assertEqual('{"f":"TWO","b":2}', t_two)



    def test_write(self):
        class G():
            def __init__(self):
                self.b = 99

        cfg = """{
          "z": 888,
          "y":{"q":"boo!"},
          "g": {"b":444},
          "p":[{"f":"ONE","b":1},{"f":"TWO","b":2}],
          "t":[{"f":"ONE","b":1},{"f":"TWO","b":2}]
        }
        """
        obj = {
            "g": G(),
            "p":[],
            "t":{},
        }
        n = nodeutil.Node(
            object=obj,
        )
        b = node.Browser(self.m, n)
        b.root().upsert_from(nodeutil.json_read_str(cfg))

        self.assertEqual(888, obj["z"])
        self.assertEqual("boo!", obj["y"]["q"])
        self.assertEqual(444, obj["g"].b)
        self.assertEqual(2, len(obj["p"]))
        self.assertEqual("ONE", obj["p"][0]["f"])
        self.assertEqual(2, len(obj["t"]))
        self.assertEqual("ONE", obj["t"]["ONE"]["f"])


if __name__ == '__main__':
    unittest.main()    