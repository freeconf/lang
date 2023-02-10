import unittest 
import grpc
import pb.parser_pb2
import pb.parser_pb2_grpc
import fc.meta_decoder
import fc.driver

class TestParserComm(unittest.TestCase):

    def test_load(self):
        channel = grpc.insecure_channel('unix:///tmp/foo')
        stub = pb.parser_pb2_grpc.ParserStub(channel)
        # fixme: path is relative to where fc-lang is started
        req = pb.parser_pb2.LoadModuleRequest(dir='./test/yang', name='testme-1')
        resp = stub.LoadModule(req)
        m = fc.meta_decoder.Decoder().decode(resp.module)
        m.handle = resp.handle
        self.assertEqual('testme-1', m.ident)
        self.assertEqual(2, len(m.definitions))
        self.assertEqual('z', m.definitions[0].ident)

if __name__ == '__main__':
    unittest.main()