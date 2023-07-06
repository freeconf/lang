import io
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


def json_write(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    req = freeconf.pb.fc_pb2.JSONWtrRequest(fname=fname)
    resp = d.g_nodeutil.JSONWtr(req)
    return freeconf.handles.RemoteRef(d, resp.nodeHnd)

