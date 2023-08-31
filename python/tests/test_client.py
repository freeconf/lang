#!/usr/bin/env python3
import sys
import unittest 
from freeconf import driver, parser,  node,  device,  restconf, source, nodeutil


class TestClient(unittest.TestCase):

    def test_client(self):
        drv = driver.Driver()
        drv.load()

        ypath = source.any(
            source.path("./testdata", driver=drv),
            source.restconf_internal_ypath(driver=drv),
        )

        app = {"message": "hello"}
        mgmt = nodeutil.Node(app)
        srv_dev = device.Device(ypath, driver=drv)
        mod = parser.load_module_file(ypath, "x", driver=drv)
        srv_browser = node.Browser(mod, mgmt, driver=drv)
        srv_dev.add_browser(srv_browser)

        _ = restconf.Server(srv_dev, driver=drv)
        srv_dev.apply_startup_config_str('{"fc-restconf":{"web": {"port": ":9998"}}}')

        client_dev = device.Device.client(ypath, "http://localhost:9998/restconf", driver=drv)
        client_browser = client_dev.get_browser("x")
        root = client_browser.root()
        actual =  nodeutil.json.json_write_str(root, driver=drv)
        root.release()
        self.assertEqual('{"message":"hello"}', actual)
        drv.unload()

if __name__ == '__main__':
    unittest.main()
