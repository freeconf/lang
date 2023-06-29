import os
import os.path
import time
import subprocess
import grpc
import freeconf.pb.fc_pb2_grpc
import freeconf.pb.fc_pb2
import freeconf.pb.fc_x_pb2_grpc
import freeconf.pb.fc_x_pb2
import freeconf.node
import weakref

# for cleanup after exit.
from concurrent import futures
import signal
import ctypes


instance = None

def shared_instance():
    """ shared instance of Driver.  Applications generally just need a single instance and so
      this would be typical access unless 
    """

    global instance
    if instance == None:
        instance = Driver()
        instance.load()
    return instance

# Start up the Go executable and create a bi-directional gRPC API with a server in Go
# and a server in python and each side creating clients to respective servers.
class Driver():

    def __init__(self, sock_file=None, x_sock_file=None):
        self.g_proc = None
        cwd = os.getcwd()
        self.sock_file = sock_file if sock_file else f'{cwd}/fc-lang.sock'
        self.x_sock_file = x_sock_file if x_sock_file else f'{cwd}/fc-x.sock'
        if os.path.exists(self.sock_file):
            os.remove(self.sock_file)
        if os.path.exists(self.x_sock_file):
            os.remove(self.x_sock_file)
        self.dbg_addr = os.environ.get('FC_LANG_DBG_ADDR')

    def load(self, test_harness=None):
        if self.g_proc:
            raise Exception("fc-lang already loaded")

        self.obj_strong = HandlePool(self, False) # objects that have an explicit release/destroy
        self.obj_weak = HandlePool(self, True) # objects that should disapear on their own

        self.start_x_server(test_harness)
        if test_harness is None:
            self.start_g_proc()
        self.wait_for_g_connection(self.dbg_addr != None)
        self.create_g_client()

    def start_g_proc(self):
        exec_bin = os.environ.get('FC_LANG_EXEC', 'fc-lang')
        cmd = [exec_bin, self.sock_file, self.x_sock_file]
        if self.dbg_addr:
            dbg = ['dlv', f'--listen={self.dbg_addr}', '--headless=true', '--api-version=2', 'exec']
            dbg.extend(cmd)
            cmd = dbg
        self.g_proc = subprocess.Popen(cmd, preexec_fn=exit_with_parent)

    def wait_for_g_connection(self, wait_forever):
        i = 0
        while i < 20 or wait_forever:
            if os.path.exists(self.sock_file):
                return
            time.sleep(0.5)
            i = i + 1
        raise Exception(f'timed out waiting for {self.sock_file} file')
    

    def create_g_client(self):
        self.g_channel = grpc.insecure_channel(f'unix://{self.sock_file}')
        self.g_handles = freeconf.pb.fc_pb2_grpc.HandlesStub(self.g_channel)
        self.g_parser = freeconf.pb.fc_pb2_grpc.ParserStub(self.g_channel)
        self.g_nodes = freeconf.pb.fc_pb2_grpc.NodeStub(self.g_channel)
        self.g_nodeutil = freeconf.pb.fc_pb2_grpc.NodeUtilStub(self.g_channel)
        self.g_device = freeconf.pb.fc_pb2_grpc.DeviceStub(self.g_channel)
        self.g_restconf = freeconf.pb.fc_pb2_grpc.RestconfStub(self.g_channel)

    def start_x_server(self, test_harness=None):
        self.x_server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        self.x_node_service = freeconf.node.XNodeServicer(self)        
        freeconf.pb.fc_x_pb2_grpc.add_XNodeServicer_to_server(self.x_node_service, self.x_server)
        if test_harness:
            freeconf.pb.fc_test_pb2_grpc.add_TestHarnessServicer_to_server(test_harness, self.x_server)
        self.x_server.add_insecure_port(f'unix://{self.x_sock_file}')
        self.x_server.start()

    def unload(self):
        self.obj_weak.release()
        self.obj_strong.release()
        self.x_server.stop(1).wait()
        self.g_proc.terminate()
        self.g_proc.wait()
        self.g_proc = None


class HandlePool:
    def __init__(self, driver, weak):
        self.weak = weak
        self.driver = driver
        if self.weak:
            self.handles = weakref.WeakValueDictionary() # objects that should disapear on their own
        else:
            self.handles = {}

    def lookup_hnd(self, id):
        return self.handles.get(id, None)

    def require_hnd(self, id):
        try:
            return self.handles[id]
        except KeyError:
            raise KeyError(f'could not resolve hnd {id}')

    def store_hnd(self, id, obj):
        self.handles[id] = obj
        if self.weak:
            weakref.finalize(obj, self.release_hnd, id)
        return id

    def release_hnd(self, id):
        if self.handles != None:
            self.driver.g_handles.Release(freeconf.pb.fc_pb2.ReleaseRequest(hnd=id))

    def release(self):
        self.handles = None


# Ensure fc-lang is terminated when this python process is terminated
# see
#  https://stackoverflow.com/questions/19447603/how-to-kill-a-python-child-process-created-with-subprocess-check-output-when-t/19448096#19448096
#
# does not work on windows, not sure about mac, need to implement different method.
#
libc = ctypes.CDLL("libc.so.6")
def exit_with_parent():
    return libc.prctl(1, signal.SIGTERM)
