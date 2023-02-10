import pb.parser_pb2
import pb.parser_pb2_grpc
import fc.meta_decoder

class Parser():

    def __init__(self, driver):
        self.stub = pb.parser_pb2_grpc.ParserStub(driver.channel)

    def load_module(self, dir, name):
        req = pb.parser_pb2.LoadModuleRequest(dir=dir, name=name)
        resp = self.stub.LoadModule(req)
        m = fc.meta_decoder.Decoder().decode(resp.module)
        m.handle = resp.handle
        return m

    def release_module(m):
        # TODO
        pass

