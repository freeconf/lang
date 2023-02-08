import unittest 
import grpc
import pb.meta_pb2
import pb.meta_pb2_grpc

class TestParserComm(unittest.TestCase):

    def test_load(self):
        channel = grpc.insecure_channel('unix:///tmp/foo')
        stub = pb.meta_pb2_grpc.ParserStub(channel)
        req = pb.meta_pb2.LoadModuleRequest(dir='../test/yang', name='testme-1')
        m = stub.LoadModule(req)
        self.assertEqual(m.ident, 'testme-1')

if __name__ == '__main__':
    unittest.main()