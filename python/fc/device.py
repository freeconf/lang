import pb.fc_pb2
import fc.handles

class Device():

    def __init__(self, driver, yangPath, hnd_id=None):
        self.driver = driver
        if not hnd_id:
            req = pb.fc_pb2.NewDeviceRequest(yangPath=yangPath)
            resp = self.driver.g_device.NewDevice(req)
            self.hnd = fc.handles.Handle(driver, resp.deviceHnd, self)
        else:
            self.hnd = fc.handles.Handle(driver, hnd_id, self)

    def add_browser(self, browser):
        req = pb.fc_pb2.DeviceAddBrowserRequest(deviceHnd=self.hnd.id, browserHnd=browser.hnd.id)
        self.driver.g_device.DeviceAddBrowser(req)

    def get_browser(self, moduleIdent):
        req = pb.fc_pb2.DeviceGetBrowserRequest(deviceHnd=self.hnd.id, moduleIdent=moduleIdent)
        resp = self.driver.g_device.DeviceGetBrowser(req)
        return fc.node.Browser.resolve(self.driver, resp.browserHnd)

    def apply_startup_config(self, configFile):
        req = pb.fc_pb2.ApplyStartupConfigRequest(deviceHnd=self.hnd.id, configFile=configFile)
        self.driver.g_device.ApplyStartupConfig(req)
