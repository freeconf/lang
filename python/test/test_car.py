#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.nodeutil
import time
import test.car

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
        root.upsert_from(fc.nodeutil.Reflect({'speed': 10}))
        root.find('start').action()
        time.sleep(0.1)
        odometer = app.miles
        self.assertGreater(odometer, 0)
        print(f'odometer={odometer}')

        drv.unload()


if __name__ == '__main__':
    unittest.main()
