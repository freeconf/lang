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


class Basic(fc.node.Node):

    def __init__(self, on_child=None, on_field=None, on_action=None):
        super(Basic, self).__init__()
        self.on_child = on_child
        self.on_field = on_field
        self.on_action = on_action


    def child(self, req):
        if not self.on_child:
            raise Exception(f'on_child not implemented in {req.path.str()}/{req.meta.ident}')
        return self.on_child(self, req)


    def field(self, req, write_val):
        if not self.on_field:
            raise Exception(f'on_field not implemented in {req.path.str()}/{req.meta.ident}')
        return self.on_field(self, req, write_val)


    def action(self, req):
        if not self.on_action:
            raise Exception(f'on_action not implemented in {req.path.str()}/{req.meta.ident}')
        return self.self.on_action(self, req)


class Reflect(fc.node.Node):

    def __init__(self, obj, object_hook=None):
        super(Reflect, self).__init__()
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


class Extend(fc.node.Node):

    def __init__(self, base, on_child=None, on_field=None, on_action=None):
        super(Extend, self).__init__()
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
            return self.on_action(self, self.base, req)
        return self.base.action(req)
