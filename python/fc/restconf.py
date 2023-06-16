import fc.pb.fc_pb2
import fc.handles

class Server():

    def __init__(self, driver, device):
        self.driver = driver
        req = fc.pb.fc_pb2.NewServerRequest(deviceHnd=device.hnd)
        resp = self.driver.g_restconf.NewServer(req)
        self.hnd = driver.obj_weak.store_hnd(resp.serverHnd, self)
