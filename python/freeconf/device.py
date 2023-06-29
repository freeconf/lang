import freeconf.pb.fc_pb2
import freeconf.handles
import freeconf.driver

class Device():

    def __init__(self, yangPath, hnd_id=None, driver=None):
        self.driver = driver if driver else freeconf.driver.shared_instance()
        if not hnd_id:
            req = freeconf.pb.fc_pb2.NewDeviceRequest(yangPath=yangPath)
            resp = self.driver.g_device.NewDevice(req)
            self.hnd = resp.deviceHnd
        else:
            self.hnd = hnd_id

    def add_browser(self, browser):
        req = freeconf.pb.fc_pb2.DeviceAddBrowserRequest(deviceHnd=self.hnd, browserHnd=browser.hnd)
        self.driver.g_device.DeviceAddBrowser(req)

    def get_browser(self, moduleIdent):
        req = freeconf.pb.fc_pb2.DeviceGetBrowserRequest(deviceHnd=self.hnd, moduleIdent=moduleIdent)
        resp = self.driver.g_device.DeviceGetBrowser(req)
        return freeconf.node.Browser.resolve(self.driver, resp.browserHnd)

    def apply_startup_config(self, configFile):
        req = freeconf.pb.fc_pb2.ApplyStartupConfigRequest(deviceHnd=self.hnd, configFile=configFile)
        self.driver.g_device.ApplyStartupConfig(req)
