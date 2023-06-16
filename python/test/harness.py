#!/usr/bin/env python3
import sys
import enum
import signal
import logging
import json
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
        self.trace_node = None


    def CreateTestCase(self, req, context):
        out = open(req.traceFile, "w")
        if req.testCase == fc.pb.fc_test_pb2.ECHO or req.testCase == fc.pb.fc_test_pb2.ADVANCED:
            e = Echo()
            n = e.node()
        elif req.testCase == fc.pb.fc_test_pb2.BASIC:
            n = fc.nodeutil.Reflect({})
        else:
            raise Exception("unimplemented test case")
        self.trace_node = fc.nodeutil.Trace(n, out)
        node_hnd = fc.node.ensure_node_hnd(self.driver, self.trace_node)
        self.trace_file = out
        return fc.pb.fc_test_pb2.CreateTestCaseResponse(nodeHnd=node_hnd)
    

    def FinalizeTestCase(self, req, context):
        if self.trace_file != None:
            self.trace_file.close()
            self.trace_file = None
        if self.trace_node != None:
            self.trace_node = None
        return fc.pb.fc_test_pb2.FinalizeTestCaseResponse()


    def ParseModule(self, req, context):
        p = fc.parser.Parser(self.driver)
        m = p.load_module(req.dir, req.moduleIdent)
        dump = {"module":meta_walk(["module"], m)}

        with open(req.dumpFile, 'w') as f:
            json.dump(dump, f, indent=4)
        return fc.pb.fc_test_pb2.ParseModuleResponse()


def meta_walk(path, val):
    """
      print every facet of a module and recurse into every meta object
      in that module. print in a format that matches Go test harness metaDump
      so that dump files can be diff'ed to ensure python is decoding modules
      properly
    """
    context = path[-1]
    if not val:
        return None
    if isinstance(val, dict):
        if len(val) == 0:
            return
        obj = {}
        for attr, attr_val in enumerate(val):
            child = meta_walk(path_push(path, attr), attr_val)
            if child:
                obj[attr] = child
        return obj
    elif isinstance(val, list):
        if len(val) == 0:
            return
        children = []
        for item in val:
            children.append(meta_walk(path, item))
        return children
    elif context == "dataDef":
        return meta_walk_datadef(path, val)
    elif context == "choice":
        return meta_walk_choice(path, val)
    elif hasattr(val, '__module__') and val.__module__ == "fc.meta":
        obj = {}
        for attr in dir(val):
            if attr.startswith("__") or attr == "hnd" or attr == "parent":
                continue
            name = lookup_alias(attr)
            child = meta_walk(path_push(path, name), getattr(val, attr))
            if child:
                obj[name] = child

        if context == "module":
            rev = val.revision()
            if rev:
                obj["revision"] = meta_walk(path_push(path, "revision"), rev)

        return obj
    elif context == "config" and val == True:
        return None
    elif isinstance(val, enum.Enum):
        return val.name.lower()
    else:
        return val


def meta_walk_datadef(path, ddef):
    def_name = lookup_alias(ddef.__class__.__name__.lower())
    data_def = {def_name:meta_walk(path_push(path, def_name), ddef)}
    for pop_ident in ["ident", "description"]:
        if pop_ident in data_def[def_name]:
            data_def[pop_ident] = data_def[def_name].pop(pop_ident)
    return data_def


def meta_walk_choice(path, choice):
     defs = []
     idents = sorted(choice.cases.keys())
     for ident in idents:
        kase = meta_walk(path_push(path, "case"), choice.cases[ident])
        del kase["ident"]
        defs.append({
            "case": kase,
            "ident": ident,
        })
     return {"dataDef": defs, "ident": choice.ident}


aliases = {
    "definitions" : "dataDef",
    "leaflist" : "leaf-list",
}


def lookup_alias(id):
    return aliases.get(id, id)

def path_push(path, item):
    copy = path.copy()
    copy.append(item)
    return copy

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
