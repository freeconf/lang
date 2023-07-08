import io
import time
import freeconf.pb.fc_pb2
import freeconf.pb.fs_pb2
import freeconf.node
import freeconf.meta
import freeconf.driver

def json_read_file(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    return __json_read(d.x_fs.new_file_handle_file(fname), d)

def json_read_str(s, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    return __json_read(d.x_fs.new_file_handle_str(s) , d)

def json_read_io(rdr, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_read(d.x_fs.new_file_handle_str(rdr), d)

def __json_read(f, d):
    req = freeconf.pb.fc_pb2.JSONRdrRequest(file=f)
    resp = d.g_nodeutil.JSONRdr(req)
    return freeconf.handles.RemoteRef(d, resp.nodeHnd)

def json_write_file(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_write(d.x_fs.new_file_handle_file(fname), d)

def json_write_io(wtr, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    return __json_write(d.x_fs.new_file_handle_io(wtr), d)

def json_write_str(sel, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()    
    wtr = io.BytesIO()
    n = __json_write(d.x_fs.new_file_handle_io(wtr), d)
    sel.upsert_into(n)
    #wtr.seek(0)
    s = wtr.read().decode('UTF-8')
    
    # BUG: This happens before python grpc has handled all the writes from Go's
    # WriteFile streaming call.  Go side is done in time, it is just the python
    # side recieving them is async and delayed.  #racecondition
    print(f"json_write_str done inserting, got {len(s)} bytes")

    return s


def __json_write(f, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    req = freeconf.pb.fc_pb2.JSONWtrRequest(file=f)
    resp = d.g_nodeutil.JSONWtr(req)
    return freeconf.handles.RemoteRef(d, resp.nodeHnd)

