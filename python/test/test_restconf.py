#!/usr/bin/env python3
import unittest 
import fc.driver
import fc.parser
import fc.node
import fc.device
import fc.restconf
import test.car
import requests

def new_car_app(drv):
    p = fc.parser.Parser(drv)
    mod = p.load_module('../test/yang', 'car')
    app = test.car.Car()
    mgmt = test.car.manage(app)
    b = fc.node.Browser(drv, mod, mgmt)
    return b, app

class TestRestconf(unittest.TestCase):


    def test_server(self):
        drv = fc.driver.Driver()
        drv.load()

        b, _ = new_car_app(drv)
        dev = fc.device.Device(drv, "../test/yang:../../restconf/yang")
        dev.add_browser(b)

        _ = fc.restconf.Server(drv, dev)
        dev.apply_startup_config("test/server-startup.json")

        resp = requests.get("http://localhost:9999/restconf/data/car:")
        self.assertEqual(200, resp.status_code)
        data = resp.json()
        self.assertEqual(10, data['speed'])

        drv.unload()


if __name__ == '__main__':
    unittest.main()
