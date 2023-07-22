import io
import freeconf.pb.fs_pb2_grpc
import traceback
import threading

# 2048 is likely low, and should be tested, only limitation is max GRPC message size
BUFF_SIZE = 2048

class FileSystemServicer():

    def __init__(self, driver):
        self.driver = driver

    def new_rdr_str(self, s):
        b = bytes(s, 'utf-8')
        req = freeconf.pb.fs_pb2.ReaderInitRequest(inlineContent=b)
        resp = self.driver.g_fs.ReaderInit(req)
        return StreamRef(self.driver, resp.streamHnd)
    
    def new_rdr_file(self, fname):
        req = freeconf.pb.fs_pb2.ReaderInitRequest(fname=fname)
        resp = self.driver.g_fs.ReaderInit(req)
        return StreamRef(self.driver, resp.streamHnd)

    def new_rdr_io(self, rdr):
        req = freeconf.pb.fs_pb2.ReaderInitRequest(stream=True)
        resp = self.driver.g_fs.ReaderInit(req)
        return StreamReader(self.driver, resp.streamHnd, rdr)

    def new_wtr_io(self, wtr):
        req = freeconf.pb.fs_pb2.WriterInitRequest(stream=True)
        resp = self.driver.g_fs.WriterInit(req)
        return StreamWriter(self.driver, resp.streamHnd, wtr)

    def new_wtr_file(self, fname):
        req = freeconf.pb.fs_pb2.WriterInitRequest(fname=fname)
        resp = self.driver.g_fs.WriterInit(req)
        return StreamRef(self.driver, resp.streamHnd)        

class StreamRef():

    def __init__(self, driver, streamHnd):
        self.hnd = streamHnd
        driver.obj_strong.store_hnd(streamHnd, self)


class StreamReader():

    def __init__(self, driver, streamHnd, delegate_rdr):
        self.hnd = streamHnd
        self.driver = driver
        self.delegate_rdr = delegate_rdr
        self.keep_running = True
        driver.obj_strong.store_hnd(streamHnd, self)
        self.t = threading.Thread(target=self.run)
        self.t.start()

    def run(self):
        self.driver.g_fs.ReaderStream(self.readall())
    
    def readall(self):
        while self.keep_running:
            # should block
            chunk = self.delegate_rdr.read(BUFF_SIZE)
            if chunk == None or len(chunk) == 0:
                return
            yield freeconf.pb.fs_pb2.ReaderStreamData(streamHnd=self.streamHnd, chunk=chunk)
    
    def close(self):
        self.keep_running = False
        self.delegate_rdr.close()

    def close_and_wait(self):
        self.close()
        self.t.join()

class StreamWriter():

    def __init__(self, driver, streamHnd, delegate_wtr):
        self.hnd = streamHnd
        self.delegate_wtr = delegate_wtr
        self.driver = driver
        self.keep_running = True
        driver.obj_strong.store_hnd(streamHnd, self)
        self.t = threading.Thread(target=self.run)
        self.t.start()

    def run(self):
        chunks = self.driver.g_fs.WriterStream(freeconf.pb.fs_pb2.WriterStreamRequest(streamHnd=self.hnd))
        for data in chunks:
            n = self.delegate_wtr.write(data.chunk)
            if n != len(data.chunk):
                raise Exception(f"partial write not supported. requested {len(data.chunk)}, wrote {n} bytes")


    def wait(self):
        self.t.join()