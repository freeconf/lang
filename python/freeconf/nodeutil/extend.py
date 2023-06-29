class Extend():

    def __init__(self, base, on_child=None, on_field=None, on_action=None, on_notification=None, 
                 on_begin_edit=None, on_end_edit=None,
                 on_release=None, on_context=None):
        self.hnd = 0
        self.base = base
        self.on_child = on_child
        self.on_field = on_field
        self.on_action = on_action
        self.on_notification = on_notification
        self.on_begin_edit = on_begin_edit
        self.on_end_edit = on_end_edit
        self.on_release = on_release
        self.on_context = on_context

    def context(self, sel):
        if self.on_context:
            return self.on_context(self.base, sel)
        return self.base.context(sel)

    def release(self, sel):
        if self.on_release:
            return self.on_release(self.base, sel)
        return self.base.release(sel)

    def child(self, r):
        if self.on_child:
            return self.on_child(self.base, r)
        return self.base.child(r)


    def field(self, r, write_val):
        if self.on_field:
            return self.on_field(self.base, r, write_val)
        return self.base.field(r, write_val)


    def action(self, r):
        if self.on_action:
            return self.on_action(self.base, r)
        return self.base.action(r)


    def notification(self, r):
        if self.on_notification:
            return self.on_notification(self.base, r)
        return self.base.notification(r)


    def begin_edit(self, r):
        if self.on_begin_edit:
            return self.on_begin_edit(self.base, r)


    def end_edit(self, r):
        if self.on_end_edit:
            return self.on_end_edit(self.base, r)
        
    def choose(self, sel, choice):
        if self.on_end_edit:
            return self.on_end_edit(self.base, sel, choice)
        return self.base.choose(sel, choice)
