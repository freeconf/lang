
import types
from re import sub
import copy
from freeconf import meta, val, node

class NodeOptions():

    def __init__(self, 
                 identities_as_strings=False,
                 enums_as_strings=False,
                 enums_as_ints=False,
                 ignore_empty=False,
                 try_plural_on_lists=False,
                 ident=None,
                 getter_prefix=None,
                 setter_prefix=None,
                 action_output_exploded=False,
                 action_input_exploded=False):        
        self.try_plural_on_lists = try_plural_on_lists
        self.ident = ident
        self.getter_prefix = getter_prefix
        self.setter_prefix = setter_prefix
        self.action_output_exploded = action_output_exploded
        self.action_input_exploded = action_input_exploded

        #TODO
        self.identities_as_strings = identities_as_strings
        self.enums_as_ints = enums_as_ints
        self.ignore_empty = ignore_empty
        self.enums_as_strings = enums_as_strings



class Node():

    def __init__(self, object, 
                 options=NodeOptions(),
                 on_options=None,
                 on_child=None, 
                 on_get_child=None,
                 on_new_child=None,
                 on_delete_child=None,                 
                 on_field=None, 
                 on_get_field=None,
                 on_set_field=None,
                 on_clear_field=None,
                 on_read=None,
                 on_write=None,
                 on_begin_edit=None, 
                 on_end_edit=None,
                 on_choose=None,
                 on_new_list_item=None,
                 on_get_by_key=None,
                 on_get_by_row=None,
                 on_delete_by_key=None,
                 on_action=None, 
                 on_notify=None, 
                 on_release=None,
                 on_new_object=None,
                 is_list_node=False):
        self.object = object
        self.on_options = on_options
        self.options = options
        self.on_field = on_field
        self.on_child = on_child
        self.on_get_child = on_get_child
        self.on_new_child = on_new_child
        self.on_delete_child = on_delete_child
        self.on_field = on_field
        self.on_get_field = on_get_field
        self.on_set_field = on_set_field
        self.on_clear_field = on_clear_field
        self.on_read = on_read
        self.on_write = on_write
        self.on_begin_edit = on_begin_edit
        self.on_end_edit = on_end_edit
        self.on_choose = on_choose
        self.on_new_list_item = on_new_list_item
        self.on_get_by_key = on_get_by_key
        self.on_get_by_row = on_get_by_row
        self.on_delete_by_key = on_delete_by_key
        self.on_action = on_action
        self.on_notify = on_notify
        self.on_release = on_release
        self.on_new_object = on_new_object
        self.hnd = None
        if is_list_node:
            self.new_list_handler()
        else:
            self.new_container_handler()

    def context(self, sel):
        pass 

    def release(self, sel):
        if self.on_release != None:
            self.on_release(self, sel)

    def child(self, r):
        if self.on_child != None:
            return self.on_child(self, r)
        return self.do_child(r)
    
    def do_child(self, r):
        if r.delete:
            if self.on_delete_child != None:
                self.on_delete_child(self, r)
                return
            self.do_delete_child(r)
            return
        if r.new:
            if self.on_new_child != None:
                return self.on_new_child(self, r)
            return self.do_new_child(r)
        if self.on_get_child != None:
            return self.on_get_child(self, r)
        return self.do_get_child(r)

    def get_options(self, meta):
        if self.on_options != None:
            opts = copy.copy(self.options)
            return self.on_options(self, meta, opts)
        return self.options
    
    def do_delete_child(self, r):
        self.container.clear(r.meta)

    def new(self, object, r=None):
        c = copy.copy(self)
        c.object = object
        c.hnd = None
        if r != None and Node.is_child_list(r):
            c.new_list_handler()
        else:
            c.new_container_handler()
        return c
    
    @classmethod
    def is_child_list(cls, r):
        return isinstance(r.meta, meta.List) and r.sel.path.meta != r.meta
    
    def new_list(self, list_object):
        c = copy.copy(self)
        c.object = list_object
        c.hnd = None
        c.new_list_handler()
        return c
    
    def new_object(self, m, inside_list):
        if self.on_new_object != None:
            return self.on_new_object(m, inside_list)
        return self.do_new_object(m, inside_list)


    def do_new_object(self, m, inside_list):
        if not inside_list and isinstance(m, meta.List):
            return []
        return {}

    
    def do_new_child(self, r):
        object = self.new_object(r.meta, Node.is_child_list(r))
        self.container.set(r.meta, object)
        return self.new(object, r)

    def do_get_child(self, r):
        obj = self.container.get(r.meta)
        if obj == None:
            return None
        
        if not r.new and self.options.ignore_empty and reflect_is_empty(obj):
            return None

        return self.new(obj, r)

    def new_list_handler(self):
        if isinstance(self.object, dict):
            self.list = DictionaryList(self.object)
        else:
            self.list = SliceList(self, self.object)

    def new_container_handler(self):
        if isinstance(self.object, dict):
            self.container = DictionaryContainer(self.object)
        else:
            self.container = ObjectContainer(self)

    def field(self, r, write_val):
        if self.on_field != None:
            return self.on_field(self, r, write_val)
        return self.do_field(r, write_val)
    
    def do_field(self, r, write_val):
        if r.clear:
            if self.on_clear_field != None:
                self.on_clear_field(self, r)
            else:
                self.do_clear_field(r)
        elif r.write:
            if self.on_set_field != None:
                self.on_set_field(self, r, write_val)
            else:                
                self.do_set_field(r, write_val)
        else:
            if self.on_get_field != None:
                return self.on_get_field(self, r)
            else:
                return self.do_get_field(r)
            
    def do_get_field(self, r):
        v = self.read_value(r.meta)
        if v == None:
            return None
        return val.Val.new(v, r.meta.type)

    def do_clear_field(self, r):
        self.container.clear(r.meta)

    def do_set_field(self, r, write_val):
        self.write_value(r.meta, write_val.v)        

    def read_value(self, meta):
        v = self.container.get(meta)
        if self.on_read != None:
            v = self.on_read(self, meta, v)
        return v

    def write_value(self, meta, v):
        if self.on_write != None:
            v = self.on_write(self, meta, v)
        self.container.set(meta, v)

    def choose(self, sel, choice):
        if self.on_choose != None:
            return self.on_choose(self, sel, choice)
        else:
            return self.do_choose(sel, choice)

    def do_choose(self, sel, choice):
        for choice_case in choice.cases.values():
            for case_ddef in choice_case.definitions:
                if self.exists(case_ddef):
                    return choice_case
        return None

    def exists(self, m):
        if isinstance(m, meta.List) or isinstance(m, meta.Container):
            r = node.ChildRequest(None, m, False, False)
            return self.child(r) != None
        r = node.FieldRequest(None, m, False, False)
        return self.field(r, None) != None
    
    def begin_edit(self, r):
        if self.on_begin_edit != None:
            self.on_begin_edit(self, r)

    def do_begin_edit(self, r):
        pass
        
    def end_edit(self, r):
        if self.on_begin_edit != None:
            self.on_end_edit(self, r)

    def do_end_edit(self, r):
        pass

    def next(self, r):
        if r.new:
            if self.on_new_list_item != None:
                return self.on_new_list_item(self, r), r.key
            else:
                return self.do_new_list_item(r), r.key
        if r.key != None and len(r.key) > 0:
            if r.delete:
                if self.on_delete_by_key != None:
                    self.on_delete_by_key(self, r)
                else:
                    self.do_delete_by_key(r)
                return
            if self.on_get_by_key != None:
                return self.on_get_by_key(self, r), r.key
            else:
                return self.do_get_by_key(r), r.key
        if self.on_get_by_row != None:
            return self.on_get_by_row(self, r)
        else:
            return self.do_get_by_row(r)


    def do_new_list_item(self, r):
        item = self.list.new_list_item(r)
        if item == None:
            return None
        return self.new(item)

    def do_delete_by_key(self, r):
        self.list.delete_by_key(r)

    def do_get_by_key(self, r):
        item = self.list.get_by_key(r)
        if item == None:
            return None
        return self.new(item)

    def do_get_by_row(self, r):
        item, key = self.list.get_by_row(r)
        if item == None:
            return None, None        
        return self.new(item), key

    def action(self, r):
        if self.on_action != None:
            return self.on_action(self, r)
        else:
            return self.do_action(r)
        
    def do_action(self, r):
        h = ActionHandler()
        opts = self.get_options(r.meta)
        resp = h.do(self, r, opts)
        if resp == None:
            return None
        return self.new(resp)

    def notify(self, r):
        if self.on_notify != None:
            return self.on_notify(self, r)
        raise NotImplemented

def reflect_is_empty(obj):
    if isinstance(obj, list) or isinstance(obj, dict) or isinstance(obj, str):
        return len(obj) == 0
    if isinstance(obj, bool):
        return True
    return obj == None

class ActionHandler():
        

    def do(self, n, r, opts):
        candidate = opts.ident
        if candidate == None:
            candidate = snake_case(r.meta.ident)
        m = getattr(n.object, candidate, None)
        if m == None or type(m) != types.MethodType:
            raise Exception(f"could not find function '{candidate}' on '{object}'")

        if r.input != None:
            input = n.new_object(r.meta.input, False)
            r.input.upsert_into(n.new(input))

            if opts.action_input_exploded:
                explode = []
                for d in r.meta.input.definitions:
                    explode.append(input.get(d.ident, None))
                resp = m(*explode)
            else:
                resp = m(input)
        else:
            resp = m()

        if resp and r.meta.output != None:
            if opts.action_output_exploded:
                implode = {}
                for i, d in enumerate(r.meta.output.definitions):
                    implode[d.ident] = resp[i]
                resp = implode

        return resp




class DictionaryList():

    def __init__(self, list):
        self.list = list
        self.index = None

    def new_list_item(self, r):
        item = {}
        self.list[self.key_val(r)] = item
        return item

    def delete_by_key(self, r):
        self.list.pop(self.key_val(r), None)

    def get_by_key(self, r):
        return self.list.get(self.key_val(r), None)

    def get_by_row(self, r):
        if self.index == None:
            self.index = [x for x in self.list.keys()]
            self.index.sort()
        if r.row < len(self.index):
            key = self.index[r.row]
            return self.list[key], [val.Val(key)]
        return None, None

    def key_val(self, r):
        if r.key == None:
            raise Exception(f"no key given for {r.path}")
        if len(r.key) != 1:
            raise Exception(f"expected single key for {r.path}")
        return r.key[0].v


class SliceList():

    def __init__(self, n, list):
        self.n = n
        self.list = list

    def new_list_item(self, r):
        item = {}
        self.list.append(item)
        return item

    def delete_by_key(self, r):
        ndx = self.find_by_key(r)
        if ndx >= 0:
            del self.list[ndx]

    def get_by_key(self, r):
        ndx = self.find_by_key(r)
        if ndx >= 0:
            return self.list[ndx]
        return None

    def get_by_row(self, r):
        if r.row >= 0 and r.row < len(self.list):
            list_item = self.list[r.row]
            key = self.get_key(list_item, r.meta.key_meta())
            return list_item, key
        return None, None

    def find_by_key(self, r):
        key_meta = r.meta.key_meta()
        if len(key_meta) == 0:
            raise Exception(f"{r.sel.path} has no keys defined")
        if len(key_meta) != len(r.key):
            raise Exception(f"{r.sel.path} requires {len(key_meta)} keys but {len(r.key)} given")
        for i, candidate in enumerate(self.list):
            key = self.get_key(candidate, key_meta)
            for j, k in enumerate(key):
                if k.v != r.key[j].v:
                    break
                last_key = j == len(key) - 1
                if last_key:
                    return i
        return -1
    
    def get_key(self, list_item, key_meta):
        if len(key_meta) == 0:
            return None
        
        key = []
        for meta in key_meta:
            child_node = self.n.new(list_item)
            fr = node.FieldRequest(None, meta, False, False)
            key_val = child_node.field(fr, None)
            if key_val == None:
                raise Exception("missing key value")
            key.append(key_val)
        return key



class DictionaryContainer():
    """Reads and writes to a dict"""

    def __init__(self, object):
        self.object = object

    def clear(self, meta):
        self.object.pop(meta.ident, None)

    def get(self, meta):
        return self.object.get(meta.ident, None)

    def set(self, meta, v):
        self.object[meta.ident] = v


class ObjectContainer():
    """Reads and writes to a object instance created from a class"""

    def __init__(self, node):
        self.node = node
        self.field_handlers = {}


    def clear(self, meta):
        self.field_handler(meta).clear()


    def get(self, meta):
        return self.field_handler(meta).get()


    def set(self, meta, v):
        self.field_handler(meta).set(v)


    def field_handler(self, meta):
        h = self.field_handlers.get(meta.ident, None)
        if h == None:
            h = self.new_field_handler(meta)
            if h == None:
                raise Exception(f"could not find field '{meta.ident}' on '{self.node.object}'")
            self.field_handlers[meta.ident] = h
        return h


    def new_field_handler(self, meta):
        opts = self.node.get_options(meta)
        h = FieldHandler(self.node.object, meta, opts)

        # support @ tags for properties?

        for candidate in FieldHandler.field_name_candidates(opts, meta):
            try:
                f = getattr(self.node.object, candidate)
                if type(f) != types.MethodType:
                    # found something, but continue on to look for getters or setters
                    h.field = candidate
                    break
            except AttributeError:
                pass


        if opts.getter_prefix:
            getter_prefix = opts.getter_prefix
        else:
            getter_prefix = "get_"
        for candidate in FieldHandler.accessor_name_candidates(opts, meta, getter_prefix):
            getter = getattr(self.node.object, candidate, None)
            if getter != None and type(getter) == types.MethodType:
                h.getter = getter
                break

        if opts.setter_prefix:
            setter_prefix = opts.setter_prefix
        else:
            setter_prefix = "set_"
        for candidate in FieldHandler.accessor_name_candidates(opts, meta, setter_prefix):
            setter = getattr(self.node.object, candidate, None)
            if setter != None and type(setter) == types.MethodType:
                h.setter = setter
                break

        if h.field == None and h.getter == None and h.setter == None:
            return None
        return h


class FieldHandler():
    
    def __init__(self, object, m, opts):
        self.object = object
        self.m = m
        self.opts = opts
        self.field = None
        self.getter = None
        self.setter = None

    @classmethod
    def field_name_candidates(cls, opts, m):
        if opts.ident != None:
            return [opts.ident]
        candidates = [snake_case(m.ident)]
        if opts.try_plural_on_lists:
            if isinstance(m, meta.List) or isinstance(m, meta.LeafList):
                candidates.append(candidates[0]+"s")
        return candidates


    @classmethod
    def accessor_name_candidates(cls, opts, m, prefix):
        field_candidates = FieldHandler.field_name_candidates(opts, m)
        candidates = []
        for candidate in field_candidates:
            candidates.append(prefix+candidate)
        candidates.extend(field_candidates)
        return candidates


    def clear(self):
        self.set(None)


    def get(self):        
        if self.getter != None:
            v = self.getter()
        else:
            v = getattr(self.object, self.field)
        if self.opts.ignore_empty and reflect_is_empty(v):
            return None
        return v

    def set(self, v):
        if self.setter != None:
            self.setter(v)
        else:
            setattr(self.object, self.field, v)


def snake_case(s):
    return '_'.join(
        sub('([A-Z][a-z]+)', r' \1',
        sub('([A-Z]+)', r' \1',
        s.replace('-', ' '))).split()).lower()