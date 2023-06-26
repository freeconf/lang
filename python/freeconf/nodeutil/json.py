import freeconf.pb.fc_pb2
import freeconf.node
import freeconf.meta

def json_read(driver, fname):
    req = freeconf.pb.fc_pb2.JSONRdrRequest(fname=fname)
    resp = driver.g_nodeutil.JSONRdr(req)
    return freeconf.handles.RemoteRef(driver, resp.nodeHnd)


def json_write(driver, fname):
    req = freeconf.pb.fc_pb2.JSONWtrRequest(fname=fname)
    resp = driver.g_nodeutil.JSONWtr(req)
    return freeconf.handles.RemoteRef(driver, resp.nodeHnd)
