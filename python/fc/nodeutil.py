import pb.fc_g_pb2
import pb.fc_g_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.node

class NodeUtilService():

    def __init__(self, driver):
        self.driver = driver
        self.stub = pb.fc_g_pb2_grpc.NodeUtilStub(driver.channel)


    def json_rdr(self, fname):
        req = pb.fc_g_pb2.JSONRdrRequest(fname=fname)
        resp = self.stub.JSONRdr(req)
        return resp.nodeHnd


class Basic():

    def __init__(self):
        self.hnd = 0

    def child(self, req):
        raise Exception(f'on_child not implemented in {req.path.str()}/{req.meta.ident}')


    def field(self, req, write_val):
        raise Exception(f'on_field not implemented in {req.path.str()}/{req.meta.ident}')


    def action(self, req):
        raise Exception(f'on_action not implemented in {req.path.str()}/{req.meta.ident}')


class Reflect():

    def __init__(self, obj, object_hook=None):
        self.hnd = 0
        self.obj = obj
        self.object_hook = object_hook
        self.is_dict = isinstance(obj, dict)


    def child(self, req):
        child = None
        if req.delete:
            if self.is_dict:
                del self.obj[req.meta.ident]
            else:
                setattr(self.obj, req.meta.ident, None)
            return None
        elif req.new:
            if self.object_hook:
                child = self.object_hook(req)
            else:
                child = {}
        else:
            if self.is_dict:
                child = self.obj.get(req.meta.ident)
            else:
                child = getattr(self.obj, req.meta.ident)

        return Reflect(child) if child else None


    def field(self, req, write_val):
        read_val = None
        if req.write:
            if self.is_dict:
                self.obj[req.meta.ident] = write_val.v
            else:
                setattr(self.obj, req.meta.ident, write_val.v)
        else:
            v = None
            if self.is_dict:
                v = self.obj.get(req.meta.ident)
            else:
                v = getattr(self.obj, req.meta.ident)
            if v:
                read_val = fc.val.Val(req.meta.type.format, v)

        print(f'field obj={self.obj}, meta={req.meta.ident}, write={req.write}, read_val={read_val}')
        return read_val 


    def action(self, req):
        if self.is_dict:
            raise Exception(f'cannot call functions on dicts in {req.path.str()}/{req.meta.ident}')

        f = getattr(self.obj, req.meta.ident)
        if not f:
            raise Exception(f'no function found for {req.path.str()}/{req.meta.ident}')

        input_data = {}
        if req.input:
            req.input.upsert_into(Reflect(input_data))
        resp = f(input_data)
        return Reflect(resp) if resp else None


class Extend():

    def __init__(self, base, on_child=None, on_field=None, on_action=None):
        self.hnd = 0
        self.base = base
        self.on_child = on_child
        self.on_field = on_field
        self.on_action = on_action


    def child(self, req):
        if self.on_child:
            return self.on_child(self, self.base, req)
        return self.base.child(req)


    def field(self, req, write_val):
        if self.on_field:
            return self.on_field(self, self.base, req, write_val)
        return self.base.field(req, write_val)


    def action(self, req):
        if self.on_action:
            return self.on_action(self.base, req)
        return self.base.action(req)
