#!/usr/bin/env python3
import sys
import signal
import logging
import fc.node
import fc.driver
from fc.nodeutil import reflect, trace
import pb.fc_test_pb2
import pb.fc_test_pb2_grpc

usage = f"""
Usage: {sys.argv[0]} fc-g-socket-file fc-x-socket-file

Starts python gRPC service with an additonal test harness service and also
connects to the running Go gRPC service.  Essentially the reverse of what normally
happens.  This lets Go pull this strings to run specific commands orchetrated by
Go unit tests.

"""

if len(sys.argv) < 3:
    print(usage)
    exit(1)

logging.basicConfig(
    filename='harness.log',
    filemode='w',
    level=logging.DEBUG,
    format='%(name)s - %(levelname)s - %(message)s')

class TestHarnessServicer(pb.fc_test_pb2_grpc.TestHarnessServicer):

    def __init__(self, driver):
        self.driver = driver

    def DumpBrowser(self, req, context):
        sel = fc.node.Selection.resolve(self.driver, req.selHnd)
        out = open(req.outputFile, "w")
        n = trace.Trace(reflect.Reflect({}), out)
        sel.upsert_into(n)
        out.close()
        return pb.fc_test_pb2.DumpResponse()

g_addr = sys.argv[1]
x_addr = sys.argv[2]
d = fc.driver.Driver(g_addr, x_addr)
test_harness = TestHarnessServicer(d)
print("loading driver...")
d.load(test_harness)
print("driver loaded.")

def signal_handler(sig, frame):
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)
signal.pause()
