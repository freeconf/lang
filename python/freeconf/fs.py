import io
import freeconf.pb.fs_pb2_grpc
import traceback
import threading

# 2048 is likely low, and should be tested, only limitation is max GRPC message size
BUFF_SIZE = 2048

# Unclear if this is optimal or even below default threshold of result GRPC message
MAX_INLINE_CONTENT = 1024 * 2048

class FileSystemServicer(freeconf.pb.fs_pb2_grpc.FileSystemServicer):

    def __init__(self, driver):
        self.driver = driver
        self.handles = {}
        self.counter = 1

    def register_stream(self, stream):
        hnd = self.counter + 1
        self.counter = hnd
        self.handles[hnd] = stream
        return hnd
    
    def new_rdr_str(self, s):
        b = bytes(s, 'utf-8')
        if len(s) < MAX_INLINE_CONTENT:
            return freeconf.pb.fs_pb2.FileHandle(inlineContent=b)
        return self.new_file_handle_io(io.BytesIO(b))
    
    def new_rdr_io(self, rdr):
        # OPTIMIZE: inspect rdr for known types and see if inline is possible
        freeconf.pb.fs_pb2
        FileReader()
        hnd = self.register_stream(rdr)
        return freeconf.pb.fs_pb2.FileHandle(streamHnd=hnd)

    def new_rdr_file(self, fname):
        req = freeconf.pb.fs_pb2.ReaderRequest(fname=fname)
        resp = self.driver.g_fs.Read(req)
        return FileRef(self.driver, resp.streamHnd)

    def require_handle(self, hnd):
        try:
            return self.obj_strong(handles[hnd]
        except KeyError:
            raise Exception(f"stream handle {hnd} not registered")

    def ReadFile(self, req, context):
        rdr = None
        try:            
            rdr = self.require_handle(req.streamHnd)
            data = rdr.read(BUFF_SIZE)
            while data and len(data) > 0:
                yield freeconf.pb.fs_pb2.ReadFileResponse(chunk=data)
                data = rdr.read(BUFF_SIZE)
        except EOFError:
            pass
        except Exception as error:
            print(traceback.format_exc())
            raise error

    def WriteFile(self, req_iter, context):
        try:          
            for req in req_iter:
                wtr = self.require_handle(req.streamHnd)
                n = wtr.write(req.chunk)
                print(f"py: wrote {n} bytes to {wtr}")
                if n != len(req.chunk):
                    raise Exception("failure to write entire contents to file")
            print("py: exiting write")
            return freeconf.pb.fs_pb2.WriteFileResponse()
        except Exception as error:
            print(traceback.format_exc())
            raise error            

    def CloseFile(self, req, context):
        try:            
            print(f"py: closing {req.streamHnd}")
            stream = self.handles.pop(req.streamHnd, None)
            if stream:
                stream.close()
        except Exception as error:
            print(traceback.format_exc())
            raise error

class FileRef():

    def __init__(self, driver, hnd):
        self.hnd = hnd
        driver.obj_strong.store_hnd(hnd, self)


class FileReader():

    def __init__(self, client, delegate_rdr):
        self.client = client
        self.delegate_rdr = delegate_rdr
        self.keep_running = True

    def readall_nonblocking(self):
        self.t = threading.Thread(target=self.readall_blocking)
        self.t.start()

    def readall_blocking(self):
        while self.keep_running:
            # should block
            chunk = self.delegate_rdr.read(BUFF_SIZE)
            if chunk == None or len(chunk) == 0:
                return
            self.client.send(freeconf.pb.fs_pb2.FileReaderResponse(chunk=chunk))
    
    def close(self, blocking=False):
        self.keep_running = False
        if blocking:
            self.t.join()



def new_file_handle_str(s, driver):
    b = bytes(s, 'utf-8')
    if len(s) < MAX_INLINE_CONTENT:
        f = freeconf.pb.fs_pb2.FileHandle(inlineContent=b)
    else:
        rdr = io.BytesIO(b)
        hnd = driver.x_fs.register_stream(rdr)
        f = freeconf.pb.fs_pb2.FileHandle(streamHnd=hnd)
