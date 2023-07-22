import freeconf.pb.fc_pb2
import freeconf.handles
import freeconf.driver

class Device():

    def __init__(self, ypath, hnd_id=None, driver=None):
        self.driver = driver if driver else freeconf.driver.shared_instance()
        if not hnd_id:
            req = freeconf.pb.fc_pb2.NewDeviceRequest(yangPathSourceHnd=ypath.hnd)
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

    def apply_startup_config_file(self, configFile):
        self.__apply_startup_config_stream(self.driver.fs.new_rdr_file(configFile))

    def apply_startup_config_str(self, str):
        self.__apply_startup_config_stream(self.driver.fs.new_rdr_str(str))

    def apply_startup_config_io(self, rdr):
        self.__apply_startup_config_stream(self.driver.fs.new_rdr_io(rdr))

    def __apply_startup_config_stream(self, stream):
        req = freeconf.pb.fc_pb2.ApplyStartupConfigRequest(deviceHnd=self.hnd, streamHnd=stream.hnd)
        self.driver.g_device.ApplyStartupConfig(req)

    @classmethod
    def client(cls, ypath, address, driver=None):
        d = driver if driver else freeconf.driver.shared_instance()
        req = freeconf.pb.fc_pb2.ClientRequest(ypathHnd=ypath.hnd, address=address)
        resp = d.g_device.Client(req)
        return Device(ypath, resp.deviceHnd, driver=d)

