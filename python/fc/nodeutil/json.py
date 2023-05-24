import pb.fc_pb2
import fc.node
import fc.meta

def json_read(driver, fname):
    req = pb.fc_pb2.JSONRdrRequest(fname=fname)
    resp = driver.g_nodeutil.JSONRdr(req)
    return fc.handles.RemoteRef(driver, resp.nodeHnd)

