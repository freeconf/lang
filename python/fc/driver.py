import os
import os.path
import time
import subprocess
import grpc
from concurrent import futures
import signal
import ctypes
libc = ctypes.CDLL("libc.so.6")

import pb.fc_lang_pb2_grpc
import pb.fc_x_pb2_grpc
import pb.fc_x_pb2
import fc.node

class Driver():

    def __init__(self, sock_file=None):
        self.fc_lang_proc = None
        cwd = os.getcwd()
        self.sock_file = sock_file if sock_file else f'{cwd}/fc-lang.sock'
        self.x_sock_file = f'{cwd}/fc-x.sock'
        if os.path.exists(self.sock_file):
            os.remove(self.sock_file)
        if os.path.exists(self.x_sock_file):
            os.remove(self.x_sock_file)


    def load(self):
        if self.fc_lang_proc:
            raise Exception("fc-lang already loaded")

        self.start_xclient_server()
        
        # emits this warning, still looking at clean way to remove this
        #  ResourceWarning: subprocess {pid} is still running
        self.fc_lang_proc = subprocess.Popen(['fc-lang', self.sock_file, self.x_sock_file], preexec_fn=exit_with_parent)

        self.wait_for_startup()
        self.channel = grpc.insecure_channel(f'unix://{self.sock_file}')
        self.stub = pb.fc_lang_pb2_grpc.HandlesStub(self.channel)

    def start_xclient_server(self):
        self.x_server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        pb.fc_x_pb2_grpc.add_XNodeServicer_to_server(fc.node.NodeHandler(), self.x_server)
        self.x_server.add_insecure_port(f'unix://{self.x_sock_file}')
        self.x_server.start()

    def wait_for_startup(self):
        for i in range(0, 20):
            if os.path.exists(self.sock_file):
                return
            time.sleep(0.01*(i*5))
        raise Exception("failure to start fc-yang")


    def unload(self):
        self.fc_lang_proc.terminate()
        self.fc_lang_proc = None

    def release(self, handle):
        req = pb.fc_lang_pb2.Handle(handle=handle)
        self.stub.Release(req)

# Ensure fc-lang is terminated when this python process is terminated
# see
#  https://stackoverflow.com/questions/19447603/how-to-kill-a-python-child-process-created-with-subprocess-check-output-when-t/19448096#19448096
#
# does not work on windows, not sure about mac, need to implement different method.
#
def exit_with_parent():
    return libc.prctl(1, signal.SIGTERM)
