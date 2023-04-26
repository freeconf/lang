class Basic():

    def __init__(self):
        self.hnd = 0

    def child(self, r):
        raise Exception(f'child not implemented in {r.path.str()}/{r.meta.ident}')


    def field(self, r, write_val):
        raise Exception(f'field not implemented in {r.path.str()}/{r.meta.ident}')


    def action(self, r):
        raise Exception(f'action not implemented in {r.path.str()}/{r.meta.ident}')


    def notification(self, r):
        raise Exception(f'notification not implemented in {r.path.str()}/{r.meta.ident}')

