import freeconf.meta
import freeconf.val

class Reflect():

    def __init__(self, obj, object_hook=None):
        self.hnd = 0
        self.obj = obj
        self.object_hook = object_hook
        self.is_dict = isinstance(obj, dict)

    def choose(self, sel, choice):
        for choice_case in choice.cases.values():
            for case_ddef in choice_case.definitions:
                if Reflect.has_value(self.obj, case_ddef):
                    return choice_case
        return None

    def context(self, sel):
        pass

    def release(self, sel):
        pass

    def child(self, r):
        child = None
        if r.delete:
            Reflect.clear_field(self.obj, r.meta)
            return None
        elif r.new:
            if self.object_hook:
                child = self.object_hook(r)            
            elif isinstance(r.meta, freeconf.meta.List) and ('key' in r.meta.key):
                child = []
            else:
                child = {}
            Reflect.write_field(self.obj, r.meta, child)
        else:
            child = Reflect.read_field(self.obj, r.meta)

        return None if child == None else self.new_child_node(child, r.meta)
    
    def new_child_node(self, child, meta):
        if isinstance(meta, freeconf.meta.List):
            return ReflectList(child, object_hook=self.object_hook)
        return Reflect(child, object_hook=self.object_hook)
    

    def new_object(self, r):
        if self.object_hook:
            return self.object_hook(r)
        if isinstance(r.meta, freeconf.meta.List):
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

    @classmethod
    def clear_field(cls, obj, meta):
        if isinstance(obj, dict):
            del obj[meta.ident]
        else:
            setattr(obj, meta.ident, None)

    @classmethod
    def has_value(cls, obj, meta):
        if isinstance(obj, dict):
            return meta.ident in obj
        else:
            hasattr(obj, meta.ident)

    def field(self, r, write_val):
        if r.write:
            if r.clear:
                Reflect.clear_field(self.obj, r.meta)
            else:
                Reflect.write_field(self.obj, r.meta, write_val.v)
        else:
            v = Reflect.read_field(self.obj, r.meta)
            if v != None:
                # TODO: coerse value...or let Go coerse it?
                return freeconf.val.Val(r.meta.type.format, v)

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

    def begin_edit(self, r):
        pass

    def end_edit(self, r):
        pass

class ReflectList():

    def __init__(self, objs, object_hook=None):
        self.hnd = 0
        self.objs = objs
        self.object_hook = object_hook
        self.is_dict = isinstance(objs, dict)
        self.is_list = isinstance(objs, list)
        if (not self.is_list) and (not self.is_dict):
            raise Exception(f'do not know how to manage list items in {type(self.objs)}')

    
    def new_object(self, r):
        if self.object_hook:
            return self.object_hook(r)
        return {}

    @classmethod
    def vals_to_key(cls, vals):
        if is_empty_list(vals):
            return None
        if len(vals) == 1:
            return vals[0].v
        v = []
        for k in vals:
            v.append(k.v)
        return tuple(v)
    
    @classmethod
    def read_fields(cls, obj, metas):
        if not is_empty_list(meta):
            vals = []
            for meta in metas:
                vals.append(Reflect.read_field(obj, meta))
        return vals

    @classmethod
    def write_fields(cls, obj, metas, vals):
        for i, val in enumerate(vals):
            Reflect.write_field(obj, metas[i], val)

    def next(self, r):
        found = None
        key = None

        if r.new:
            found = self.new_object(r)            
            if self.is_list:
                self.objs.append(found)
            else:
                if is_empty_list(r.key):
                    raise Exception(f'no key and do not know how to add to {type(self.objs)}')
                self.objs[self.vals_to_key(r.key)] = found

            if not is_empty_list(key):
                # while this will happen anyway, setting them on init means
                # subsequent sets will have these key values already set
                Reflect.write_fields(found, r.meta.keyMeta(), key)

        elif not is_empty_list(r.key):
            item_key = ReflectList.vals_to_key(r.key)
            if self.is_dict:
                if r.delete:
                    del self.objs[item_key]
                    return
                found = self.objs.get(item_key)
            else:
                # brute force iterate list until val matches
                for i, candidate in enumerate(self.objs):
                    candidate_vals = ReflectList.read_fields(candidate, r.meta.key)
                    candidate_key = ReflectList.vals_to_key(candidate_vals)
                    if item_key == candidate_key:
                        if r.delete:
                            del self.objs[i]
                            return
                        found = candidate
                        break

        elif r.row < len(self.objs):            
            if self.is_list:
                found = self.objs[r.row]
            else:
                # TODO: test if this is efficient when run w/large dict
                found = self.objs.keys()[r.row]
            key = ReflectList.read_fields(found, r.meta.key)

        if found != None:
            return Reflect(found, object_hook=self.object_hook), key
        return None, None

    def begin_edit(self, r):
        pass

    def end_edit(self, r):
        pass

def is_empty_list(list):
    return list == None or len(list) == 0
