#!/usr/bin/env python3
import sys
import signal
import logging
import fc.node
import fc.driver
from fc.nodeutil import reflect, trace, extend
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
        self.trace_file = None


    def CreateTestCase(self, req, context):
        out = open(req.traceFile, "w")
        if req.testCase == pb.fc_test_pb2.ECHO:
            n = echo_node()
        elif req.testCase == pb.fc_test_pb2.BASIC:
            n = reflect.Reflect({})
        else:
            raise Exception("unimplemented test case")
        n = trace.Trace(n, out)
        fc.node.resolve_node(self.driver, n)
        self.trace_file = out
        return pb.fc_test_pb2.CreateTestCaseResponse(nodeHnd=n.hnd.id)
    

    def FinalizeTestCase(self, req, context):
        self.trace_file.close()
        return pb.fc_test_pb2.FinalizeTestCaseResponse()


def echo_node():
    base = reflect.Reflect({})
    def on_action(parent, r):
        n = reflect.Reflect({})
        print(f'r.input[{r.input.hnd.id}].insert_into()')
        r.input.insert_into(n)
        return n
    return extend.Extend(base, on_action=on_action)


g_addr = sys.argv[1]
x_addr = sys.argv[2]
d = fc.driver.Driver(g_addr, x_addr)
test_harness = TestHarnessServicer(d)
d.load(test_harness)

def signal_handler(sig, frame):
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)
signal.pause()
