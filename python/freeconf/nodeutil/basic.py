class Basic():

    def __init__(self):
        self.hnd = 0

    def context(self, sel):
        raise Exception(f'context not implemented in {sel.path.str()}')

    def release(self, sel):
        raise Exception(f'release not implemented in {sel.path.str()}')

    def child(self, r):
        raise Exception(f'child not implemented in {r.path.str()}/{r.meta.ident}')


    def field(self, r, write_val):
        raise Exception(f'field not implemented in {r.path.str()}/{r.meta.ident}')


    def action(self, r):
        raise Exception(f'action not implemented in {r.path.str()}/{r.meta.ident}')


    def notification(self, r):
        raise Exception(f'notification not implemented in {r.path.str()}/{r.meta.ident}')


    def begin_edit(self, r):
        pass

    def end_edit(self, r):
        pass
