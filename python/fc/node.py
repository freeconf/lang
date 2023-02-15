import pb.fc_lang_pb2
import pb.fc_lang_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc

class Node():
    __stubs__ = [
        "handle"
    ]

    def child(self):
        pass


class NodeService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_lang_pb2_grpc.NodeStub(driver.channel)

    def new_browser(self, m, n):
        req = pb.fc_lang_pb2.NewBrowserRequest(moduleHandle=m.handle, nodeHandle=1000)
        hnd = self.stub.NewBrowser(req)
        print(f'handle={hnd.handle}')


class NodeHandler(pb.fc_x_pb2_grpc.XNodeServicer):

    def __init__(self, *args, **kwargs):
        pass

    def Child(self, request, context):
        print("HERE")
        return pb.fc_x_pb2.ChildResponse(handle=2000)

