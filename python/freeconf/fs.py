import io
import freeconf.pb.fs_pb2_grpc
import traceback

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
    
    def new_file_handle_str(self, s):
        b = bytes(s, 'utf-8')
        if len(s) < MAX_INLINE_CONTENT:
            return freeconf.pb.fs_pb2.FileHandle(inlineContent=b)
        return self.new_file_handle_io(io.BytesIO(b))
    
    def new_file_handle_io(self, rdr):
        # OPTIMIZE: inspect rdr for known types and see if inline is possible
        hnd = self.driver.x_fs.register_stream(rdr)
        return freeconf.pb.fs_pb2.FileHandle(streamHnd=hnd)

    def new_file_handle_file(self, fname):
        return freeconf.pb.fs_pb2.FileHandle(fname=fname)

    def require_handle(self, hnd):
        try:
            return self.handles[hnd]
        except KeyError:
            raise Exception(f"stream handle {hnd} not registered")

    def ReadFile(self, req, context):
        rdr = None
        try:            
            rdr = self.require_handle(req.streamHnd)
            data = rdr.read(BUFF_SIZE)
            while data and len(data) > 0:
                print(f"sending {len(data)} bytes")
                yield freeconf.pb.fs_pb2.ReadFileResponse(chunk=data)
                data = rdr.read(BUFF_SIZE)
        except EOFError:
            pass
        except Exception as error:
            print(traceback.format_exc())
            raise error
            

    def CloseFile(self, req, context):
        try:            
            rdr = self.handles.pop(req.streamHnd, None)
            if rdr:
                rdr.close()
        except Exception as error:
            print(traceback.format_exc())
            raise error


def new_file_handle_str(s, driver):
    b = bytes(s, 'utf-8')
    if len(s) < MAX_INLINE_CONTENT:
        f = freeconf.pb.fs_pb2.FileHandle(inlineContent=b)
    else:
        rdr = io.BytesIO(b)
        hnd = driver.x_fs.register_stream(rdr)
        f = freeconf.pb.fs_pb2.FileHandle(streamHnd=hnd)
