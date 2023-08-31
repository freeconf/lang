#!/usr/bin/env python3
import sys, traceback, threading
import unittest 
from freeconf import nodeutil, driver, node, parser, source
import time

sys.path.append(".")
import car

def dump_threads():
    thread_names = {t.ident: t.name for t in threading.enumerate()}
    for thread_id, frame in sys._current_frames().items():
        print("Thread %s:" % thread_names.get(thread_id, thread_id))
        traceback.print_stack(frame)
        print()

class TestCar(unittest.TestCase):

    def test_start_car(self):
        drv = driver.Driver()
        drv.load()

        ypath = source.path('testdata', driver=drv)
        schema = parser.load_module_file(ypath, 'car', driver=drv)
        app = car.Car()
        mgmt = car.manage(app)
        b = node.Browser(schema, mgmt, driver=drv)
        root = b.root()
        update_called = False
        def update_listener(_msg):
            nonlocal update_called
            update_called = True
        update_sel = root.find('update')
        unsubscribe = update_sel.notification(update_listener)
        root.upsert_from(nodeutil.Node({'speed': 10}))
        update_sel.release()
        start = root.find('start')
        start.action()
        start.release()
        root.release()
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
