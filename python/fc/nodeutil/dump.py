from fc.nodeutil import basic

class Dump(basic.Basic):

    def __init__(self, target, out, padding=''):
        super(Dump, self).__init__()
        self.target = target
        self.out = out
        self.padding = padding
        self.write = out.write

    def next(self, r):
        self.write(f'{self.padding}{r.meta.ident}:\n')
        next, key = self.target.next(r)
        if next == None:
            return None, key
        return Dump(next, self.out, self.padding), key

    def child(self, r):
        self.write(f'{self.padding}{r.meta.ident}:\n')
        child = self.target.child(r)
        if child == None:
            self.write("found=false\n")
            return None
        self.write("found=true\n")
        return Dump(child, self.out, self.padding + '  ')

    def field(self, r, hnd):
        if r.write:
            self.write(f'{self.padding}->{r.meta.ident}=({hnd.format.name.lower()}.{hnd.v})\n')
            self.target.field(r, hnd)
        else:
            self.target.field(r, hnd)
            self.write(f'{self.padding}<-{r.meta.ident}=({hnd.format.name.lower()}.{hnd.v})\n')
