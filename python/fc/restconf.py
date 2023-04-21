import pb.fc_pb2
import fc.handles

class Server():

    def __init__(self, driver, device, hnd_id=None):
        self.driver = driver
        if not hnd_id:
            req = pb.fc_pb2.NewServerRequest(deviceHnd=device.hnd.id)
            resp = self.driver.g_restconf.NewServer(req)
            self.hnd = fc.handles.Handle(driver, resp.serverHnd, self)
        else:
            self.hnd = fc.handles.Handle(driver, hnd_id, self)
