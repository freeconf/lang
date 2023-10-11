#!/usr/bin/env python3
import sys
import enum
import signal
import logging
from freeconf import node, driver, nodeutil, parser, source, meta, val
import freeconf.pb.fc_test_pb2
import freeconf.pb.fc_test_pb2_grpc

sys.path.append(".")
import schema_dumper


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

class TestHarnessServicer(freeconf.pb.fc_test_pb2_grpc.TestHarnessServicer):

    def __init__(self, driver):
        self.driver = driver
        self.trace_file = None
        self.trace_node = None


    def CreateTestCase(self, req, context):
        out = open(req.traceFile, "w")
        if req.testCase == freeconf.pb.fc_test_pb2.ECHO:
            e = Echo()
            n = e.node()
        elif req.testCase == freeconf.pb.fc_test_pb2.BASIC or req.testCase == freeconf.pb.fc_test_pb2.ADVANCED:
            n = nodeutil.Node({})
        else:
            raise Exception("unimplemented test case")
        self.trace_node = nodeutil.Trace(n, out)
        node_hnd = node.ensure_node_hnd(self.driver, self.trace_node)
        self.trace_file = out
        return freeconf.pb.fc_test_pb2.CreateTestCaseResponse(nodeHnd=node_hnd)
    

    def FinalizeTestCase(self, req, context):
        if self.trace_file != None:
            self.trace_file.close()
            self.trace_file = None
        if self.trace_node != None:
            self.trace_node = None
        return freeconf.pb.fc_test_pb2.FinalizeTestCaseResponse()


    def ParseModule(self, req, context):
        ypath = source.path(req.dir, driver=self.driver)
        m = parser.load_module_file(ypath, req.moduleIdent, driver=self.driver)
        dumper = schema_dumper.Dumper().node(m)
        node_hnd = node.ensure_node_hnd(self.driver, dumper)
        return freeconf.pb.fc_test_pb2.ParseModuleResponse(schemaNodeHnd=node_hnd)


class Echo:
    """
    Implements the test harness compliance tests for "Echo" which includes:
      1.) action that returns back into
      2.) action that triggers a notification that send action input as event message  
    """

    def __init__(self):
        self.listeners = []

    def on_update(self, listener):
        self.listeners.append(listener)
        def closer():
            self.listeners.remove(listener)
        return closer

    def update_listeners(self, data):
        print(f'update_listeners called')
        for l in self.listeners:
            l(data)

    def echo(self, data):
        return data
    
    def send(self, data):
        self.update_listeners(data)

    def node(self):
        
        def notify(p, r):
            if r.meta.ident == "recv":
                def listener(data):
                    r.send(nodeutil.Node(data))
                return self.on_update(listener)
            return None
        
        return nodeutil.Node(self, on_notify=notify)


g_addr = sys.argv[1]
x_addr = sys.argv[2]
d = driver.Driver(g_addr, x_addr)
test_harness = TestHarnessServicer(d)
d.load(test_harness)

def signal_handler(sig, frame):
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)
signal.pause()
