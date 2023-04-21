import pb.fc_pb2
import pb.fc_pb2_grpc
import pb.fc_x_pb2
import pb.fc_x_pb2_grpc
import fc.node

def json_rdr(driver, fname):
    req = pb.fc_pb2.JSONRdrRequest(fname=fname)
    resp = driver.g_nodeutil.JSONRdr(req)
    return fc.handles.RemoteRef(driver, resp.nodeHnd)

class Basic():

    def __init__(self):
        self.hnd = 0

    def child(self, req):
        raise Exception(f'child not implemented in {req.path.str()}/{req.meta.ident}')


    def field(self, req, write_val):
        raise Exception(f'field not implemented in {req.path.str()}/{req.meta.ident}')


    def action(self, req):
        raise Exception(f'action not implemented in {req.path.str()}/{req.meta.ident}')


    def notification(self, req):
        raise Exception(f'notification not implemented in {req.path.str()}/{req.meta.ident}')


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

        return None if child == None else Reflect(child)


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
                # TODO: coerse value
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


class Extend():

    def __init__(self, base, on_child=None, on_field=None, on_action=None, on_notification=None):
        self.hnd = 0
        self.base = base
        self.on_child = on_child
        self.on_field = on_field
        self.on_action = on_action
        self.on_notification = on_notification


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


    def notification(self, req):
        if self.on_notification:
            return self.on_notification(self.base, req)
        return self.base.notification(req)
