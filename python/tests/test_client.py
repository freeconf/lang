#!/usr/bin/env python3
import sys
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.node
import freeconf.device
import freeconf.restconf
import freeconf.source
import  freeconf.nodeutil.json


class TestClient(unittest.TestCase):

    def test_client(self):
        drv = freeconf.driver.Driver()
        drv.load()

        ypath = freeconf.source.any(
            freeconf.source.path("./testdata", driver=drv),
            freeconf.source.restconf_internal_ypath(driver=drv),
        )

        app = {"message": "hello"}
        mgmt = freeconf.nodeutil.reflect.Reflect(app)
        srv_dev = freeconf.device.Device(ypath, driver=drv)
        mod = freeconf.parser.load_module_file(ypath, "x", driver=drv)
        srv_browser = freeconf.node.Browser(mod, mgmt, driver=drv)
        srv_dev.add_browser(srv_browser)

        _ = freeconf.restconf.Server(srv_dev, driver=drv)
        srv_dev.apply_startup_config_str('{"fc-restconf":{"web": {"port": ":9998"}}}')

        client_dev = freeconf.device.Device.client(ypath, "http://localhost:9998/restconf", driver=drv)
        client_browser = client_dev.get_browser("x")
        root = client_browser.root()
        actual =  freeconf.nodeutil.json.json_write_str(root, driver=drv)
        root.release()
        self.assertEqual('{"message":"hello"}', actual)
        drv.unload()

if __name__ == '__main__':
    unittest.main()
