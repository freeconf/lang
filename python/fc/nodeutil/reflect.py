# import pb.fc_pb2_grpc
# import pb.fc_x_pb2
# import pb.fc_x_pb2_grpc
# import fc.node
import fc.meta
import fc.val

class Reflect():

    def __init__(self, obj, object_hook=None):
        self.hnd = 0
        self.obj = obj
        self.object_hook = object_hook
        self.is_dict = isinstance(obj, dict)

    def child(self, r):
        child = None
        if r.delete:
            if self.is_dict:
                del self.obj[r.meta.ident]
            else:
                Reflect.write_field(self.obj, r.meta, None)
            return None
        elif r.new:
            if self.object_hook:
                child = self.object_hook(r)            
            elif isinstance(r.meta, fc.meta.List) and ('key' in r.meta.key):
                child = []
            else:
                child = {}
            Reflect.write_field(self.obj, r.meta, child)            
        else:
            child = Reflect.read_field(self.obj, r.meta)

        return None if child == None else self.new_child_node(child, r.meta)
    
    def new_child_node(self, child, meta):
        if isinstance(meta, fc.meta.List):
            return ReflectList(child, object_hook=self.object_hook)
        return Reflect(child, object_hook=self.object_hook)
    

    def new_object(self, r):
        if self.object_hook:
            return self.object_hook(r)
        if isinstance(r.meta, fc.meta.List):
            return []
        return {}

    @classmethod
    def read_field(cls, obj, meta):
        if isinstance(obj, dict):
            return obj.get(meta.ident)
        return getattr(obj, meta.ident)
     
    @classmethod
    def write_field(cls, obj, meta, v):
        if isinstance(obj, dict):            
            obj[meta.ident] = v
        else:
            setattr(obj, meta.ident, v)

    def field(self, r, write_val):
        read_val = None
        if r.write:
            Reflect.write_field(self.obj, r.meta, write_val.v)
        else:
            v = Reflect.read_field(self.obj, r.meta)
            if v != None:
                # TODO: coerse value
                return fc.val.Val(r.meta.type.format, v)

        return None

    def action(self, r):
        if self.is_dict:
            raise Exception(f'cannot call functions on dicts in {r.path.str()}/{r.meta.ident}')

        f = getattr(self.obj, r.meta.ident)
        if not f:
            raise Exception(f'no function found for {r.path.str()}/{r.meta.ident}')

        input_data = {}
        if r.input:
            r.input.upsert_into(Reflect(input_data))
        resp = f(input_data)
        return Reflect(resp) if resp else None


class ReflectList():

    def __init__(self, objs, object_hook=None):
        self.hnd = 0
        self.objs = objs
        self.object_hook = object_hook
        self.is_dict = isinstance(objs, dict)
        self.is_list = isinstance(objs, list)
        
    def get_single_key(self, key):
        l = len(key) 
        if l == 0:
            return None
        if l == 1:
            return key[0].v
        raise Exception(f'do not know how to use compound keys to add to list')

    def require_single_key(self, key):
        l = len(key) 
        if l == 0:
            return None
        if l == 1:
            return key[0].v
        raise Exception(f'do not know how to use compound keys to add to list')
    
    def new_object(self, r):
        if self.object_hook:
            return self.object_hook(r)
        return {}

    def meta_key_to_item_key(self, key):
        if self.is_key_empty(key):
            return None
        if len(key) == 1:
            return key[0].v
        v = []
        for k in key:
            v.append(k.v)
        return tuple(v)
    
    def read_fields_by_metas(self, obj, metas):
        if not is_empty_list(meta):
            vals = []
            for meta in metas:
                vals.append(self.read_field(obj, meta))
        return vals

    def item_vals_to_meta_vals(self, vals, meta):
        meta_vals = []
        for v in vals:
            meta_vals.append(fc.val.Val(meta.type.format, v))
        return meta_vals

    
    def next(self, r):
        found = None
        key = None
        if (not self.is_list) and (not self.is_dict):
            raise Exception(f'do not know how to manage list items in {type(self.objs)}')
        
        if r.new:
            found = self.new_object()            
            if self.is_list:
                self.objs.append(found)
            else:
                if is_empty_list(r.key):
                    raise Exception(f'no key and do not know how to add to {type(self.objs)}')
                self.objs[self.meta_key_to_item_key(r.key)] = found

            if not self.is_key_empty(key):
                for k in r.key:
                    # while this will happen anyway, setting them on init means
                    # subsequent sets will have these key values already set
                    self.write_field(found, k)

        elif not is_empty_list(r.key):
            item_key = self.meta_key_to_item_key(r.key)
            if self.is_dict:
                if r.delete:
                    del self.objs[item_key]
                    return
                found = self.objs.get(item_key)
            else:
                # brute force iterate list until val matches then remove it
                for item in self.objs:
                    if item_key == self.read_key(item, r.meta.key):
                        found = item
                        break
        elif r.row < len(self.objs):            
            if self.is_list:
                found = self.objs[r.row]
            else:
                # TODO: test if this is efficient when run w/large map
                found = self.objs.keys()[r.row]
            key = self.read_key(found, r.meta.key)

        if found:
            return self.new_copy(found), key
        return None, None


def is_empty_list(self, list):
    return list == None or len(list) == 0
