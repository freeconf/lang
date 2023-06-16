import fc.pb.fc_pb2
import fc.handles

class Device():

    def __init__(self, driver, yangPath, hnd_id=None):
        self.driver = driver
        if not hnd_id:
            req = fc.pb.fc_pb2.NewDeviceRequest(yangPath=yangPath)
            resp = self.driver.g_device.NewDevice(req)
            self.hnd = resp.deviceHnd
        else:
            self.hnd = hnd_id

    def add_browser(self, browser):
        req = fc.pb.fc_pb2.DeviceAddBrowserRequest(deviceHnd=self.hnd, browserHnd=browser.hnd)
        self.driver.g_device.DeviceAddBrowser(req)

    def get_browser(self, moduleIdent):
        req = fc.pb.fc_pb2.DeviceGetBrowserRequest(deviceHnd=self.hnd, moduleIdent=moduleIdent)
        resp = self.driver.g_device.DeviceGetBrowser(req)
        return fc.node.Browser.resolve(self.driver, resp.browserHnd)

    def apply_startup_config(self, configFile):
        req = fc.pb.fc_pb2.ApplyStartupConfigRequest(deviceHnd=self.hnd, configFile=configFile)
        self.driver.g_device.ApplyStartupConfig(req)
