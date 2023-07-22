import io
import time
import freeconf.pb.fc_pb2
import freeconf.pb.fs_pb2
import freeconf.node
import freeconf.meta
import freeconf.driver

def json_read_file(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    return __json_read(d.fs.new_rdr_file(fname), d)

def json_read_str(s, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    return __json_read(d.fs.new_rdr_str(s) , d)

def json_read_io(rdr, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_read(d.fs.new_rdr_io(rdr), d)

def __json_read(stream, driver):
    req = freeconf.pb.fc_pb2.JSONRdrRequest(streamHnd=stream.hnd)
    resp = driver.g_nodeutil.JSONRdr(req)
    return freeconf.handles.RemoteRef(driver, resp.nodeHnd)


def json_write_file(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_write(d.fs.new_wtr_file(fname), d)

def json_write_io(wtr, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_write(d.fs.new_wtr_io(wtr), d)

def json_write_str(sel, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    wtr = io.BytesIO()
    stream = d.fs.new_wtr_io(wtr)
    n = __json_write(stream, d)
    sel.upsert_into(n)
    d.g_fs.CloseStream(freeconf.pb.fs_pb2.CloseStreamRequest(streamHnd=stream.hnd))
    wtr.seek(0)
    return wtr.read().decode('UTF-8')

def __json_write(stream, driver):
    req = freeconf.pb.fc_pb2.JSONWtrRequest(streamHnd=stream.hnd)
    resp = driver.g_nodeutil.JSONWtr(req)
    return freeconf.handles.RemoteRef(driver, resp.nodeHnd)
