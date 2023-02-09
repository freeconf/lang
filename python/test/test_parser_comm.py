import unittest 
import grpc
import pb.parser_pb2
import pb.parser_pb2_grpc

class TestParserComm(unittest.TestCase):

    def test_load(self):
        channel = grpc.insecure_channel('unix:///tmp/foo')
        stub = pb.parser_pb2_grpc.ParserStub(channel)
        req = pb.parser_pb2.LoadModuleRequest(dir='./test/yang', name='testme-1')
        m = stub.LoadModule(req)
        self.assertEqual('testme-1', m.ident)
        self.assertEqual(2, len(m.definitions))
        self.assertEqual('z', m.definitions[0].container.ident)

if __name__ == '__main__':
    unittest.main()