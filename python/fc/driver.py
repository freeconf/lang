import os
import os.path
import time
import subprocess
import grpc
import pb.fc_lang_pb2_grpc

class Driver():

    def __init__(self, sock_file=None, fallover_callback=None):
        self.fc_lang_proc = None
        cwd = os.getcwd()
        self.sock_file = sock_file if sock_file else f'{cwd}/fc-lang.sock'

    def load(self):
        if self.fc_lang_proc:
            raise Exception("fc-lang already loaded")
        self.fc_lang_proc = subprocess.Popen(['fc-lang', self.sock_file])
        self.wait_for_startup()
        self.channel = grpc.insecure_channel(f'unix://{self.sock_file}')
        self.stub = pb.fc_lang_pb2_grpc.DriverStub(self.channel)

    def wait_for_startup(self):
        for i in range(0, 10):
            if os.path.exists(self.sock_file):
                return
            time.sleep(0.01*(i*5))
            print(i)
        raise Exception("failure to start fc-yang")


    def unload(self):
        self.fc_lang_proc.terminate()
        self.fc_lang_proc = None

    def release(self, handle):
        req = pb.fc_lang_pb2.ReleaseRequest(handle=handle)
        self.stub.Release(req)

