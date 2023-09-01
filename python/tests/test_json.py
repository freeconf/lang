#!/usr/bin/env python3
import sys
import unittest 
from freeconf import driver, parser, node, source, nodeutil

sys.path.append(".")
import car


class TestJson(unittest.TestCase):

    def test_read_car_json(self):
        drv = driver.Driver()
        drv.load()

        ypath = source.path("testdata", driver=drv)
        schema = parser.load_module_file(ypath, 'car', driver=drv)
        app = car.Car()
        mgmt = car.manage(app)
        b = node.Browser(schema, mgmt, driver=drv)
        root = b.root()
        cfg = root.find("?content=config")
        actual = nodeutil.json_write_str(cfg, driver=drv)
        cfg.release()
        self.assertEqual('{"speed":0}', actual)
        root.release()

        drv.unload()
        
        # useful if test won't exit
        # dump_threads()

if __name__ == '__main__':
    unittest.main()
