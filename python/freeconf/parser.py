import freeconf.pb.fc_pb2
import freeconf.pb.fc_pb2_grpc
import freeconf.pb.fs_pb2
import freeconf.meta_decoder
import freeconf.handles
import freeconf.driver

class Parser():

    def __init__(self, driver=None):
        self.driver = driver if driver else freeconf.driver.shared_instance()

    def load_module_file(self, path, name):
        req = freeconf.pb.fc_pb2.LoadModuleRequest(path=path, name=name)
        resp = self.driver.g_parser.LoadModule(req)
        m = freeconf.meta_decoder.Decoder().decode(resp.module)
        m.hnd = self.driver.obj_weak.store_hnd(resp.moduleHnd, m)
        return m

    def load_module_io(self, path, rdr):
        stream = self.driver.fs.new_rdr_io(rdr)
        req = freeconf.pb.fc_pb2.LoadModuleRequest(path=path, streamHnd=stream.hnd)
        resp = self.driver.g_parser.LoadModule(req)
        m = freeconf.meta_decoder.Decoder().decode(resp.module)
        m.hnd = self.driver.obj_weak.store_hnd(resp.moduleHnd, m)
        return m

    def load_module_str(self, path, module_str):
        stream = self.driver.fs.new_rdr_str(module_str)
        req = freeconf.pb.fc_pb2.LoadModuleRequest(path=path, streamHnd=stream.hnd)
        resp = self.driver.g_parser.LoadModule(req)
        m = freeconf.meta_decoder.Decoder().decode(resp.module)
        m.hnd = self.driver.obj_weak.store_hnd(resp.moduleHnd, m)
        return m

    @classmethod
    def resolve_module(cls, driver, module_hnd_id):
        m = driver.obj_weak.lookup_hnd(module_hnd_id)
        if m == None:
            req = freeconf.pb.fc_pb2.GetModuleRequest(moduleHnd=module_hnd_id)
            resp = driver.g_nodes.GetModule(req)
            m = freeconf.meta_decoder.Decoder().decode(resp.module)
            m.hnd = driver.obj_weak.store_hnd(module_hnd_id, m)
        return m
