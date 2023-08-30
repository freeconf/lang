import freeconf.pb.fc_pb2
import freeconf.handles
import freeconf.driver

class Server():

    def __init__(self, device, driver=None):
        self.driver = driver if driver else freeconf.driver.shared_instance()
        req = freeconf.pb.fc_pb2.GnmiServerRequest(deviceHnd=device.hnd)
        resp = self.driver.g_proto.GnmiServer(req)
        self.hnd = self.driver.obj_weak.store_hnd(resp.serverHnd, self)
