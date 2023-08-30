#!/usr/bin/env python3
import unittest 
from freeconf import nodeutil, node, parser

class TestAction(unittest.TestCase):

    def test_action(self):
        mstr = """module x {
            rpc a {}

            rpc echo {
                input {
                    leaf x {
                        type string;
                    }
                }
                output {
                    leaf x {
                        type string;
                    }
                }
            }

            rpc echoExplode {
                input {
                    leaf x1 {
                        type string;
                    }
                    leaf x2 {
                        type int32;
                    }
                }
                output {
                    leaf x1 {
                        type string;
                    }
                    leaf x2 {
                        type int32;
                    }
                }
            }            
        }
        """
        m = parser.load_module_str(None, mstr)
        class X():
            def __init__(self):
                self.a_called = False
            
            def a(self):
                self.a_called = True

            def echo(self, input):
                return input
            
            def echo_explode(self, x1, x2):
                return x1, x2

        obj = X()

        n = nodeutil.Node(
            object=obj,
        )
        b = node.Browser(m, n)

        b.root().find("a").action(None)
        self.assertTrue(obj.a_called)

        resp = b.root().find("echo").action(nodeutil.Node({"x":"hello"}))
        actual = nodeutil.json_write_str(resp)
        self.assertEqual('{"x":"hello"}', actual)

        n = nodeutil.Node(
            object=obj,
            options=nodeutil.NodeOptions(action_input_exploded=True, action_output_exploded=True),
        )
        b = node.Browser(m, n)
        resp = b.root().find("echoExplode").action(nodeutil.Node({"x1":"hello", "x2": 10}))
        actual = nodeutil.json_write_str(resp)
        self.assertEqual('{"x1":"hello","x2":10}', actual)

if __name__ == '__main__':
    unittest.main()            