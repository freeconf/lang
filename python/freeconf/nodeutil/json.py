import freeconf.pb.fc_pb2
import freeconf.node
import freeconf.meta
import freeconf.driver

def json_read(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    req = freeconf.pb.fc_pb2.JSONRdrRequest(fname=fname)
    resp = d.g_nodeutil.JSONRdr(req)
    return freeconf.handles.RemoteRef(d, resp.nodeHnd)


def json_write(fname, driver=None):
    d = driver if driver else freeconf.driver.shared_instance()
    req = freeconf.pb.fc_pb2.JSONWtrRequest(fname=fname)
    resp = d.g_nodeutil.JSONWtr(req)
    return freeconf.handles.RemoteRef(d, resp.nodeHnd)
