#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil
import time

class Car():


    def __init__(self):
        self.speed = 0
        self.miles = 0
        self.running = False


    def start(self, running):
        self.running = running


    def reset(self):
        self.miles = 0


    def run(self):
        while True:
            time.sleep(0.01)
            if self.running:
                self.miles = self.miles + self.speed


def manage_car(c):


    def on_action(node, req):
        if req.meta.ident == 'stop':
            c.start(False)
        elif req.meta.ident == 'start':
            c.start(True)
        else:
            return node.action(req)
        return None


    return fc.nodeutil.Extend(
        base = fc.nodeutil.Reflect(c),
        on_action = on_action)


class TestDevice(unittest.TestCase):


    def test_device(self):
        d = fc.driver.Driver()
        d.load()

        node = fc.node.NodeService(d)
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'car')
        c = Car()    
        b = node.new_browser(m, manage_car(c))
        root = b.root()
        root.upsert_from(fc.nodeutil.Reflect({'speed': 10}))
        root.find('start').action()
        time.sleep(0.1)
        self.assertGreater(c.miles, 0)

        d.unload()


if __name__ == '__main__':
    unittest.main()
