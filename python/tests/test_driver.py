#!/usr/bin/env python3
import unittest 
import freeconf.driver
import freeconf.parser
import freeconf.pb.fc_pb2

class TestDriver(unittest.TestCase):
    def test_load(self):
        d = freeconf.driver.Driver()
        d.load()
        d.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=2))
        d.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=3))
        d.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=4))
        d.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=5))
        d.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=6))
        d.unload()

if __name__ == '__main__':
    unittest.main()