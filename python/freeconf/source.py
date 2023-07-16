import freeconf.driver
import freeconf.pb.fc_pb2

class Source:
    """
     Handle to collections of YANG files
    """
    def __init__(self, hnd):
        self.hnd = hnd

def path(p, driver=None):
    """
      series of directories to find files separated by ':' 
    """
    return __new_source(driver, freeconf.pb.fc_pb2.SourceRequest(path=p))

def any(*openers, driver=None):
    """
      append two or more sources together
    """
    hnds = [o.hnd for o in openers]
    return __new_source(driver, freeconf.pb.fc_pb2.SourceRequest(any=hnds))

def yang_internal_ypath(driver=None):
    """
      access internal YANG files from the freeconf "yang" package like 
      schema or docs 
    """
    return __new_source(driver, freeconf.pb.fc_pb2.SourceRequest(yangInternalYpath=True))

def restconf_internal_ypath(driver=None):
    """
      access internal YANG files from the freeconf "restconf" package like
      fc-restconf and ietf-inet-types among others
    """
    return __new_source(driver, freeconf.pb.fc_pb2.SourceRequest(restconfInternalYpath=True))

def __new_source(driver, req):
    d = driver if driver else freeconf.driver.shared_instance()
    resp = d.g_parser.Source(req)
    return Source(resp.sourceHnd)
