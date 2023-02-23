import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc

class NodeUtilService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.NodeUtilStub(driver.channel)

    def json_rdr(self, fname):
        req = pb.fc_g_pb2.JSONRdrRequest(fname=fname)
        resp = self.stub.JSONRdr(req)
        return resp.gNodeHnd

