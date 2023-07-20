import freeconf.pb.fc_pb2
import freeconf.pb.fc_pb2_grpc
import freeconf.pb.fs_pb2
import freeconf.meta_decoder
import freeconf.handles
import freeconf.driver

## Parse YANG files into freeconf.meta.Module 

def load_module_file(ypath, name, driver=None):
    """
    Parse a YANG file and return parsed results as a freeconf.meta.Module.  

    :param ypath: freeconf.source.Source representing where to find YANG names
    :param name: name of the YANG module w/o the ".yang" file extention
    """
    d = driver if driver else freeconf.driver.shared_instance()
    req = freeconf.pb.fc_pb2.LoadModuleRequest(name=name)
    if ypath:
        req.sourceHnd = ypath.hnd
    resp = d.g_parser.LoadModule(req)
    m = freeconf.meta_decoder.Decoder().decode(resp.module)
    m.hnd = d.obj_weak.store_hnd(resp.moduleHnd, m)
    return m

def load_module_io(ypath, rdr, driver=None):
    """
    Parse a YANG file and return parsed results as a freeconf.meta.Module.  

    :param ypath: freeconf.source.Source representing where to find YANG names
    :param rdr: file-like reader with contents of YANG definition
    """
    d = driver if driver else freeconf.driver.shared_instance()
    stream = d.fs.new_rdr_io(rdr)
    req = freeconf.pb.fc_pb2.LoadModuleRequest(streamHnd=stream.hnd)
    if ypath:
        req.sourceHnd = ypath.hnd
    resp = d.g_parser.LoadModule(req)
    m = freeconf.meta_decoder.Decoder().decode(resp.module)
    m.hnd = d.obj_weak.store_hnd(resp.moduleHnd, m)
    return m

def load_module_str(ypath, module_str, driver=None):
    """
    Parse a YANG file and return parsed results as a freeconf.meta.Module.  

    :param ypath: freeconf.source.Source representing where to find YANG names
    :param module_str: contents of YANG definition
    """
    d = driver if driver else freeconf.driver.shared_instance()
    stream = d.fs.new_rdr_str(module_str)
    req = freeconf.pb.fc_pb2.LoadModuleRequest(streamHnd=stream.hnd)
    if ypath:
        req.sourceHnd = ypath.hnd
    resp = d.g_parser.LoadModule(req)
    m = freeconf.meta_decoder.Decoder().decode(resp.module)
    m.hnd = d.obj_weak.store_hnd(resp.moduleHnd, m)
    return m

def resolve_module(driver, module_hnd_id):
    m = driver.obj_weak.lookup_hnd(module_hnd_id)
    if m == None:
        req = freeconf.pb.fc_pb2.GetModuleRequest(moduleHnd=module_hnd_id)
        resp = driver.g_nodes.GetModule(req)
        m = freeconf.meta_decoder.Decoder().decode(resp.module)
        m.hnd = driver.obj_weak.store_hnd(module_hnd_id, m)
    return m
