import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import fc.meta_decoder

class Parser():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.ParserStub(driver.channel)        

    def load_module(self, dir, name):
        req = pb.fc_g_pb2.LoadModuleRequest(dir=dir, name=name)
        resp = self.stub.LoadModule(req)
        m = fc.meta_decoder.Decoder().decode(resp.module)
        m.hnd = resp.moduleHnd
        return m

    def release_module(self, m):
        self.driver.release(m.hnd)

