from fc.nodeutil import basic

# Wrap a node with a trace and see all data into and out of node structure without
# changing operation. 
#
# This follows FreeCONF Go's nodeutil/trace.go and must match same output format
# to be consistent but output format is used in unit tests
class Trace(basic.Basic):

    def __init__(self, target, out, level=0):
        super(Trace, self).__init__()
        self.target = target
        self.out = out
        self.write = out.write
        self.level = level

    def next(self, r):
        if r.new:
            self.trace(self.level, f'next.new[{r.row}]', self.ident(r.sel.path))
        elif r.delete:
            self.trace(self.level, f'next.delete[{r.row}]', self.ident(r.sel.path))
        else:
            self.trace(self.level, f'next.read[{r.row}]', self.ident(r.sel.path))
        next, key = self.target.next(r)
        if next == None:
            return None, key
        return Trace(next, self.out, level=self.level+1), key

    def child(self, r):
        if r.new:
            self.trace(self.level, "child.new", self.ident(r.sel.path))
        elif r.delete:
            self.trace(self.level, "child.delete", self.ident(r.sel.path))
        else:
            self.trace(self.level, "child.read", self.ident(r.sel.path))
        child = self.target.child(r)
        self.trace(self.level+1, "found", child!=None)
        if child == None:
            return None
        return Trace(child, self.out, level=self.level+1)

    def field(self, r, write_val):
        if r.write:
            self.trace(self.level, "field.write", r.meta.ident)
            self.trace_val(self.level+1, "val", write_val)
            self.target.field(r, write_val)
        else:
            self.trace(self.level, "field.read", r.meta.ident)
            read_val = self.target.field(r, write_val)
            self.trace_val(self.level+1, "val", read_val)
            return read_val


    def trace(self, level, key, val):
        if val == None:
            val_str = "nil"
        elif isinstance(val, bool):
            val_str = "true" if val else "false"
        else:
            val_str = str(val)
        self.out.write(f'{"    "*level}{key}: {val_str}\n')


    def trace_val(self, level, key, val):
        if val == None:
            self.trace(level, key, "nil")
        else:
            self.trace(level, key, f'{val.format.name.lower()}({val.v})')


    def ident(self, path):
        if path.key != None and len(path.key) > 0:
            strs = []
            for key in path.key:
                strs.append(str(key.v))
            return f'{path.meta.ident}={strs.join(",")}'
        return path.meta.ident

