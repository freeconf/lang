#!/usr/bin/env python3
import sys
import signal
import logging
import fc.node
import fc.driver
import fc.nodeutil
import fc.pb.fc_test_pb2
import fc.pb.fc_test_pb2_grpc
import fc.parser

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

class TestHarnessServicer(fc.pb.fc_test_pb2_grpc.TestHarnessServicer):

    def __init__(self, driver):
        self.driver = driver
        self.trace_file = None


    def CreateTestCase(self, req, context):
        out = open(req.traceFile, "w")
        if req.testCase == fc.pb.fc_test_pb2.ECHO or req.testCase == fc.pb.fc_test_pb2.ADVANCED:
            e = Echo()
            n = e.node()
        elif req.testCase == fc.pb.fc_test_pb2.BASIC:
            n = fc.nodeutil.Reflect({})
        else:
            raise Exception("unimplemented test case")
        n = fc.nodeutil.Trace(n, out)
        fc.node.resolve_node(self.driver, n)
        self.trace_file = out
        return fc.pb.fc_test_pb2.CreateTestCaseResponse(nodeHnd=n.hnd.id)
    

    def FinalizeTestCase(self, req, context):
        if self.trace_file != None:
            self.trace_file.close()
            self.trace_file = None
        return fc.pb.fc_test_pb2.FinalizeTestCaseResponse()


    def ParseModule(self, req, context):
        p = fc.parser.Parser(self.driver)
        p.load_module(req.dir, req.moduleIdent)
        return fc.pb.fc_test_pb2.ParseModuleResponse()


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

    def update_listeners(self, n):
        print(f'update_listeners called')
        for l in self.listeners:
            l(n)

    def node(self):
        base = fc.nodeutil.Reflect({})

        def action(parent, r):
            if r.meta.ident == "echo":
                n = fc.nodeutil.Reflect({})
                r.input.insert_into(n)
                return n
            elif r.meta.ident == "send":
                self.update_listeners(r.input.node)
            return None
        
        def notification(node, r):
            print(f'notify called')
            if r.meta.ident == "recv":
                def listener(n):
                    r.send(n)
                return self.on_update(listener)
        
        return fc.nodeutil.Extend(base, on_action=action, on_notification=notification)


g_addr = sys.argv[1]
x_addr = sys.argv[2]
d = fc.driver.Driver(g_addr, x_addr)
test_harness = TestHarnessServicer(d)
d.load(test_harness)

def signal_handler(sig, frame):
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)
signal.pause()
