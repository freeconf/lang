from freeconf.nodeutil import reflect

# Wrap a node with a trace and see all data into and out of node structure without
# changing operation. 
#
# This follows FreeCONF Go's nodeutil/trace.go and must match same output format
# to be consistent but output format is used in unit tests
class Trace():

    def __init__(self, target, out, level=0):
        self.hnd = 0
        self.target = target
        self.out = out
        self.write = out.write
        self.level = level

    def context(self, sel):
        self.trace(self.level, 'context', "")

    def release(self, sel):
        self.trace(self.level, 'release', "")

    def next(self, r):
        if r.new:
            self.trace(self.level, f'next.new[{r.row}]', self.ident_str(r.meta, r.key))
        elif r.delete:
            self.trace(self.level, f'next.delete[{r.row}]', self.ident_str(r.meta, r.key))
        else:
            self.trace(self.level, f'next.read[{r.row}]', self.ident_str(r.meta, r.key))
        next, key = self.target.next(r)
        self.trace(self.level+1, "found", next!=None)
        if not r.new:
            self.trace_vals(self.level+1, "response.key", key)
        if next == None:
            return None, key        
        return Trace(next, self.out, level=self.level+1), key


    def meta_str(self, meta):
        if hasattr(meta, 'ident'):
            return meta.ident
        return str(type(meta))

    def child(self, r):
        if r.new:
            self.trace(self.level, "child.new", self.meta_str(r.meta))
        elif r.delete:
            self.trace(self.level, "child.delete", self.meta_str(r.meta))
        else:
            self.trace(self.level, "child.read", self.meta_str(r.meta))
        child = self.target.child(r)
        self.trace(self.level+1, "found", child!=None)
        if child == None:
            return None
        return Trace(child, self.out, level=self.level+1)


    def field(self, r, write_val):
        if r.write:
            self.trace(self.level, "field.write", r.meta.ident)
            if r.clear:
                self.trace(self.level+1, "clear", "true") 
            else:
                self.trace_val(self.level+1, "val", write_val)
            self.target.field(r, write_val)
        else:
            self.trace(self.level, "field.read", r.meta.ident)
            read_val = self.target.field(r, write_val)
            self.trace_val(self.level+1, "val", read_val)
            return read_val


    def begin_edit(self, r):
        self.trace(self.level, f"edit.begin", self.path_str(r.sel.path))
        if r.new:
            self.trace(self.level+1, "new", "true") 
        if r.delete:
            self.trace(self.level+1, "delete", r.delete)
        return self.target.begin_edit(r)


    def end_edit(self, r):
        self.trace(self.level, "edit.end", self.path_str(r.sel.path))
        if r.new:
            self.trace(self.level+1, "new", "true") 
        if r.delete:
            self.trace(self.level+1, "delete", r.delete)
        return self.target.end_edit(r)
    
    def action(self, r):
        self.trace(self.level, "action", r.meta.ident)
        if r.input is None:
            self.trace(self.level+1, "input", "nil")
        else:
            self.trace(self.level+1, "input", "true")
            
        output = self.target.action(r)
        if output is None:
            self.trace(self.level+1, "output", "nil")
        else:
            self.trace(self.level+1, "output", "true")
        return output
    

    def notify(self, r):
        self.trace(self.level, "notify", r.meta.ident)
        return self.target.notify(r)


    def choose(self, sel, choice):
        self.trace(self.level, "choose", choice.ident)
        choosen = self.target.choose(sel, choice)
        if choosen is None:
            self.trace(self.level+1, "choosen", "nil")
        else:
            self.trace(self.level+1, "choosen", choosen.ident)

        return choosen

    def trace(self, level, key, val):
        if val == None:
            val_str = "nil"
        elif isinstance(val, bool):
            val_str = "true" if val else "false"
        else:
            val_str = str(val)
        self.out.write(f'{"    "*level}{key}: {val_str}\n')


    def trace_vals(self, level, key, vals):
        if vals is None:
            self.trace(level, key, None)
        else:
            for i, val in enumerate(vals):
                self.trace_val(level, f'{key}[{i}]', val)


    def trace_val(self, level, key, val):
        if val == None:
            self.trace(level, key, None)
        else:
            self.trace(level, key, f'{val.format.name.lower()}({val.v})')


    def path_str(self, p):
        return self.ident_str(p.meta, p.key)


    def ident_str(self, meta, keys):
        if keys != None and len(keys) > 0:
            strs = []
            for key in keys:
                strs.append(str(key.v))
            return f'{self.meta_str(meta)}={",".join(strs)}'
        return self.meta_str(meta)
