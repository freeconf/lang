import os
import os.path
import time
import subprocess
import grpc
import freeconf.pb.fc_pb2_grpc
import freeconf.pb.fc_pb2
import freeconf.pb.fc_x_pb2_grpc
import freeconf.pb.fc_x_pb2
import freeconf.pb.fs_pb2_grpc
import freeconf.pb.fs_pb2
import freeconf.node
import freeconf.fs
import freeconf
import weakref
import platform
import distutils.spawn

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

def exe_fname():
    os = platform.system().lower()
    py_arch = platform.machine().lower()
    exe_ext = ""
    if os == "windows":
        exe_ext = ".exe"
    arch = {
        "x86_64": "amd64",
    }.get(py_arch, py_arch)
    fname = f'fc-lang-{freeconf.__version__}-{os}-{arch}{exe_ext}'
    return fname

def home_bin_dir():
    return os.path.expanduser("~/.freeconf/bin") # works on windows too


class ExecNotFoundException(Exception):
    def __init__(self, msg):
        super().__init__(msg)


def path_to_exe(verbose=False):
    """
    Rules for finding fc-lang exe.  We are someone flexible because we want to make it
    usable in all environments without too much hassle, but not at the expense of being
    too magical.  Hopefully this is the right balance.
    
        1. If explicitly set exact exec filename using FC_LANG_EXEC env var, use that only
        2. If explicitly set dir to set of exes using FC_LANG_DIR env var, use that only
        3. Look in ~/.freeconf/bin 
        4. Look in PATH
        5. Fail
    """
    file_path = os.environ.get('FC_LANG_EXEC', None)
    if file_path:
        if verbose:
            print(f"FC_LANG_EXEC set. checking {file_path}...")
        if not os.path.isfile(file_path):
            raise ExecNotFoundException(f"FC_LANG_EXEC={file_path} does not point to a valid file")
        return file_path
    elif verbose:
        print("FC_LANG_EXEC not set")

    fname = exe_fname()
    fc_lang_dir = os.environ.get('FC_LANG_DIR', None)
    if verbose:
        print(f"FC_LANG_DIR set" if fc_lang_dir != None else "FC_LANG_DIR not set, checking home dir")
    fc_dir = home_bin_dir() if fc_lang_dir == None else fc_lang_dir
    file_path = os.path.join(fc_dir, fname)
    if verbose:
        print(f"checking {file_path}")
    if os.path.isfile(file_path):
        return file_path
    elif fc_lang_dir != None:
        # if they explicitly set FC_LANG_DIR and we didn't find file, then exit as this is likely a misconfig
        raise ExecNotFoundException(f"{file_path} was not found in {fc_lang_dir}")
    
    if verbose:
        print("checking PATH")
    full_path = distutils.spawn.find_executable(fname)
    if not full_path:
        raise ExecNotFoundException(f"{fname} was not found in PATH or any of the other documented locations")

    return full_path


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
        exec_bin = path_to_exe()
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
        self.g_fs = freeconf.pb.fs_pb2_grpc.FileSystemStub(self.g_channel)
        self.fs = freeconf.fs.FileSystemServicer(self)

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
        if id == 0:
            raise Exception("0 id not valid")
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
