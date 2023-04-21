import fc.nodeutil

class Dump(fc.nodeutil.Basic):

    def __init__(self, target, out, padding=''):
        super(Dump, self).__init__()
        self.target = target
        self.out = out
        self.padding = padding
        self.write = out.write

    def child(self, r):
        self.write(f'{self.padding}{r.meta.ident}:\n')
        child = self.target.child(r)
        print(f'dump child new={r.new} {r.meta.ident}={child}')
        if child == None:
            self.write(", !found\n")
            return None
        self.write(", found")
        return Dump(self.target, self.out, self.padding + '  ')

    def field(self, r, hnd):
        if r.write:
            self.write(f'{self.padding}->{r.meta.ident}=({hnd.format}.{hnd.v}\n')
            self.target.field(r, hnd)
        else:
            self.target.field(r, hnd)
            self.write(f'{self.padding}<-{r.meta.ident}=({hnd.format}.{hnd.v}\n')
