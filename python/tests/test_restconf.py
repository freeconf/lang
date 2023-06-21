#!/usr/bin/env python3
import sys
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.node
import freeconf.device
import freeconf.restconf
import requests

sys.path.append(".")
import car

def new_car_app(drv):
    p = freeconf.parser.Parser(drv)
    mod = p.load_module('./testdata', 'car')
    app = car.Car()
    mgmt = car.manage(app)
    b = freeconf.node.Browser(drv, mod, mgmt)
    return b, app

class TestRestconf(unittest.TestCase):


    def test_server(self):
        drv = freeconf.driver.Driver()
        drv.load()

        b, _ = new_car_app(drv)
        dev = freeconf.device.Device(drv, "./testdata:../yang")
        dev.add_browser(b)

        _ = freeconf.restconf.Server(drv, dev)
        dev.apply_startup_config("testdata/server-startup.json")

        resp = requests.get("http://localhost:9999/restconf/data/car:")
        self.assertEqual(200, resp.status_code)
        data = resp.json()
        self.assertEqual(10, data['speed'])

        drv.unload()


if __name__ == '__main__':
    unittest.main()
