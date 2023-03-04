#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil
import time
import threading

class Car():

    def __init__(self):
        self.speed = 0
        self.miles = 0
        self.running = False
        self.thread = None

    def start(self, running):
        if self.running != running:
            self.running = running
            if running:
                self.thread = threading.Thread(target=self.run, name="Car")
                self.thread.start()
            if not running:
                self.thread = None

    def reset(self):
        self.miles = 0

    def run(self):
        while self.running:
            time.sleep(0.01)
            self.miles = self.miles + self.speed


def manage_car(c):

    def action(node, req):
        if req.meta.ident == 'stop':
            c.start(False)
        elif req.meta.ident == 'start':
            c.start(True)
        else:
            return node.action(req)
        return None

    return fc.nodeutil.Extend(
        base = fc.nodeutil.Reflect(c),
        on_action = action)


class TestDevice(unittest.TestCase):


    def test_device(self):
        d = fc.driver.Driver()
        d.load()

        node = fc.node.NodeService(d)
        p = fc.parser.Parser(d)
        m = p.load_module('../test/yang', 'car')
        c = Car()  
        n = manage_car(c)
        b = node.new_browser(m, n)
        root = b.root()
        cfg = fc.nodeutil.Reflect({'speed': 10})
        root.upsert_from(cfg)
        print(f'cfg cfg.hnd={cfg.hnd}, manage.hnd={n.hnd}')
        root.find('start').action()
        time.sleep(1)
        self.assertGreater(c.miles, 0)

        d.unload()


if __name__ == '__main__':
    unittest.main()
