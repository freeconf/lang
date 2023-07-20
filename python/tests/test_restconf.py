#!/usr/bin/env python3
import sys
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.node
import freeconf.device
import freeconf.restconf
import freeconf.source
import requests

sys.path.append(".")
import car

def new_car_app(drv, ypath):
    mod = freeconf.parser.load_module_file(ypath, 'car', driver=drv)
    app = car.Car()
    mgmt = car.manage(app)
    b = freeconf.node.Browser(mod, mgmt, driver=drv)
    return b, app

class TestRestconf(unittest.TestCase):


    def test_server(self):
        drv = freeconf.driver.Driver()
        drv.load()

        ypath = freeconf.source.any(
            freeconf.source.restconf_internal_ypath(driver=drv),
            freeconf.source.path("testdata")
        )
        b, _ = new_car_app(drv, ypath)
        dev = freeconf.device.Device(ypath, driver=drv)
        dev.add_browser(b)

        _ = freeconf.restconf.Server(dev, driver=drv)
        dev.apply_startup_config("testdata/server-startup.json")

        resp = requests.get("http://localhost:9999/restconf/data/car:")
        self.assertEqual(200, resp.status_code)
        data = resp.json()
        self.assertEqual(10, data['speed'])

        drv.unload()


if __name__ == '__main__':
    unittest.main()
