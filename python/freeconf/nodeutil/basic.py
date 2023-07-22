class Basic():

    def __init__(self, on_context=None, on_release=None, on_child=None, on_next=None,                 
                 on_field=None, on_action=None, on_notification=None, on_begin_edit=None,
                 on_end_edit=None):        
        self.hnd = 0
        self.on_context = on_context
        self.on_release = on_release
        self.on_child = on_child
        self.on_next = on_next
        self.on_field = on_field
        self.on_action = on_action
        self.on_notification = on_notification
        self.on_begin_edit = on_begin_edit
        self.on_end_edit = on_end_edit

    def context(self, sel):
        if self.on_context != None:
            return self.on_context(sel)

    def release(self, sel):
        if self.on_release != None:
            self.on_release(sel)

    def child(self, r):
        if self.on_child != None:
            return self.on_child(r)
        raise Exception(f'child not implemented in {r.path.str()}/{r.meta.ident}')
    
    def next(self, r):
        if self.on_next != None:
            return self.on_next(r)
        raise Exception(f'next not implemented in {r.path.str()}/{r.meta.ident}')

    def field(self, r, write_val):
        if self.on_field != None:
            return self.on_field(r, write_val)
        raise Exception(f'field not implemented in {r.path.str()}/{r.meta.ident}')


    def action(self, r):
        if self.on_action != None:
            return self.on_action(r)
        raise Exception(f'action not implemented in {r.path.str()}/{r.meta.ident}')


    def notification(self, r):
        if self.on_notification != None:
            return self.on_notification(r)
        raise Exception(f'notification not implemented in {r.path.str()}/{r.meta.ident}')


    def begin_edit(self, r):
        if self.on_begin_edit != None:
            self.on_begin_edit(r)

    def end_edit(self, r):
        if self.on_edit_edit != None:
            self.on_edit_edit(r)
