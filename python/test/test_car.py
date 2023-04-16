#!/usr/bin/env python3
import sys, traceback, threading
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil
import time
import test.car

def dump_threads():
    thread_names = {t.ident: t.name for t in threading.enumerate()}
    for thread_id, frame in sys._current_frames().items():
        print("Thread %s:" % thread_names.get(thread_id, thread_id))
        traceback.print_stack(frame)
        print()

class TestCar(unittest.TestCase):

    def test_start_car(self):
        drv = fc.driver.Driver()
        drv.load()

        p = fc.parser.Parser(drv)
        schema = p.load_module('../test/yang', 'car')
        app = test.car.Car()
        mgmt = test.car.manage(app)
        b = fc.node.Browser(drv, schema, mgmt)
        root = b.root()
        update_called = False
        def update_listener(_msg):
            nonlocal update_called
            print(f'update_called = {update_called}')
            update_called = True
        update_sel = root.find('update')
        unsubscribe = update_sel.notification(update_listener)
        root.upsert_from(fc.nodeutil.Reflect({'speed': 10}))
        root.find('start').action()
        print("waiting....")
        time.sleep(0.1)
        self.assertTrue(update_called)
        odometer = app.miles
        self.assertGreater(odometer, 0)
        print(f'odometer={odometer}')
        unsubscribe()
        drv.unload()
        
        # useful if test won't exit
        # dump_threads()

if __name__ == '__main__':
    unittest.main()
